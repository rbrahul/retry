package retry

import (
	"errors"
	"math"
	"reflect"
	"time"
)

type Retryable func() bool
type BackoffFunc func(uint64) uint64

type retryManager struct {
	maxRetry    int64
	backoffFn   BackoffFunc
	retryUntil  time.Duration
	startedAt   time.Time
	delay       time.Duration
	lastBackoff uint64 // in number of seconds
}

var (
	ErrMaxRetryOrRetryUntilInvalidArg = errors.New("invalid argument type. maxRetry can be either integer, time.Duration or func(uint64) uint64")
	ErrDelayOrBackOffFuncInvalidArg   = errors.New("invalid argument type. delay can be either time.Duration or `func(uint64) uint64`")
	ErrDeadlineExceeded               = errors.New("retry deadline has been exceeded")
	ErrMaximumRetryExceeded           = errors.New("maximum retry has been exceeded")
)

const defaultDelayDuration = 1 * time.Second

func (rm *retryManager) parseParams(args ...interface{}) error {
	if len(args) > 0 {
		firstArgKind := reflect.TypeOf(args[0]).Kind()
		if reflect.TypeOf(args[0]).String() == "time.Duration" {
			maxDuration := args[0].(time.Duration)
			rm.retryUntil = maxDuration
		} else if isIntKind(firstArgKind) {
			rm.maxRetry = int64(math.Abs(float64(reflect.ValueOf(args[0]).Int())))
		} else {
			return ErrMaxRetryOrRetryUntilInvalidArg
		}

		// delay in time.Duration or backOffFunc as func(uint64) uint64
		if len(args) > 1 {
			if reflect.TypeOf(args[1]).String() == "time.Duration" {
				rm.delay = args[1].(time.Duration)
			} else if reflect.TypeOf(args[1]).String() == "func(uint64) uint64" {
				rm.backoffFn = args[1].(func(uint64) uint64)
			} else {
				return ErrDelayOrBackOffFuncInvalidArg
			}
		}
	}
	if rm.backoffFn == nil && rm.delay == 0 {
		rm.delay = defaultDelayDuration
	}
	return nil
}

func (rm *retryManager) addDelay() {
	var delayInBetween time.Duration
	if rm.lastBackoff == 0 {
		rm.lastBackoff = 1
	}
	if rm.backoffFn != nil {
		numberOfSeconds := rm.backoffFn(rm.lastBackoff)
		rm.lastBackoff = numberOfSeconds
		delayInBetween = time.Duration(numberOfSeconds) * time.Second
	} else {
		delayInBetween = rm.delay
	}
	withJitter := addJitter(delayInBetween)
	time.Sleep(withJitter)
}

func (rm *retryManager) execute(fn Retryable) error {
	shouldRetry := fn()
	for shouldRetry {
		if rm.retryUntil > 0 {
			rm.addDelay()
			deadLineExceeded := time.Now().After(rm.startedAt.Add(rm.retryUntil))
			if deadLineExceeded {
				return ErrDeadlineExceeded
			}
			shouldRetry = fn()
			continue
		}
		// If maxRetry is set 5 then 5-1 time will be retried. Because initially function was already excuted once.
		// Which means: 1 (Initial Call) + 4 retries == 5 maxRetry
		if rm.maxRetry > 0 {
			rm.maxRetry -= 1
			if rm.maxRetry == 0 {
				return ErrMaximumRetryExceeded
			}
		}

		rm.addDelay()
		shouldRetry = fn()
	}

	return nil
}

// Retry excutes the function fn in an interval specified as delay or backoffFunc. And it will run until it reaches the maximum number of times specified as maxNumberOfRetry or retryUntil as duration. Retry will be terminated if the fn returns false.
//
// Accepted Parameters:
//
// Retry(fn func() bool, maxNumberOfRetry int | retryUntil time.Duration,  delay time.Duration | backOffFn func(uint64) uint64))
//
// Example:
//
//	err := retry.Retry(func() bool {
//			fmt.Println("Hello world!")
//			return true
//		}, 10, 2*time.Second)
//
// Default Parameters:
//
//	maxNumberOfRetry -> Infinity
//	delay -> 1 * time.Second
//
// Usage:
//
// Retry(fn func() bool, maxNumberOfRetry int, delay time.Duration)
//
// Retry(fn func() bool, retryUntil time.Duration, delay time.Duration)
//
// Retry(fn func() bool, retryUntil time.Duration, backOffFn func(uint64) uint64)
//
// Retry(fn func() bool, maxNumberOfRetry int)
//
// Retry(fn)
func Retry(fn Retryable, args ...interface{}) error {
	retryManager := &retryManager{
		startedAt: time.Now(),
	}
	err := retryManager.parseParams(args...)
	if err != nil {
		return err
	}
	return retryManager.execute(fn)
}

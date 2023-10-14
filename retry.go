package retry

import (
	"errors"
	"fmt"
	"math"
	"reflect"
	"time"
)

// implement retry function who can retry a function if any

// retry() --> true, data,

type Retryable func() bool
type BackOffFn func(uint64) uint64

type retryManger struct {
	maxRetry    int64
	delayerFn   BackOffFn
	retryUntil  time.Duration
	startedAt   time.Time
	delay       time.Duration
	lastBackOff uint64 // in number of seconds
}

const defaultDelayDuration = 1 * time.Second

// Retry(retriableFn func() bool, maxNumberOfRetry int | retryUntil time.Duration,  delay time.Duration | backOffFn func(uint64) uint64))
// Default:
//		maxNumberOfRetry -> Infinity
// 		delay -> 1 * time.Second

// Usage:
// Retry(retriableFn func() bool, maxNumberOfRetry int | retryUntil time.Duration,  delay time.Duration | backOffFn func(uint64) uint64))
// Retry(retriableFn func() bool, maxNumberOfRetry int, delay time.Duration)
// Retry(retriableFn func() bool, retryUntil time.Duration, delay time.Duration)
// Retry(retriableFn func() bool, retryUntil time.Duration, backOffFn func(uint64) uint64)
// Retry(retriableFn func() bool, maxNumberOfRetry int)
// Retry(retriableFn, backOffFn func(uint64) uint64)) [If 2nd argument is function then it will be treated as backOffFn and maxNumberOfRetry will be considered Infinity]
// Retry(retriableFn)

func (rm *retryManger) parseParams(args ...interface{}) error {
	if len(args) > 0 {
		firstArgKind := reflect.TypeOf(args[0]).Kind()
		if reflect.TypeOf(args[0]).String() == "time.Duration" {
			maxDuration := args[0].(time.Duration)
			rm.retryUntil = maxDuration
		} else if isIntKind(firstArgKind) {
			rm.maxRetry = int64(math.Abs(float64(reflect.ValueOf(args[0]).Int())))
		} else {
			return errors.New("invalid argument type. maxRetry can be either integer or time.Duration")
		}

		if len(args) == 1 && reflect.TypeOf(args[0]).String() == "func(uint64) uint64" {
			rm.delayerFn = args[0].(func(uint64) uint64)
		}

		// delay in duration /backOff as func --> exponentialBackOff, randInt, customBackOffFn
		if len(args) > 1 {
			if reflect.TypeOf(args[1]).String() == "time.Duration" {
				rm.delay = args[1].(time.Duration)
			} else if reflect.TypeOf(args[1]).String() == "func(uint64) uint64" {
				rm.delayerFn = args[1].(func(uint64) uint64)
			} else {
				return errors.New("invalid argument type for delay. delay can be either time.Duration or `func(uint64) uint64`")
			}
		}
	}
	if rm.delayerFn == nil && rm.delay == 0 {
		rm.delay = defaultDelayDuration
	}
	return nil
}

func (rm *retryManger) addDelay() {
	var delayInBetween time.Duration
	if rm.lastBackOff == 0 {
		rm.lastBackOff = 1
	}
	if rm.delayerFn != nil {
		numberOfSeconds := rm.delayerFn(rm.lastBackOff)
		rm.lastBackOff = numberOfSeconds
		delayInBetween = time.Duration(numberOfSeconds) * time.Second
	} else {
		delayInBetween = rm.delay
	}
	fmt.Println("delayInBetween:", delayInBetween)
	withJitter := addJitter(delayInBetween)
	fmt.Println("With Jitter:", withJitter)
	time.Sleep(withJitter)
}

func (rm *retryManger) execute(fn Retryable) error {
	shouldRetry := fn()
	for shouldRetry {
		if rm.retryUntil > 0 {
			deadLineExceeded := time.Now().After(rm.startedAt.Add(rm.retryUntil))
			if deadLineExceeded {
				return errors.New("retry deadline has been exceeded")
			}
			rm.addDelay()
			shouldRetry = fn()
			continue
		}

		if rm.maxRetry > 0 {
			rm.maxRetry -= 1
			if rm.maxRetry == 0 {
				return errors.New("maximum retry has been exceeded")
			}
		}

		rm.addDelay()
		shouldRetry = fn()
	}

	return nil
}

func Retry(fn Retryable, args ...interface{}) error {
	retryManager := &retryManger{
		startedAt: time.Now(),
	}
	err := retryManager.parseParams(args...)
	if err != nil {
		return err
	}
	return retryManager.execute(fn)
}

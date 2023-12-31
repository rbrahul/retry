package retry

import (
	"math/rand"
	"reflect"
	"time"
)

func isIntKind(kind reflect.Kind) bool {
	intKinds := []reflect.Kind{reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64}
	for _, kindItem := range intKinds {
		if kindItem == kind {
			return true
		}
	}
	return false
}

func randomer() *rand.Rand {
	seed := rand.NewSource(time.Now().UnixNano())
	return rand.New(seed)
}

func applyJitter(delay time.Duration) time.Duration {
	// Added jitter between 10% to 30%
	return delay + time.Duration(float64(delay)*randomFloatWithinRange(.1, .3))
}

var addJitter = applyJitter

func randomIntWithinRange(lower, upper int) uint64 {
	return uint64(lower + randomer().Intn(upper-lower))
}

func randomFloatWithinRange(lower, upper float64) float64 {
	return lower + rand.Float64()*(upper-lower)
}

// It returns a random integer value between lower and upper range. And will be considered as the delay in the number of seconds.
//
// For Example:
//
// RandomBackoff(3,7) will return a dealy between 3 and 7
func RandomBackoff(lower, upper int) func(uint64) uint64 {
	return func(_ uint64) uint64 {
		return randomIntWithinRange(lower, upper)
	}
}

// It returns an exponential delay as number of seconds for each retry execution. If delay is > maxBackoff then maxBackoff will be the used as delay for next retry.
//
// For Example:
//
// ExponentialBackoff(10) will return the following intervals in number seconds as long as the retry is executed.
//
// [1 -> 2 -> 4 -> 8 -> 10 -> 10 ...]
func ExponentialBackoff(maxBackoff int) func(uint64) uint64 {
	return func(lastBackOff uint64) uint64 {
		exponentialbackOff := lastBackOff * 2
		if exponentialbackOff > uint64(maxBackoff) {
			return uint64(maxBackoff)
		}
		return exponentialbackOff
	}
}

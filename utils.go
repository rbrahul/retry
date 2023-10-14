package retry

import (
	"fmt"
	"math/rand"
	"reflect"
	"time"
)

func isIntKind(kind reflect.Kind) bool {
	fmt.Println("KIND:", kind)
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

func addJitter(delay time.Duration) time.Duration {
	// Added jitter between 10% to 30%
	return delay + time.Duration(float64(delay)*randomFloatWithinRange(.1, .3))
}

func randomIntWithinRange(lower, upper int) uint64 {
	return uint64(lower + randomer().Intn(upper-lower))
}

func randomFloatWithinRange(lower, upper float64) float64 {
	return lower + rand.Float64()*(upper-lower)
}

func RandInt(lower, upper int) func(uint64) uint64 {
	return func(_ uint64) uint64 {
		return randomIntWithinRange(lower, upper)
	}
}

func ExponentialBackoff(maxBackOff int) func(uint64) uint64 {
	return func(lastBackOff uint64) uint64 {
		exponentialbackOff := lastBackOff * 2
		if exponentialbackOff > uint64(maxBackOff) {
			return uint64(maxBackOff)
		}
		return exponentialbackOff
	}
}

func Parcentage(percentage float32, maxBackOff int) func(uint64) uint64 {
	return func(lastBackOff uint64) uint64 {
		exponentialbackOff := float32(lastBackOff) + float32(lastBackOff)*percentage/100
		if uint64(exponentialbackOff) > uint64(maxBackOff) {
			return uint64(maxBackOff)
		}
		return uint64(exponentialbackOff)
	}
}

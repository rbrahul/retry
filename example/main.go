package main

import (
	"fmt"
	"time"

	"github.com/rbrahul/retry"
)

func main() {
	i := 0
	var err error
	// Example: It will Retry infinitely with 1 second default delay until the operation succeeds
	err = retry.Retry(func() bool {
		defer func() {
			i++
		}()
		fmt.Println("I will be called in every second until I succeed")
		return i < 10 // if i>=20 retry succeeds and will not retry further
	})

	if err == nil {
		fmt.Println("Operation suceeded")
	}

	// Example: Retry maximum 5 times with 1 second interval
	err = retry.Retry(func() bool {
		fmt.Println("I will be called in every second but maximum 5 times")
		return true
	}, 5, 1*time.Second)

	if err == retry.ErrMaximumRetryExceeded {
		fmt.Println("RETRY ERROR:", retry.ErrMaximumRetryExceeded.Error())
	}

	// Example: Will be called until 10 seconds from now in every 2 seconds
	err = retry.Retry(func() bool {
		fmt.Println("I will be called until 10 seconds from now in every 2 seconds")
		return true
	}, 10*time.Second, 2*time.Second)

	if err != nil {
		fmt.Println("RETRY ERROR:", err)
	}

	// Example: It will be called maxium 5 times with a Random interval between 2 to 5 seconds using the retry.RandomBackoff backoff function.
	err = retry.Retry(func() bool {
		fmt.Println("I will be called maxium 5 times with a random interval from 2 to 5 seconds")
		return true
	}, 5, retry.RandomBackoff(2, 5))
	if err != nil {
		fmt.Println("RETRY ERROR:", err)
	}

	// Example: It will be called maxium 5 times with an Exponential interval but maxium 10 seconds using the retry.ExponentialBackoff backoff function.
	err = retry.Retry(func() bool {
		fmt.Println("I will be called maxium 5 times with an Exponential interval but maximum delay will be 10 seconds")
		return true
	}, 5, retry.ExponentialBackoff(10))
	if err != nil {
		fmt.Println("RETRY ERROR:", err)
	}

	// Example: It will retry maxium 10 times with a Custom backoff function that returns 2 seconds delay
	err = retry.Retry(func() bool {
		fmt.Println("I will be called maxium 10 times in every 2 seconds")
		return true
	}, 10, func(_ uint64) uint64 {
		return uint64(2)
	})

	if err != nil {
		fmt.Println("RETRY ERROR:", err)
	}
}

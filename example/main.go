package main

import (
	"fmt"

	"github.com/rbrahul/retry"
)

func main() {
	//i := 0
	var err error

	// example 1: retry.Retry(retriableFn, maxRetry, delay)
	/*fmt.Println("Example 1: retry.Retry(retriableFn, maxRetry, delay)")
	retry.Retry(func() bool {
		defer func() {
			i++
		}()
		fmt.Println("Hello")
		return true
	}, 10, 1*time.Second)*/

	// example 2: retry.Retry(retriableFn, maxRetry, delayFunc)
	/*err := retry.Retry(func() bool {
		defer func() {
			i++
		}()
		fmt.Println("Hello")
		return true
	}, 10, func() uint64 {
		return 2
	})
	if err != nil {
		fmt.Println("RETRY ERROR:", err)
	}*/

	// example 3: retry.Retry(retriableFn, maxRetry{int or time.Duration }, delay{time.Duration})

	/*err = retry.Retry(func() bool {
		defer func() {
			i++
		}()
		fmt.Println("Hello")
		return true
	}, 10, 2*time.Second)

	if err != nil {
		fmt.Println("RETRY ERROR Example 3:", err)
	}*/

	/*t1 := time.Now()
	err = retry.Retry(func() bool {
		defer func() {
			i++
		}()
		fmt.Println("Hello")
		return i < 10
	}, 25, 5*time.Second)
	if err != nil {
		fmt.Println("RETRY ERROR Example 3:", err)
	}
	diff := time.Now().Unix() - t1.Unix()
	fmt.Println("DIFF:", diff)
	*/

	err = retry.Retry(func() bool {
		fmt.Println("Hello")
		return true
	}, func(_ uint64) uint64 {
		return uint64(2)
	})
	if err != nil {
		fmt.Println("RETRY ERROR Example 3:", err)
	}
}

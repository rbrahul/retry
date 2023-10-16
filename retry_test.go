package retry

import (
	"fmt"
	"testing"
	"time"
)

type callEntry struct {
	calledAt time.Time
}

type callRecorder struct {
	fn        func() bool
	callStack []callEntry
}

func NewCallRecorder(fn func() bool) *callRecorder {
	recorder := &callRecorder{
		callStack: []callEntry{},
	}
	logger := func() bool {
		recorder.callStack = append(recorder.callStack, callEntry{
			calledAt: time.Now(),
		})
		return fn()
	}
	recorder.fn = logger

	return recorder
}

func TestRetryFailsWithWrongParameters(t *testing.T) {
	err := Retry(func() bool {
		return true
	}, "wrong argument", false)

	if err != ErrMaxRetryOrRetryUntilInvalidArg {
		t.Fatalf("Error did not match. Found %+v.\nexpected %+v ", err.Error(), ErrMaxRetryOrRetryUntilInvalidArg.Error())
	}
}

func TestRetryFailsWithInvalidDelayParameter(t *testing.T) {
	err := Retry(func() bool {
		return true
	}, 1, "Invalid argument")

	if err != ErrDelayOrBackOffFuncInvalidArg {
		t.Fatalf("Error did not match. Found %+v.\nexpected %+v ", err.Error(), ErrDelayOrBackOffFuncInvalidArg.Error())
	}
}

func TestRetryCallsExactNumberOfTimesWhichWasProvidedAsParamters(t *testing.T) {
	recorder := NewCallRecorder(func() bool {
		return true
	})

	Retry(recorder.fn, 3)
	if len(recorder.callStack) != 3 {
		t.Fatalf("Expected retriable function should be executed %d times but was executed %d times", 1, len(recorder.callStack))
	}
}

func TestRetryWasExecutedWithExpectedDelay(t *testing.T) {
	recorder := NewCallRecorder(func() bool {
		return true
	})

	expectedDelay := 2

	originalJitter := addJitter

	defer func() {
		addJitter = originalJitter
	}()

	// Doesn't apply any jitter
	addJitter = func(delay time.Duration) time.Duration {
		return delay
	}

	Retry(recorder.fn, 5, time.Duration(expectedDelay)*time.Second)

	diffInSeconds := recorder.callStack[1].calledAt.Second() - recorder.callStack[0].calledAt.Second()

	if diffInSeconds != 2 {
		t.Fatalf("Expected delay is %d but found %d", expectedDelay, diffInSeconds)
	}
}

func TestRetryDeadLineWasExeecedAfterSpecifiedDeadline(t *testing.T) {
	recorder := NewCallRecorder(func() bool {
		return true
	})

	startedAt := time.Now()

	expectedDuration := 5 * time.Second
	deadLine := startedAt.Add(expectedDuration)

	err := Retry(recorder.fn, expectedDuration, 1*time.Second)

	lastCalledAt := recorder.callStack[len(recorder.callStack)-1].calledAt
	if err != ErrDeadlineExceeded {
		t.Fatalf("Expected error should be %s. But found: %s", ErrDeadlineExceeded.Error(), err.Error())
	}
	if lastCalledAt.After(deadLine) {
		t.Fatal("Retry should not call the function after the deadline is finished")
	}
}

func TestRetryHasADelaySpecifiedByBackOffFunction(t *testing.T) {
	recorder := NewCallRecorder(func() bool {
		return true
	})

	expectedDelay := 2 // seconds

	originalJitter := addJitter

	defer func() {
		addJitter = originalJitter
	}()

	// Doesn't apply any jitter
	addJitter = func(delay time.Duration) time.Duration {
		return delay
	}

	err := Retry(recorder.fn, 5, func(_ uint64) uint64 {
		return uint64(expectedDelay)
	})

	diffInSeconds := recorder.callStack[1].calledAt.Second() - recorder.callStack[0].calledAt.Second()
	if err != ErrMaximumRetryExceeded {
		t.Fatalf("Expected error should be %s. But found: %s", ErrMaximumRetryExceeded.Error(), err.Error())
	}

	if diffInSeconds != expectedDelay {
		t.Fatalf("Expected delay is %d but found %d", expectedDelay, diffInSeconds)
	}
}

func TestRetryCallsFunctionExpectedNumberOfTimes(t *testing.T) {
	expectedNumberOfCalls := 5
	i := 0
	recorder := NewCallRecorder(func() bool {
		i++
		return i < expectedNumberOfCalls
	})

	err := Retry(recorder.fn)

	if err != nil {
		t.Fatal("Should not return any error")
	}
	if len(recorder.callStack) != expectedNumberOfCalls {
		t.Fatalf("Expected retriable function should be executed %d times but was executed %d times", expectedNumberOfCalls, len(recorder.callStack))
	}
}

func TestRetryExitsWithErrorAfterSpecifiedNumberOfTryWithSpecifiedDelayInBetween(t *testing.T) {
	expectedNumberOfCalls := 10
	expectedDelay := 3
	recorder := NewCallRecorder(func() bool {
		return true
	})

	originalJitter := addJitter
	defer func() {
		addJitter = originalJitter
	}()

	// Doesn't apply any jitter
	addJitter = func(delay time.Duration) time.Duration {
		return delay
	}

	err := Retry(recorder.fn, expectedNumberOfCalls, time.Duration(expectedDelay)*time.Second)

	diffInSeconds := recorder.callStack[1].calledAt.Second() - recorder.callStack[0].calledAt.Second()

	if err != ErrMaximumRetryExceeded {
		t.Fatalf("Expected error should be %s. But found: %s", ErrMaximumRetryExceeded.Error(), err.Error())
	}

	if diffInSeconds != expectedDelay {
		t.Fatalf("Expected delay is %d but found %d", expectedDelay, diffInSeconds)
	}

	if len(recorder.callStack) != expectedNumberOfCalls {
		t.Fatalf("Expected retriable function should be executed %d times but was executed %d times", expectedNumberOfCalls, len(recorder.callStack))
	}
}

func ExampleRetry() {
	// Example: Will be called maximum 10 times in every 2 seconds
	err := Retry(func() bool {
		return true
	}, 10, 2*time.Second)

	if err != nil {
		fmt.Println("Error in Retry:", err)
	}
}

func ExampleRetry_withDeadline() {
	// Example: will be called until 10 seconds from now in every 2 seconds
	err := Retry(func() bool {
		return true
	}, 10*time.Second, 2*time.Second)

	if err != nil {
		fmt.Println("Error in Retry:", err)
	}
}

func ExampleRetry_withRandomBackoff() {
	// Example: will be called until 1 minute from now with random interval between 3 to 7 seconds
	err := Retry(func() bool {
		return true
	}, 1*time.Minute, RandomBackoff(3, 7))

	if err != nil {
		fmt.Println("Error in Retry:", err)
	}
}

func ExampleRetry_withExponentialBackoff() {
	// Example: will be called until 1 minute from now with exponential backoff of maximum 10 seconds delay
	// The exponential delay will be [1, 2, 4, 8, 10, 10, 10, 10]
	err := Retry(func() bool {
		return true
	}, 1*time.Minute, ExponentialBackoff(10))

	if err != nil {
		fmt.Println("Error in Retry:", err)
	}
}

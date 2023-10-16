![Go test workflow](https://github.com/rbrahul/retry/actions/workflows/go.yaml/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/rbrahul/retry)](https://goreportcard.com/report/github.com/rbrahul/retry)
[![codecov](https://codecov.io/gh/rbrahul/retry/branch/main/graph/badge.svg?token=XSJW5M8JC4)](https://codecov.io/gh/rbrahul/retry)
[![Go Reference](https://pkg.go.dev/badge/github.com/rbrahul/retry.svg)](https://pkg.go.dev/github.com/rbrahul/retry)

# Retry

A reliable and powerful retry library for Golang to build fault-tolerant system.

## Installation:

```bash
$ go get github.com/rbrahul/retry
```

## Usage:

```golang
import("github.com/rbrahul/retry")
```

## Examples:

#### Retry the failed operation maximum of 10 times with a 3 second delay in between.

```golang
    err := retry.Retry(func() bool {
	    operationError := doesHeavyLifting()
            if operationError != nil {
                return true // retry operation
            }
            return false // No need to retry
        }, 10, 3*time.Second)

    if err != nil {
        fmt.Error("Maxium retry exceeded")
    }
```

#### Retry failed operations with a deadline of 1 minute and with a random interval of 2 to 10 seconds.
 RandomBackoff returns a random integer value between lower and upper range. And will be considered as the delay in the number of seconds.

For Example:

RandomBackoff(2,10) will return a dealy between 2 and 10

```golang
    err := retry.Retry(func() bool {
	operationError := doesHeavyLifting()
        if operationError != nil {
            return true // retry operation
        }
		return false // No need to retry
	}, 1 * time.Minute, retry.RandomBackoff(2, 10))

    if err == retry.ErrDeadlineExceeded {
        fmt.Error("Retry deadline exceeded")
    }

```

#### Retry failed operations with a deadline of 1 minute and with ExponentialBackoff. Every retry the delay will be twice compared to the previous delay. But the maximum delay will be 10 seconds.

 ExponentialBackoff returns an exponential delay as number of seconds for each retry execution. If delay is > maxBackoff then maxBackoff will be the used as delay for next retry.

For Example:

ExponentialBackoff(10) will return the following intervals in number seconds as long as the retry is executed.

```
[1 -> 2 -> 4 -> 8 -> 10 -> 10 ...]
```
 

```golang
    err := retry.Retry(func() bool {
	operationError := doesHeavyLifting()
        if operationError != nil {
            return true // retry operation
        }
		return false // No need to retry
	}, 1 * time.Minute, retry.ExponentialBackoff(10))

    if err == retry.ErrDeadlineExceeded {
        fmt.Error("Retry deadline exceeded")
    }

```
#### Retry failed operations with a maximum 10 retries and with a custom Backoff function. Delay will be 1.5 times compared to the previous delay.

```golang
    err := retry.Retry(func() bool {
	operationError := doesHeavyLifting()
        if operationError != nil {
            return true // retry operation
        }
		return false // No need to retry
	}, 10, func(previousDelay uint64) uint64 {
		return previousDelay * 1.5
	})

    if err == retry.ErrMaximumRetryExceeded {
        fmt.Error("Maxium retry exceeded")
    }

```
#### A note related to the tiny latency:

For every delay there is a Jitter applied to avoid collisions among schedules to resolve [Thundering herd problem](https://en.wikipedia.org/wiki/Thundering_herd_problem).
 The extended delay could be between +10% to +30%.

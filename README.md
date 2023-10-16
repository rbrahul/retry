![Go test workflow](https://github.com/rbrahul/retry/actions/workflows/go.yaml/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/rbrahul/retry)](https://goreportcard.com/report/github.com/rbrahul/retry)
[![codecov](https://codecov.io/gh/rbrahul/retry/branch/main/graph/badge.svg?token=XSJW5M8JC4)](https://codecov.io/gh/rbrahul/retry)
[![Go Reference](https://pkg.go.dev/badge/github.com/rbrahul/retry.svg)](https://pkg.go.dev/github.com/rbrahul/retry)

# Retry

An essential retry-operation related library for Golang to build fault-tolerant system.

## Installation and Usage:

```golang
import("github.com/rbrahul/retry")
```

### Retry the failed operation maximum of 10 times with a 3 second delay in between.

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

### Retry failed operations with a deadline of 1 minute and with a random interval of 2 to 10 seconds.

```golang
    err := retry.Retry(func() bool {
	operationError := doesHeavyLifting()
        if operationError != nil {
            return true // retry operation
        }
		return false // No need to retry
	}, 1 * time.Minute(), retry.RandomBackoff(2, 10))

    if err == retry.ErrDeadlineExceeded {
        fmt.Error("Retry deadline exceeded")
    }

```

### Retry failed operations with a deadline of 1 minute and with ExponentialBackoff. Every retry the delay will be twice compared to the previous delay. But the maximum delay will be 10 seconds.

```golang
    err := retry.Retry(func() bool {
	operationError := doesHeavyLifting()
        if operationError != nil {
            return true // retry operation
        }
		return false // No need to retry
	}, 1 * time.Minute(), retry.ExponentialBackoff(10))

    if err == retry.ErrDeadlineExceeded {
        fmt.Error("Retry deadline exceeded")
    }

```
### Retry failed operations with a maximum 10 retries and with a custom Backoff function. Delay will be 1.5 times compared to the previous delay.

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
### A note related to the tiny latency:

For every delay there is a Jitter applied to avoid collisions among schedules to resolve [Thundering herd problem](https://en.wikipedia.org/wiki/Thundering_herd_problem).
 The extended delay could be between +10% to +30%.

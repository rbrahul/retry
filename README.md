[Go test workflow](https://github.com/rbrahul/exception/actions/workflows/go.yaml/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/rbrahul/retry)](https://goreportcard.com/report/github.com/rbrahul/retry)
[![codecov](https://codecov.io/gh/rbrahul/exception/branch/main/graph/badge.svg?token=CW54A6HWS6)](https://codecov.io/gh/rbrahul/retry)
[![Go Reference](https://pkg.go.dev/badge/github.com/rbrahul/retry.svg)](https://pkg.go.dev/github.com/rbrahul/retry)

# Retry

An essential retry-operation related library for Golang to build fault-tolerant system.

## Usages:

```go
import("github.com/rbrahul/retry")
```

### Retry the failed operation maximum of 10 times with a 3 second delay in between.

```go

    err := retry.Retry(func() bool {
		err := doesHeavyLifting()
        if err != nil {
            return true // retry operation
        }
		return false // No need to retry
	}, 10, 3*time.Second)

    if err != nil {
        fmt.Error("Maxium retry exceeded")
    }
```

### Retry failed operations with a deadline of 1 minute and with a random interval of 2 to 10 seconds.

```go
    err := retry.Retry(func() bool {
		err := doesHeavyLifting()
        if err != nil {
            return true // retry operation
        }
		return false // No need to retry
	}, 1 * time.Minute(), retry.RandomBackoff(2, 10))

    if err == retry.ErrDeadlineExceeded {
        fmt.Error("Retry deadline exceeded")
    }

```

### Retry failed operations with a deadline of 1 minute and with ExponentialBackoff. Every retry the delay will be twice compared to the previous delay. But the maximum delay will be 10 seconds.

```go
    err := retry.Retry(func() bool {
		err := doesHeavyLifting()
        if err != nil {
            return true // retry operation
        }
		return false // No need to retry
	}, 1 * time.Minute(), retry.ExponentialBackoff(10))

    if err == retry.ErrDeadlineExceeded {
        fmt.Error("Retry deadline exceeded")
    }

```
### Retry failed operations with a maximum 10 retries and with an Custom Backoff function. Delay will be 1.5 times compared to the previous delay.

```go
    err := retry.Retry(func() bool {
		err := doesHeavyLifting()
        if err != nil {
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

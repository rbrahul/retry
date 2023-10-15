# Retry

An essential retry-operation related library for Golang to build fault-tolerant system.

## Usages:

```go
import("github.com/rbrahul/retry")
```

### Retry the operation maximum 10 times in 3 seconds delay in between
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
### Retry failed operations with a maximum 10 retries and with an Custom Backoff function.

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

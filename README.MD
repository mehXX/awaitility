# awaitility

[![Go Reference](https://pkg.go.dev/badge/github.com/go-awaitility/awaitility.svg)](https://pkg.go.dev/github.com/go-awaitility/awaitility)

Package awaitility provides a simple mechanism to poll for conditions with a general timeout.

This package is inspired by https://github.com/datenhahn/golang-awaitility. However, since the original author is no longer maintaining it, 
I decided to refactor the code, fix some bugs, and integrate go.mod for easier dependency management.

I plan to support the library, so feel free to open issues or PRs.

### Await
Await - Repeatedly calls the provided until function at intervals specified by pollInterval until the function returns true or the specified atMost timeout is reached. 
If the timeout is exceeded before the condition is met, an error is returned. This function is useful for waiting on asynchronous conditions to be met.
```go
err := awaitility.Await(100 * time.Millisecond, 1000 * time.Millisecond, func() bool {
// e.g check for message in message queue
    return true
})

if err != nil {
    t.Errorf("Unexpected error during await: %s", err)
}
```

[[play (ok)](https://go.dev/play/p/Wr8mUtPxkMW)]

[[play (error on timeout)](https://go.dev/play/p/gpPLuKmCXq2)]

### AwaitCtx
AwaitCtx - Similar to Await, but also respects the provided context for cancellation. The waiting operation will terminate early if the context is canceled.
```go
ctx := context.Background()
err := awaitility.Await(ctx, 100 * time.Millisecond, 1000 * time.Millisecond, func() bool {
// e.g check for message in message queue
    return true
})

if err != nil {
    t.Errorf("Unexpected error during await: %s", err)
}
```


[[play (ok)](https://go.dev/play/p/FGvbGJ8I0mV)]

[[play (error on timeout)](https://go.dev/play/p/PjZ3fLnAvmy)]

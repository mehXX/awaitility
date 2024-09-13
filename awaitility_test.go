package awaitility_test

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/go-awaitility/awaitility"
)

func TestAwaitFalseThenOk(t *testing.T) {
	t.Parallel()
	testVar := false
	var mu sync.Mutex

	go func() {
		time.Sleep(150 * time.Millisecond)
		mu.Lock()
		testVar = true
		mu.Unlock()
	}()

	err := awaitility.Await(100*time.Millisecond, 300*time.Millisecond, func() bool {
		mu.Lock()
		defer mu.Unlock()
		return testVar
	})
	if err != nil {
		t.Errorf("Await ended with unexpected error: %s", err.Error())
	}
}

func TestAwaitFalse(t *testing.T) {
	t.Parallel()
	startTime := time.Now()

	err := awaitility.Await(10*time.Millisecond, 20*time.Millisecond, func() bool {
		return false
	})

	delta := time.Since(startTime) / time.Millisecond

	if delta > 25 {
		t.Errorf("Took too long to cancel, took %dms, expected to be under 20ms", delta)
	}

	if err == nil {
		t.Errorf("Expected await to end with error, but ended ok")
	}

	if !errors.Is(err, awaitility.ErrAwaitilityConditionTimeout) {
		t.Errorf("Expected a Timeout Error, actual error is: %s", err.Error())
	}
}

func TestPassTimeInFunc(t *testing.T) {
	t.Parallel()
	startTime := time.Now()

	err := awaitility.Await(10*time.Millisecond, 100*time.Millisecond, func() bool {
		time.Sleep(200 * time.Millisecond)
		return true
	})

	delta := time.Since(startTime) / time.Millisecond

	if delta > 105 {
		t.Errorf("Took too long to cancel, took %dms, expected to be under 100ms", delta)
	}

	if err == nil {
		t.Errorf("Expected await to cancel")
	}
}

func TestAwaitLimits(t *testing.T) {
	t.Parallel()
	err := awaitility.Await(0, 10*time.Millisecond, func() bool {
		return false
	})

	if err == nil {
		t.Errorf("Expected error when 0 poll interval but none received")
	}
	if !errors.Is(err, awaitility.ErrPollIntervalValidation) {
		t.Errorf("Expected message '%s' but got '%s'", awaitility.ErrPollIntervalValidation, err.Error())
	}

	err = awaitility.Await(10*time.Millisecond, 0, func() bool {
		return false
	})

	if err == nil {
		t.Errorf("Expected error when 0 poll interval but none received")
	}
	if !errors.Is(err, awaitility.ErrAtMostValidation) {
		t.Errorf("Expected message '%s' but got '%s'", awaitility.ErrAtMostValidation, err.Error())
	}

	err = awaitility.Await(20*time.Millisecond, 10*time.Millisecond, func() bool {
		return false
	})

	if err == nil {
		t.Errorf("Expected error when 0 poll interval but none received")
	}
	if !errors.Is(err, awaitility.ErrPollIntervalMustBeSmallerThanAtMost) {
		t.Errorf("Expected message '%s' but got '%s'", awaitility.ErrPollIntervalMustBeSmallerThanAtMost, err.Error())
	}
}

func TestAwaitCtxWithContextTimeout(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	atMost := 200 * time.Millisecond
	pollInterval := 10 * time.Millisecond

	startTime := time.Now()

	err := awaitility.AwaitCtx(ctx, pollInterval, atMost, func() bool {
		return false
	})

	delta := time.Since(startTime)

	if delta >= atMost {
		t.Errorf("Expected to timeout due to context, but took the full atMost duration")
	}

	if err == nil {
		t.Errorf("Expected an error due to context timeout, but no error was returned")
	}

	if !errors.Is(err, context.DeadlineExceeded) {
		t.Errorf("Expected context.DeadlineExceeded error, but got: %v", err)
	}

	if !errors.Is(err, awaitility.ErrAwaitilityConditionTimeout) {
		t.Errorf("Expected awaitility.ErrAwaitilityConditionTimeout error, but got: %v", err)
	}
}

func TestAwaitPanicInUntilFunction(t *testing.T) {
	t.Parallel()

	err := awaitility.Await(10*time.Millisecond, 100*time.Millisecond, func() bool {
		panic("test panic")
	})

	if err == nil {
		t.Errorf("Expected an error due to panic, but no error was returned")
	}

	if !errors.Is(err, awaitility.ErrPanicInUntilFunction) {
		t.Errorf("Expected ErrPanicInUntilFunction error, but got: %v", err)
	}
}

func TestAwaitImmediateSuccess(t *testing.T) {
	t.Parallel()

	startTime := time.Now()
	err := awaitility.Await(100*time.Millisecond, 1*time.Second, func() bool {
		return true
	})

	delta := time.Since(startTime)

	if err != nil {
		t.Errorf("Unexpected error: %s", err.Error())
	}

	if delta >= 100*time.Millisecond {
		t.Errorf("Expected immediate return, but took %v", delta)
	}
}

func TestAwaitCtxCancelledBeforeExecution(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := awaitility.AwaitCtx(ctx, 100*time.Millisecond, 1*time.Second, func() bool {
		return false
	})

	if err == nil {
		t.Errorf("Expected an error due to context cancellation, but no error was returned")
	}

	if !errors.Is(err, context.Canceled) {
		t.Errorf("Expected context.Canceled error, but got: %v", err)
	}
}

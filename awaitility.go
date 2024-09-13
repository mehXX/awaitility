package awaitility

import (
	"context"
	"errors"
	"fmt"
	"time"

	"golang.org/x/sync/errgroup"
)

var (
	// ErrAwaitilityConditionTimeout is returned when the condition does not become true
	// within the specified timeout.
	ErrAwaitilityConditionTimeout = errors.New("await condition did not return true, limit reached")

	// ErrPollIntervalValidation is returned when the pollInterval is set to 0 or a negative value.
	ErrPollIntervalValidation = errors.New("pollInterval cannot be 0 or below")

	// ErrAtMostValidation is returned when the atMost timeout is set to 0 or a negative value.
	ErrAtMostValidation = errors.New("atMost timeout cannot be 0 or below")

	// ErrPollIntervalMustBeSmallerThanAtMost is returned when the pollInterval is greater than the atMost timeout.
	ErrPollIntervalMustBeSmallerThanAtMost = errors.New("pollInterval must be smaller than atMost timeout")

	// ErrPanicInUntilFunction is returned when a panic occurs in the until function.
	ErrPanicInUntilFunction = errors.New("panic occurred in until function")
)

// Await calls the "until" function initially and then at intervals specified by "pollInterval"
// until the function returns true or the total time spent exceeds the "atMost" limit.
// If the time limit is exceeded, an error is returned. Note that the "until" function will
// continue to run in the background even if the Await function exits after the timeout.
// Play: https://go.dev/play/p/Wr8mUtPxkMW
func Await(pollInterval time.Duration, atMost time.Duration, until func() bool) error {
	ctx := context.Background()
	return AwaitCtx(ctx, pollInterval, atMost, until)
}

// AwaitCtx is the same as Await, but it also respects the provided context for cancellation.
// Play: https://go.dev/play/p/FGvbGJ8I0mV
func AwaitCtx(ctx context.Context, pollInterval time.Duration, atMost time.Duration, until func() bool) error {
	if err := validateAwaitArgs(pollInterval, atMost); err != nil {
		return fmt.Errorf("validation: %w", err)
	}

	ctx, cancel := context.WithTimeout(ctx, atMost)
	defer cancel()

	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()

	resultCh := make(chan bool, 1)
	errCh := make(chan error, 1)

	eg, ctx := errgroup.WithContext(ctx)
	eg.SetLimit(1)

	runUntil := func() error {
		defer func() {
			if err := recover(); err != nil {
				resultCh <- false
				errCh <- fmt.Errorf("%w: %v", ErrPanicInUntilFunction, err)
			}
		}()
		resultCh <- until()
		return nil
	}

	eg.TryGo(runUntil)

	for {
		select {
		case <-ticker.C:
			eg.TryGo(runUntil)
		case result := <-resultCh:
			if result {
				return nil
			}
		case err := <-errCh:
			return err
		case <-ctx.Done():
			return fmt.Errorf("%w: %w", ErrAwaitilityConditionTimeout, ctx.Err())
		}
	}
}

func validateAwaitArgs(pollInterval, atMost time.Duration) error {
	if pollInterval <= 0 {
		return fmt.Errorf("%w, got: %v", ErrPollIntervalValidation, pollInterval)
	}

	if atMost <= 0 {
		return fmt.Errorf("%w, got: %v", ErrAtMostValidation, atMost)
	}

	if pollInterval > atMost {
		return fmt.Errorf("%w, got: pollInterval=%v, atMost=%v", ErrPollIntervalMustBeSmallerThanAtMost, pollInterval, atMost)
	}

	return nil
}

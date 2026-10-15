package worker

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPool_Run_Sync(t *testing.T) {
	p := NewPool[int](1, 1, time.Minute)

	ctx := context.Background()
	worker := func(ctx context.Context, i int) {
		time.Sleep(time.Millisecond)
	}

	err := p.Run(ctx, worker)

	require.NoError(t, err, "unexpected error for status not_initialized")

	err = p.Run(ctx, worker)

	assert.ErrorIs(t, err, ErrAlreadyRunning, "incorrect error for status running")

	err = p.Close()

	require.NoError(t, err, "unexpected error on Close")

	err = p.Run(ctx, worker)

	assert.NoError(t, err, "unexpected error for status closed")
}

func TestPool_AssignJob_Sync(t *testing.T) {
	p := NewPool[int](1, 1, time.Minute)

	ctx := context.Background()

	assert.ErrorIs(t, p.AssignJob(1), ErrNotInitialized, "incorrect error for status not_initialized")

	require.NoError(t, p.Run(ctx, func(ctx context.Context, i int) {
		time.Sleep(time.Millisecond * time.Duration(i))
	}), "unexpected error on Run")

	assert.NoError(t, p.AssignJob(1), "unexpected error on AssignJob for status running")

	p.AssignJob(100)                                // fill the queue
	p.AssignJob(100)                                // fill the queue
	assert.ErrorIs(t, p.AssignJob(1), ErrQueueFull) // backpressure check

	require.NoError(t, p.Close(), "unexpected error on Close")

	time.Sleep(10 * time.Millisecond)

	assert.ErrorIs(t, p.AssignJob(1), ErrQueueClosed, "incorrect error for status closed")
}

func TestPool_Close_Sync(t *testing.T) {
	p := NewPool[int](1, 1, time.Minute)

	ctx := context.Background()

	assert.ErrorIs(t, p.Close(), ErrNotInitialized, "incorrect error for status not_initialized")

	require.NoError(t, p.Run(ctx, func(ctx context.Context, i int) {
		time.Sleep(time.Millisecond)
	}), "unexpected error on Run")

	require.NoError(t, p.Close(), "unexpected error for status running")
	assert.ErrorIs(t, p.Close(), ErrQueueClosed, "incorrect error for status closed")
}

func TestPool_AssignJob_Close_Race(t *testing.T) {
	const (
		iterations = 100
		assigners  = 10
		jobs       = 100_000 // avoid infinite loop in assigners if Close fails
	)

	for range iterations {
		p := NewPool[int](2, 4, time.Minute)

		var wg sync.WaitGroup
		var counter atomic.Int32

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		require.NoError(t, p.Run(ctx, func(ctx context.Context, i int) {
			time.Sleep(time.Millisecond)
		}), "unexpected error on Run")

		for range assigners {
			wg.Add(1)

			go func() {
				defer wg.Done()

				defer func() {
					if r := recover(); r != nil {
						t.Errorf("panic on AssignJob: %v", r)
						counter.Add(1)
					}
				}()

				for range jobs {
					err := p.AssignJob(1)

					if err == ErrQueueClosed {
						return
					}

					if err != nil && err != ErrQueueFull {
						t.Errorf("unexpected error on AssignJob: %v", err)
						return
					}
				}
			}()
		}

		wg.Add(1)
		go func() {
			defer wg.Done()

			time.Sleep(time.Millisecond)

			assert.NoError(t, p.Close())
		}()

		wg.Wait()
		require.Zero(t, counter.Load(), "AssignJob panic detected")
	}
}

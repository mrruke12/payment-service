package worker

import (
	"context"
	"fmt"
	"log/slog"
	"runtime/debug"
	"sync"
	"time"
)

type PoolStatus string

const (
	PoolStatusNotInitialized PoolStatus = "not_initialized"
	PoolStatusRunning        PoolStatus = "running"
	PoolStatusClosed         PoolStatus = "closed"
)

// A thread-safe implementation of Worker Pool to run workers
type Pool[T any] struct {
	status  PoolStatus
	mu      sync.Mutex
	wg      sync.WaitGroup
	jobs    chan T
	workers int
	buffer  int
	timeout time.Duration
}

func NewPool[T any](workers, buffer int, timeout time.Duration) *Pool[T] {
	return &Pool[T]{
		status:  PoolStatusNotInitialized,
		jobs:    nil,
		workers: workers,
		buffer:  buffer,
		timeout: timeout,
	}
}

// Run workers
func (p *Pool[T]) Run(ctx context.Context, worker func(context.Context, T)) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.status == PoolStatusRunning {
		return ErrAlreadyRunning
	}

	p.status = PoolStatusRunning

	p.jobs = make(chan T, p.buffer)
	jobs := p.jobs // capture the channel to avoid concurrent nil referrence when Close called before goroutines are ran

	for range p.workers {
		fmt.Println(1)
		p.wg.Add(1)

		go func() {
			defer p.wg.Done()

			for {
				select {
				case v, ok := <-jobs:
					if !ok {
						return
					}
					func() {
						defer func() {
							if r := recover(); r != nil {
								slog.Error("worker panicked",
									"recover", r,
									"stack", string(debug.Stack()))
							}
						}()

						c, cancel := context.WithTimeout(ctx, p.timeout)
						defer cancel()

						worker(c, v)
					}()
				case <-ctx.Done():
					return
				}
			}
		}()
	}

	return nil
}

// Assign a job to the queue. Implements at-most-once delivery.
//
// Returns worker.ErrQueueFull if the buffer is full, hence retries are required.
func (p *Pool[T]) AssignJob(job T) error {
	p.mu.Lock()
	defer p.mu.Unlock() // hold the lock until the job is enqueued to prevent panic if the queue is closed before

	switch p.status {
	case PoolStatusNotInitialized:
		return ErrNotInitialized
	case PoolStatusClosed:
		return ErrQueueClosed
	default:
		break
	}

	select {
	case p.jobs <- job:
		return nil
	default:
		return ErrQueueFull
	}
}

// Close the queue and wait for workers to complete
func (p *Pool[T]) Close() error {
	p.mu.Lock()

	switch p.status {
	case PoolStatusNotInitialized:
		p.mu.Unlock()
		return ErrNotInitialized
	case PoolStatusClosed:
		p.mu.Unlock()
		return ErrQueueClosed
	default:
		break
	}

	p.status = PoolStatusClosed

	p.mu.Unlock()

	close(p.jobs)

	p.jobs = nil

	p.wg.Wait()

	return nil
}

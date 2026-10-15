package worker

import "errors"

var (
	ErrQueueFull      error = errors.New("queue is full")
	ErrAlreadyRunning error = errors.New("workers are already running")
	ErrQueueClosed    error = errors.New("queue is closed")
	ErrNotInitialized error = errors.New("workers are not initialized")
)

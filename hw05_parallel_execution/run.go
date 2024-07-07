package hw05parallelexecution

import (
	"errors"
)

var (
	ErrErrorsLimitExceeded = errors.New("errors limit exceeded")
	ErrErrorsNoWorkers     = errors.New("no worker")
)

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	if n <= 0 {
		return ErrErrorsNoWorkers
	}
	threadPoolExecutor := newPool(tasks, n, int64(m))

	return threadPoolExecutor.run()
}

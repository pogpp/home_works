package hw05parallelexecution

import (
	"sync"
)

type worker struct {
	taskCh chan Task
}

func newTasksSolver(chTask chan Task) *worker {
	return &worker{
		taskCh: chTask,
	}
}

func (w *worker) Start(wg *sync.WaitGroup, limiter *limiter) {
	wg.Add(1)

	go func() {
		defer wg.Done()
		for task := range w.taskCh {
			if task() != nil {
				limiter.increment()
				if limiter.isLimitExceeded() {
					return
				}
			}
		}
	}()
}

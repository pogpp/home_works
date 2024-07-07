package hw05parallelexecution

import (
	"sync"
	"sync/atomic"
)

const chanSize = 1000

type limiter struct {
	count int64
	limit int64
}

type pool struct {
	tasks         []Task
	routinesCount int
	collector     chan Task
	wg            sync.WaitGroup
	limiter       *limiter
}

func newLimiter(limit int64) *limiter {
	return &limiter{
		limit: limit,
	}
}

func newPool(tasks []Task, rCount int, maxErrorCount int64) *pool {
	return &pool{
		tasks:         tasks,
		routinesCount: rCount,
		collector:     make(chan Task, chanSize),
		limiter:       newLimiter(maxErrorCount),
	}
}

func (l *limiter) increment() {
	atomic.AddInt64(&l.count, 1)
}

func (l *limiter) isLimitExceeded() bool {
	return atomic.LoadInt64(&l.count) >= atomic.LoadInt64(&l.limit)
}

func (p *pool) run() error {
	for i := 0; i < p.routinesCount; i++ {
		w := newTasksSolver(p.collector)
		w.Start(&p.wg, p.limiter)
	}

	for _, task := range p.tasks {
		p.collector <- task
	}
	close(p.collector)

	p.wg.Wait()

	if p.limiter.isLimitExceeded() {
		return ErrErrorsLimitExceeded
	}

	return nil
}

package worker

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/pablo-banker/taskforge/internal/queue"
	"github.com/pablo-banker/taskforge/internal/store"
)

type Pool struct {
	store       store.Store
	queue       *queue.Manager
	workers     int
	taskTimeout time.Duration
	handler     Handler
}

func NewPool(
	store store.Store,
	queue *queue.Manager,
	workers int,
	taskTimeout time.Duration,
	handler Handler,
) *Pool {
	if workers <= 0 {
		workers = 1
	}

	if taskTimeout <= 0 {
		taskTimeout = 5 * time.Second
	}

	return &Pool{
		store:       store,
		queue:       queue,
		workers:     workers,
		taskTimeout: taskTimeout,
		handler:     handler,
	}
}

func (p *Pool) Start(ctx context.Context) {
	var wg sync.WaitGroup

	for i := 1; i <= p.workers; i++ {
		wg.Add(1)

		go func(workerID int) {
			defer wg.Done()
			p.runWorker(ctx, workerID)
		}(i)
	}

	<-ctx.Done()

	wg.Wait()
}

func (p *Pool) runWorker(ctx context.Context, workerID int) {
	for {
		select {
		case job := <-p.queue.Jobs():
			p.process(ctx, workerID, job)

		case <-ctx.Done():
			fmt.Printf("[worker-%d] stopped\n", workerID)
			return
		}
	}
}

func (p *Pool) process(ctx context.Context, workerID int, job queue.Job) {
	t, started, err := p.store.StartTask(job.TaskID)
	if err != nil {
		fmt.Printf("[worker-%d] failed to start task=%s error=%v\n", workerID, job.TaskID, err)
		return
	}

	if !started {
		return
	}

	fmt.Printf("[worker-%d] running task=%s attempt=%d/%d\n", workerID, t.ID, t.Attempts, t.MaxAttempts)

	taskCtx, cancel := context.WithTimeout(ctx, p.taskTimeout)
	defer cancel()

	if err := p.handler.Handle(taskCtx, t); err != nil {
		nextRunAt := time.Now().Add(backoff(t.Attempts))

		updated, updateErr := p.store.FailTask(t.ID, err.Error(), nextRunAt)
		if updateErr != nil {
			fmt.Printf("[worker-%d] failed to update failed task=%s error=%v\n", workerID, t.ID, updateErr)
			return
		}

		fmt.Printf(
			"[worker-%d] task=%s failed status=%s error=%q\n",
			workerID,
			updated.ID,
			updated.Status,
			updated.LastError,
		)

		return
	}

	updated, err := p.store.CompleteTask(t.ID)
	if err != nil {
		fmt.Printf("[worker-%d] failed to complete task=%s error=%v\n", workerID, t.ID, err)
		return
	}

	fmt.Printf("[worker-%d] completed task=%s status=%s\n", workerID, updated.ID, updated.Status)
}

func backoff(attempt int) time.Duration {
	if attempt <= 0 {
		attempt = 1
	}

	return time.Duration(attempt*attempt) * time.Second
}

package scheduler

import (
	"context"
	"fmt"
	"time"

	"github.com/pablo-banker/taskforge/internal/queue"
	"github.com/pablo-banker/taskforge/internal/store"
)

type Scheduler struct {
	store    store.Store
	queue    *queue.Manager
	interval time.Duration
}

func New(store store.Store, queue *queue.Manager, interval time.Duration) *Scheduler {
	if interval <= 0 {
		interval = time.Second
	}

	return &Scheduler{
		store:    store,
		queue:    queue,
		interval: interval,
	}
}

func (s *Scheduler) Start(ctx context.Context) {
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	s.dispatchReady(ctx)

	for {
		select {
		case <-ticker.C:
			s.dispatchReady(ctx)

		case <-ctx.Done():
			fmt.Println("[scheduler] stopped")
			return
		}
	}
}

func (s *Scheduler) dispatchReady(ctx context.Context) {
	readyTasks := s.store.ReadyScheduled(time.Now())

	for _, t := range readyTasks {
		if err := s.queue.Publish(ctx, t.Queue, t.ID); err != nil {
			fmt.Printf("[scheduler] failed to publish task=%s error=%v\n", t.ID, err)
			continue
		}

		fmt.Printf("[scheduler] queued task=%s queue=%s\n", t.ID, t.Queue)
	}
}

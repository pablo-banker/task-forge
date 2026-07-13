package queue

import "context"

type Job struct {
	Queue  string
	TaskID string
}

type Manager struct {
	jobs chan Job
}

func NewManager(buffer int) *Manager {
	return &Manager{
		jobs: make(chan Job, buffer),
	}
}

func (m *Manager) Publish(ctx context.Context, queueName string, taskID string) error {
	select {
	case m.jobs <- Job{
		Queue:  queueName,
		TaskID: taskID,
	}:
		return nil

	case <-ctx.Done():
		return ctx.Err()
	}
}

func (m *Manager) Jobs() <-chan Job {
	return m.jobs
}

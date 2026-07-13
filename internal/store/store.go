package store

import (
	"context"
	"time"

	"github.com/pablo-banker/taskforge/internal/task"
)

type Store interface {
	Create(ctx context.Context, t task.Task) (task.Task, error)
	Get(ctx context.Context, id string) (task.Task, error)
	List(ctx context.Context, queueName string, status task.Status) ([]task.Task, error)

	ReadyScheduled(now time.Time) []task.Task

	StartTask(id string) (task.Task, bool, error)
	CompleteTask(id string) (task.Task, error)
	FailTask(id string, message string, nextRunAt time.Time) (task.Task, error)
	CancelTask(id string) (task.Task, error)

	Subscribe(id string) (<-chan task.Task, func())
}

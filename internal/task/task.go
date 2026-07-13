package task

import "time"

type Status string

const (
	StatusPending    Status = "pending"
	StatusScheduled  Status = "scheduled"
	StatusQueued     Status = "queued"
	StatusRunning    Status = "running"
	StatusCompleted  Status = "completed"
	StatusFailed     Status = "failed"
	StatusCancelled  Status = "cancelled"
	StatusDeadLetter Status = "dead_letter"
)

type Priority int

const (
	PriorityLow Priority = iota + 1
	PriorityNormal
	PriorityHigh
)

type Task struct {
	ID          string
	Queue       string
	Name        string
	Payload     string
	Status      Status
	Priority    Priority
	Attempts    int
	MaxAttempts int
	RunAt       time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
	LastError   string
}

func (t Task) IsTerminal() bool {
	return t.Status == StatusCompleted ||
		t.Status == StatusCancelled ||
		t.Status == StatusDeadLetter
}

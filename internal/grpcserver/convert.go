package grpcserver

import (
	"time"

	"github.com/pablo-banker/taskforge/internal/task"
	pb "github.com/pablo-banker/taskforge/proto/taskforge/v1"
)

func taskToProto(t task.Task) *pb.Task {
	return &pb.Task{
		Id:            t.ID,
		Queue:         t.Queue,
		Name:          t.Name,
		Payload:       t.Payload,
		Status:        statusToProto(t.Status),
		Priority:      priorityToProto(t.Priority),
		Attempts:      int32(t.Attempts),
		MaxAttempts:   int32(t.MaxAttempts),
		RunAtUnix:     t.RunAt.Unix(),
		CreatedAtUnix: t.CreatedAt.Unix(),
		UpdatedAtUnix: t.UpdatedAt.Unix(),
		LastError:     t.LastError,
	}
}

func statusToProto(status task.Status) pb.TaskStatus {
	switch status {
	case task.StatusPending:
		return pb.TaskStatus_TASK_STATUS_PENDING
	case task.StatusScheduled:
		return pb.TaskStatus_TASK_STATUS_SCHEDULED
	case task.StatusQueued:
		return pb.TaskStatus_TASK_STATUS_QUEUED
	case task.StatusRunning:
		return pb.TaskStatus_TASK_STATUS_RUNNING
	case task.StatusCompleted:
		return pb.TaskStatus_TASK_STATUS_COMPLETED
	case task.StatusFailed:
		return pb.TaskStatus_TASK_STATUS_FAILED
	case task.StatusCancelled:
		return pb.TaskStatus_TASK_STATUS_CANCELLED
	case task.StatusDeadLetter:
		return pb.TaskStatus_TASK_STATUS_DEAD_LETTER
	default:
		return pb.TaskStatus_TASK_STATUS_UNSPECIFIED
	}
}

func statusFromProto(status pb.TaskStatus) task.Status {
	switch status {
	case pb.TaskStatus_TASK_STATUS_PENDING:
		return task.StatusPending
	case pb.TaskStatus_TASK_STATUS_SCHEDULED:
		return task.StatusScheduled
	case pb.TaskStatus_TASK_STATUS_QUEUED:
		return task.StatusQueued
	case pb.TaskStatus_TASK_STATUS_RUNNING:
		return task.StatusRunning
	case pb.TaskStatus_TASK_STATUS_COMPLETED:
		return task.StatusCompleted
	case pb.TaskStatus_TASK_STATUS_FAILED:
		return task.StatusFailed
	case pb.TaskStatus_TASK_STATUS_CANCELLED:
		return task.StatusCancelled
	case pb.TaskStatus_TASK_STATUS_DEAD_LETTER:
		return task.StatusDeadLetter
	default:
		return ""
	}
}

func priorityToProto(priority task.Priority) pb.TaskPriority {
	switch priority {
	case task.PriorityLow:
		return pb.TaskPriority_TASK_PRIORITY_LOW
	case task.PriorityNormal:
		return pb.TaskPriority_TASK_PRIORITY_NORMAL
	case task.PriorityHigh:
		return pb.TaskPriority_TASK_PRIORITY_HIGH
	default:
		return pb.TaskPriority_TASK_PRIORITY_UNSPECIFIED
	}
}

func priorityFromProto(priority pb.TaskPriority) task.Priority {
	switch priority {
	case pb.TaskPriority_TASK_PRIORITY_LOW:
		return task.PriorityLow
	case pb.TaskPriority_TASK_PRIORITY_HIGH:
		return task.PriorityHigh
	case pb.TaskPriority_TASK_PRIORITY_NORMAL:
		return task.PriorityNormal
	default:
		return task.PriorityNormal
	}
}

func unixOrNow(value int64) time.Time {
	if value <= 0 {
		return time.Now()
	}

	return time.Unix(value, 0)
}

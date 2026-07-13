package grpcserver

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/pablo-banker/taskforge/internal/queue"
	"github.com/pablo-banker/taskforge/internal/store"
	"github.com/pablo-banker/taskforge/internal/task"
	pb "github.com/pablo-banker/taskforge/proto/taskforge/v1"
	"google.golang.org/grpc/codes"
	grpcstatus "google.golang.org/grpc/status"
)

type Service struct {
	pb.UnimplementedTaskServiceServer

	store store.Store
	queue *queue.Manager
}

func NewService(store store.Store, queue *queue.Manager) *Service {
	return &Service{
		store: store,
		queue: queue,
	}
}

func (s *Service) CreateTask(ctx context.Context, req *pb.CreateTaskRequest) (*pb.CreateTaskResponse, error) {
	queueName := strings.TrimSpace(req.GetQueue())
	if queueName == "" {
		queueName = "default"
	}

	name := strings.TrimSpace(req.GetName())
	if name == "" {
		return nil, grpcstatus.Error(codes.InvalidArgument, "task name is required")
	}

	payload := req.GetPayload()

	maxAttempts := int(req.GetMaxAttempts())
	if maxAttempts <= 0 {
		maxAttempts = 3
	}

	now := time.Now()
	runAt := unixOrNow(req.GetRunAtUnix())

	status := task.StatusPending
	if runAt.After(now) {
		status = task.StatusScheduled
	}

	newTask := task.Task{
		ID:          uuid.NewString(),
		Queue:       queueName,
		Name:        name,
		Payload:     payload,
		Status:      status,
		Priority:    priorityFromProto(req.GetPriority()),
		Attempts:    0,
		MaxAttempts: maxAttempts,
		RunAt:       runAt,
		CreatedAt:   now,
		UpdatedAt:   now,
		LastError:   "",
	}

	created, err := s.store.Create(ctx, newTask)
	if err != nil {
		return nil, grpcstatus.Error(codes.Internal, err.Error())
	}

	if created.Status == task.StatusPending {
		if err := s.queue.Publish(ctx, created.Queue, created.ID); err != nil {
			return nil, grpcstatus.Error(codes.Internal, err.Error())
		}
	}

	return &pb.CreateTaskResponse{
		Task: taskToProto(created),
	}, nil
}

func (s *Service) GetTask(ctx context.Context, req *pb.GetTaskRequest) (*pb.GetTaskResponse, error) {
	id := strings.TrimSpace(req.GetId())
	if id == "" {
		return nil, grpcstatus.Error(codes.InvalidArgument, "task id is required")
	}

	t, err := s.store.Get(ctx, id)
	if err != nil {
		return nil, mapStoreError(err)
	}

	return &pb.GetTaskResponse{
		Task: taskToProto(t),
	}, nil
}

func (s *Service) ListTasks(ctx context.Context, req *pb.ListTasksRequest) (*pb.ListTasksResponse, error) {
	queueName := strings.TrimSpace(req.GetQueue())
	status := statusFromProto(req.GetStatus())

	tasks, err := s.store.List(ctx, queueName, status)
	if err != nil {
		return nil, grpcstatus.Error(codes.Internal, err.Error())
	}

	response := &pb.ListTasksResponse{
		Tasks: make([]*pb.Task, 0, len(tasks)),
	}

	for _, t := range tasks {
		response.Tasks = append(response.Tasks, taskToProto(t))
	}

	return response, nil
}

func (s *Service) CancelTask(ctx context.Context, req *pb.CancelTaskRequest) (*pb.CancelTaskResponse, error) {
	id := strings.TrimSpace(req.GetId())
	if id == "" {
		return nil, grpcstatus.Error(codes.InvalidArgument, "task id is required")
	}

	t, err := s.store.CancelTask(id)
	if err != nil {
		return nil, mapStoreError(err)
	}

	return &pb.CancelTaskResponse{
		Task: taskToProto(t),
	}, nil
}

func (s *Service) WatchTask(req *pb.WatchTaskRequest, stream pb.TaskService_WatchTaskServer) error {
	id := strings.TrimSpace(req.GetId())
	if id == "" {
		return grpcstatus.Error(codes.InvalidArgument, "task id is required")
	}

	current, err := s.store.Get(stream.Context(), id)
	if err != nil {
		return mapStoreError(err)
	}

	if err := stream.Send(&pb.WatchTaskResponse{
		Task: taskToProto(current),
	}); err != nil {
		return err
	}

	updates, unsubscribe := s.store.Subscribe(id)
	defer unsubscribe()

	for {
		select {
		case t := <-updates:
			if err := stream.Send(&pb.WatchTaskResponse{
				Task: taskToProto(t),
			}); err != nil {
				return err
			}

		case <-stream.Context().Done():
			return nil
		}
	}
}

func mapStoreError(err error) error {
	switch {
	case errors.Is(err, store.ErrTaskNotFound):
		return grpcstatus.Error(codes.NotFound, err.Error())

	case errors.Is(err, store.ErrTaskAlreadyExists):
		return grpcstatus.Error(codes.AlreadyExists, err.Error())

	case errors.Is(err, store.ErrTaskNotCancellable):
		return grpcstatus.Error(codes.FailedPrecondition, err.Error())

	default:
		return grpcstatus.Error(codes.Internal, err.Error())
	}
}

package app

import (
	"context"
	"fmt"
	"net"

	"github.com/pablo-banker/taskforge/internal/cliui"
	"github.com/pablo-banker/taskforge/internal/config"
	"github.com/pablo-banker/taskforge/internal/grpcserver"
	"github.com/pablo-banker/taskforge/internal/queue"
	"github.com/pablo-banker/taskforge/internal/scheduler"
	"github.com/pablo-banker/taskforge/internal/store"
	"github.com/pablo-banker/taskforge/internal/worker"
	pb "github.com/pablo-banker/taskforge/proto/taskforge/v1"
	"google.golang.org/grpc"
)

type App struct {
	cfg config.Config
}

func New(cfg config.Config) *App {
	return &App{
		cfg: cfg,
	}
}

func (a *App) Run(ctx context.Context) error {
	listener, err := net.Listen("tcp", a.cfg.Addr)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", a.cfg.Addr, err)
	}

	memoryStore := store.NewMemoryStore()
	queueManager := queue.NewManager(a.cfg.QueueBuffer)

	taskScheduler := scheduler.New(
		memoryStore,
		queueManager,
		a.cfg.SchedulerInterval,
	)

	workerPool := worker.NewPool(
		memoryStore,
		queueManager,
		a.cfg.Workers,
		a.cfg.TaskTimeout,
		worker.SimulatedHandler{},
	)

	grpcServer := grpc.NewServer()

	taskService := grpcserver.NewService(memoryStore, queueManager)
	pb.RegisterTaskServiceServer(grpcServer, taskService)

	cliui.PrintServerBanner(cliui.ServerBannerConfig{
		Addr:              a.cfg.Addr,
		Workers:           a.cfg.Workers,
		TaskTimeout:       a.cfg.TaskTimeout,
		SchedulerInterval: a.cfg.SchedulerInterval,
		QueueBuffer:       a.cfg.QueueBuffer,
		Store:             a.cfg.Store,
	})

	componentCtx, cancelComponents := context.WithCancel(ctx)
	defer cancelComponents()

	go taskScheduler.Start(componentCtx)
	go workerPool.Start(componentCtx)

	serverErrors := make(chan error, 1)

	go func() {
		serverErrors <- grpcServer.Serve(listener)
	}()

	select {
	case <-ctx.Done():
		cliui.PrintWarning("Shutting down TaskForge server...")
		grpcServer.GracefulStop()
		cancelComponents()
		cliui.PrintSuccess("Server stopped gracefully.")
		return nil

	case err := <-serverErrors:
		cancelComponents()

		if err != nil {
			return fmt.Errorf("gRPC server stopped unexpectedly: %w", err)
		}

		return nil
	}
}

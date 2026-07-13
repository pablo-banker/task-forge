package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/pablo-banker/taskforge/internal/cliui"
	"github.com/pablo-banker/taskforge/internal/tui"
	pb "github.com/pablo-banker/taskforge/proto/taskforge/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "create":
		createCommand(os.Args[2:])

	case "get":
		getCommand(os.Args[2:])

	case "list":
		listCommand(os.Args[2:])

	case "cancel":
		cancelCommand(os.Args[2:])

	case "watch":
		watchCommand(os.Args[2:])

	case "tui":
		tuiCommand(os.Args[2:])

	default:
		usage()
		os.Exit(1)
	}
}

func usage() {
	fmt.Println(`TaskForge CLI

Usage:
  taskforge create --name welcome --payload "hello"
  taskforge create --queue emails --name welcome --payload "send email" --delay 10s
  taskforge create --name broken --payload "fail" --max-attempts 3
  taskforge get --id <task-id>
  taskforge list
  taskforge list --queue emails
  taskforge list --status completed
  taskforge cancel --id <task-id>
  taskforge watch --id <task-id>
  taskforge tui
`)
}

func tuiCommand(args []string) {
	fs := flag.NewFlagSet("tui", flag.ExitOnError)

	addr := fs.String("addr", "localhost:50051", "server address")

	_ = fs.Parse(args)

	if err := tui.Run(*addr); err != nil {
		log.Fatalf("tui failed: %v", err)
	}
}

func createCommand(args []string) {
	fs := flag.NewFlagSet("create", flag.ExitOnError)

	addr := fs.String("addr", "localhost:50051", "server address")
	queueName := fs.String("queue", "default", "queue name")
	name := fs.String("name", "", "task name")
	payload := fs.String("payload", "", "task payload")
	priority := fs.String("priority", "normal", "priority: low, normal, high")
	maxAttempts := fs.Int("max-attempts", 3, "max attempts")
	delay := fs.Duration("delay", 0, "delay before running task")

	_ = fs.Parse(args)

	client, conn := mustClient(*addr)
	defer conn.Close()

	runAtUnix := int64(0)
	if *delay > 0 {
		runAtUnix = time.Now().Add(*delay).Unix()
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := client.CreateTask(ctx, &pb.CreateTaskRequest{
		Queue:       *queueName,
		Name:        *name,
		Payload:     *payload,
		Priority:    parsePriority(*priority),
		MaxAttempts: int32(*maxAttempts),
		RunAtUnix:   runAtUnix,
	})
	if err != nil {
		log.Fatalf("create task failed: %v", err)
	}

	cliui.PrintTask(res.GetTask())
}

func getCommand(args []string) {
	fs := flag.NewFlagSet("get", flag.ExitOnError)

	addr := fs.String("addr", "localhost:50051", "server address")
	id := fs.String("id", "", "task id")

	_ = fs.Parse(args)

	client, conn := mustClient(*addr)
	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := client.GetTask(ctx, &pb.GetTaskRequest{
		Id: *id,
	})
	if err != nil {
		log.Fatalf("get task failed: %v", err)
	}

	cliui.PrintTask(res.GetTask())
}

func listCommand(args []string) {
	fs := flag.NewFlagSet("list", flag.ExitOnError)

	addr := fs.String("addr", "localhost:50051", "server address")
	queueName := fs.String("queue", "", "queue name")
	status := fs.String("status", "", "status filter")

	_ = fs.Parse(args)

	client, conn := mustClient(*addr)
	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := client.ListTasks(ctx, &pb.ListTasksRequest{
		Queue:  *queueName,
		Status: parseStatus(*status),
	})
	if err != nil {
		log.Fatalf("list tasks failed: %v", err)
	}

	cliui.PrintTaskList(res.GetTasks())
}

func cancelCommand(args []string) {
	fs := flag.NewFlagSet("cancel", flag.ExitOnError)

	addr := fs.String("addr", "localhost:50051", "server address")
	id := fs.String("id", "", "task id")

	_ = fs.Parse(args)

	client, conn := mustClient(*addr)
	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := client.CancelTask(ctx, &pb.CancelTaskRequest{
		Id: *id,
	})
	if err != nil {
		log.Fatalf("cancel task failed: %v", err)
	}

	cliui.PrintTask(res.GetTask())
}

func watchCommand(args []string) {
	fs := flag.NewFlagSet("watch", flag.ExitOnError)

	addr := fs.String("addr", "localhost:50051", "server address")
	id := fs.String("id", "", "task id")

	_ = fs.Parse(args)

	client, conn := mustClient(*addr)
	defer conn.Close()

	stream, err := client.WatchTask(context.Background(), &pb.WatchTaskRequest{
		Id: *id,
	})
	if err != nil {
		log.Fatalf("watch task failed: %v", err)
	}

	for {
		res, err := stream.Recv()
		if err == io.EOF {
			return
		}

		if err != nil {
			log.Fatalf("watch stream failed: %v", err)
		}

		cliui.PrintTask(res.GetTask())
	}
}

func mustClient(addr string) (pb.TaskServiceClient, *grpc.ClientConn) {
	conn, err := grpc.NewClient(
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("failed to create gRPC client: %v", err)
	}

	return pb.NewTaskServiceClient(conn), conn
}

func parsePriority(value string) pb.TaskPriority {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "low":
		return pb.TaskPriority_TASK_PRIORITY_LOW
	case "high":
		return pb.TaskPriority_TASK_PRIORITY_HIGH
	default:
		return pb.TaskPriority_TASK_PRIORITY_NORMAL
	}
}

func parseStatus(value string) pb.TaskStatus {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "pending":
		return pb.TaskStatus_TASK_STATUS_PENDING
	case "scheduled":
		return pb.TaskStatus_TASK_STATUS_SCHEDULED
	case "queued":
		return pb.TaskStatus_TASK_STATUS_QUEUED
	case "running":
		return pb.TaskStatus_TASK_STATUS_RUNNING
	case "completed":
		return pb.TaskStatus_TASK_STATUS_COMPLETED
	case "failed":
		return pb.TaskStatus_TASK_STATUS_FAILED
	case "cancelled":
		return pb.TaskStatus_TASK_STATUS_CANCELLED
	case "dead_letter":
		return pb.TaskStatus_TASK_STATUS_DEAD_LETTER
	default:
		return pb.TaskStatus_TASK_STATUS_UNSPECIFIED
	}
}

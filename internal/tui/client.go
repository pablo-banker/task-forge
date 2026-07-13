package tui

import (
	"context"
	"time"

	pb "github.com/pablo-banker/taskforge/proto/taskforge/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	addr string
	conn *grpc.ClientConn
	api  pb.TaskServiceClient
}

func NewClient(addr string) (*Client, error) {
	conn, err := grpc.NewClient(
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}

	return &Client{
		addr: addr,
		conn: conn,
		api:  pb.NewTaskServiceClient(conn),
	}, nil
}

func (c *Client) Close() error {
	if c.conn == nil {
		return nil
	}

	return c.conn.Close()
}

func (c *Client) ListTasks(status pb.TaskStatus) ([]*pb.Task, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := c.api.ListTasks(ctx, &pb.ListTasksRequest{
		Status: status,
	})
	if err != nil {
		return nil, err
	}

	return res.GetTasks(), nil
}

func (c *Client) GetTask(id string) (*pb.Task, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := c.api.GetTask(ctx, &pb.GetTaskRequest{
		Id: id,
	})
	if err != nil {
		return nil, err
	}

	return res.GetTask(), nil
}

func (c *Client) CancelTask(id string) (*pb.Task, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := c.api.CancelTask(ctx, &pb.CancelTaskRequest{
		Id: id,
	})
	if err != nil {
		return nil, err
	}

	return res.GetTask(), nil
}

func (c *Client) CreateTask(req *pb.CreateTaskRequest) (*pb.Task, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := c.api.CreateTask(ctx, req)
	if err != nil {
		return nil, err
	}

	return res.GetTask(), nil
}

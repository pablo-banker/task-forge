package tui

import pb "github.com/pablo-banker/taskforge/proto/taskforge/v1"

type tasksLoadedMsg struct {
	tasks []*pb.Task
	err   error
}

type taskLoadedMsg struct {
	task *pb.Task
	err  error
}

type taskCreatedMsg struct {
	task *pb.Task
	err  error
}

type taskCancelledMsg struct {
	task *pb.Task
	err  error
}

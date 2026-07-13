package tui

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	pb "github.com/pablo-banker/taskforge/proto/taskforge/v1"
)

func loadTasksCmd(client *Client) tea.Cmd {
	return func() tea.Msg {
		tasks, err := client.ListTasks(pb.TaskStatus_TASK_STATUS_UNSPECIFIED)

		return tasksLoadedMsg{
			tasks: tasks,
			err:   err,
		}
	}
}

func getTaskCmd(client *Client, id string) tea.Cmd {
	return func() tea.Msg {
		task, err := client.GetTask(id)

		return taskLoadedMsg{
			task: task,
			err:  err,
		}
	}
}

func cancelTaskCmd(client *Client, id string) tea.Cmd {
	return func() tea.Msg {
		task, err := client.CancelTask(id)

		return taskCancelledMsg{
			task: task,
			err:  err,
		}
	}
}

func createTaskCmd(client *Client, form CreateForm) tea.Cmd {
	return func() tea.Msg {
		delay, err := form.Delay()
		if err != nil {
			return taskCreatedMsg{
				err: err,
			}
		}

		maxAttempts, err := form.MaxAttempts()
		if err != nil {
			return taskCreatedMsg{
				err: err,
			}
		}

		runAtUnix := int64(0)
		if delay > 0 {
			runAtUnix = time.Now().Add(delay).Unix()
		}

		task, err := client.CreateTask(&pb.CreateTaskRequest{
			Queue:       form.Queue(),
			Name:        form.Name(),
			Payload:     form.Payload(),
			Priority:    pb.TaskPriority_TASK_PRIORITY_NORMAL,
			MaxAttempts: int32(maxAttempts),
			RunAtUnix:   runAtUnix,
		})

		return taskCreatedMsg{
			task: task,
			err:  err,
		}
	}
}

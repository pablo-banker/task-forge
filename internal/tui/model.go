package tui

import (
	"strconv"
	"strings"
	"time"

	pb "github.com/pablo-banker/taskforge/proto/taskforge/v1"
)

type Page int

const (
	PageDashboard Page = iota
	PageTasks
	PageCreate
	PageDetails
	PageDeadLetter
)

type Model struct {
	addr   string
	client *Client

	width  int
	height int

	page   Page
	cursor int

	tasks    []*pb.Task
	selected *pb.Task

	loading bool
	err     string
	notice  string

	form CreateForm
}

type CreateForm struct {
	fields []FormField
	index  int
}

type FormField struct {
	Label       string
	Placeholder string
	Value       string
}

func NewModel(addr string, client *Client) Model {
	return Model{
		addr:    addr,
		client:  client,
		page:    PageDashboard,
		cursor:  0,
		loading: true,
		form:    NewCreateForm(),
	}
}

func NewCreateForm() CreateForm {
	return CreateForm{
		index: 0,
		fields: []FormField{
			{
				Label:       "Queue",
				Placeholder: "default",
				Value:       "default",
			},
			{
				Label:       "Name",
				Placeholder: "welcome-email",
			},
			{
				Label:       "Payload",
				Placeholder: "send welcome email",
			},
			{
				Label:       "Delay",
				Placeholder: "10s, 1m, 0",
			},
			{
				Label:       "Max Attempts",
				Placeholder: "3",
				Value:       "3",
			},
		},
	}
}

func (f CreateForm) Queue() string {
	value := strings.TrimSpace(f.fields[0].Value)
	if value == "" {
		return "default"
	}

	return value
}

func (f CreateForm) Name() string {
	return strings.TrimSpace(f.fields[1].Value)
}

func (f CreateForm) Payload() string {
	return f.fields[2].Value
}

func (f CreateForm) Delay() (time.Duration, error) {
	value := strings.TrimSpace(f.fields[3].Value)
	if value == "" || value == "0" {
		return 0, nil
	}

	return time.ParseDuration(value)
}

func (f CreateForm) MaxAttempts() (int, error) {
	value := strings.TrimSpace(f.fields[4].Value)
	if value == "" {
		return 3, nil
	}

	attempts, err := strconv.Atoi(value)
	if err != nil {
		return 0, err
	}

	if attempts <= 0 {
		return 3, nil
	}

	return attempts, nil
}

func (m Model) CurrentTasks() []*pb.Task {
	if m.page == PageDeadLetter {
		result := make([]*pb.Task, 0)

		for _, task := range m.tasks {
			if task.GetStatus() == pb.TaskStatus_TASK_STATUS_DEAD_LETTER {
				result = append(result, task)
			}
		}

		return result
	}

	return m.tasks
}

func (m Model) CurrentTask() *pb.Task {
	tasks := m.CurrentTasks()
	if len(tasks) == 0 {
		return nil
	}

	if m.cursor < 0 || m.cursor >= len(tasks) {
		return nil
	}

	return tasks[m.cursor]
}

func pageTitle(page Page) string {
	switch page {
	case PageDashboard:
		return "Dashboard"
	case PageTasks:
		return "Tasks"
	case PageCreate:
		return "Create Task"
	case PageDetails:
		return "Task Details"
	case PageDeadLetter:
		return "Dead Letter"
	default:
		return "TaskForge"
	}
}

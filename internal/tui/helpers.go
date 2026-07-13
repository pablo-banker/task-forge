package tui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	pb "github.com/pablo-banker/taskforge/proto/taskforge/v1"
)

func shortID(id string) string {
	if len(id) <= 8 {
		return id
	}

	return id[:8]
}

func statusText(status pb.TaskStatus) string {
	value := strings.TrimPrefix(status.String(), "TASK_STATUS_")
	value = strings.ToLower(value)

	return value
}

func priorityText(priority pb.TaskPriority) string {
	value := strings.TrimPrefix(priority.String(), "TASK_PRIORITY_")
	value = strings.ToLower(value)

	return value
}

func statusStyled(status pb.TaskStatus) string {
	text := "● " + statusText(status)

	switch status {
	case pb.TaskStatus_TASK_STATUS_COMPLETED:
		return successStyle.Render(text)
	case pb.TaskStatus_TASK_STATUS_RUNNING:
		return accentStyle.Render(text)
	case pb.TaskStatus_TASK_STATUS_SCHEDULED:
		return lipglossColor(colorYellow, text)
	case pb.TaskStatus_TASK_STATUS_DEAD_LETTER:
		return errorStyle.Render(text)
	case pb.TaskStatus_TASK_STATUS_CANCELLED:
		return mutedStyle.Render(text)
	default:
		return lipglossColor(colorMagenta, text)
	}
}

func lipglossColor(color lipgloss.Color, value string) string {
	return lipgloss.NewStyle().Foreground(color).Render(value)
}

func formatUnix(value int64) string {
	if value <= 0 {
		return "-"
	}

	return time.Unix(value, 0).Format(time.RFC3339)
}

func countStatus(tasks []*pb.Task, status pb.TaskStatus) int {
	total := 0

	for _, task := range tasks {
		if task.GetStatus() == status {
			total++
		}
	}

	return total
}

func taskRow(index int, cursor int, task *pb.Task) string {
	prefix := " "
	if index == cursor {
		prefix = accentStyle.Render("›")
	}

	return fmt.Sprintf(
		"%s %-10s %-18s %-12s %-28s %d/%d",
		prefix,
		shortID(task.GetId()),
		statusText(task.GetStatus()),
		task.GetQueue(),
		truncate(task.GetName(), 28),
		task.GetAttempts(),
		task.GetMaxAttempts(),
	)
}

func truncate(value string, max int) string {
	runes := []rune(value)

	if len(runes) <= max {
		return value
	}

	return string(runes[:max-1]) + "…"
}

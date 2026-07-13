package cliui

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"

	pb "github.com/pablo-banker/taskforge/proto/taskforge/v1"
)

const (
	reset   = "\033[0m"
	bold    = "\033[1m"
	dim     = "\033[2m"
	cyan    = "\033[36m"
	green   = "\033[32m"
	yellow  = "\033[33m"
	red     = "\033[31m"
	magenta = "\033[35m"
	gray    = "\033[90m"
)

var ansiRegex = regexp.MustCompile(`\x1b\[[0-9;]*m`)

func color(code string, value string) string {
	if os.Getenv("NO_COLOR") != "" {
		return value
	}

	return code + value + reset
}

func PrintTask(t *pb.Task) {
	if t == nil {
		fmt.Println(color(red, "Task not found."))
		return
	}

	status := cleanEnum(t.GetStatus().String(), "TASK_STATUS_")
	priority := cleanEnum(t.GetPriority().String(), "TASK_PRIORITY_")

	rows := []struct {
		label string
		value string
	}{
		{"ID", t.GetId()},
		{"Queue", t.GetQueue()},
		{"Name", t.GetName()},
		{"Payload", t.GetPayload()},
		{"Status", statusBadge(status)},
		{"Priority", priority},
		{"Attempts", fmt.Sprintf("%d/%d", t.GetAttempts(), t.GetMaxAttempts())},
		{"Run at", formatUnix(t.GetRunAtUnix())},
		{"Created", formatUnix(t.GetCreatedAtUnix())},
		{"Updated", formatUnix(t.GetUpdatedAtUnix())},
	}

	if t.GetLastError() != "" {
		rows = append(rows, struct {
			label string
			value string
		}{
			label: "Error",
			value: color(red, t.GetLastError()),
		})
	}

	width := 86
	title := fmt.Sprintf(" Task %s ", shortID(t.GetId()))

	printTop(title, width)

	for _, row := range rows {
		printRow(row.label, row.value, width)
	}

	printBottom(width)
}

func PrintTaskList(tasks []*pb.Task) {
	if len(tasks) == 0 {
		fmt.Println(color(gray, "No tasks found."))
		return
	}

	width := 100

	printTop(" Tasks ", width)
	printTableHeader(width)

	for _, t := range tasks {
		id := shortID(t.GetId())
		status := statusBadge(cleanEnum(t.GetStatus().String(), "TASK_STATUS_"))
		queue := truncate(t.GetQueue(), 14)
		name := truncate(t.GetName(), 28)
		attempts := fmt.Sprintf("%d/%d", t.GetAttempts(), t.GetMaxAttempts())

		line := fmt.Sprintf(
			"│ %-10s %-22s %-14s %-28s %-8s",
			id,
			status,
			queue,
			name,
			attempts,
		)

		fmt.Println(line + spaces(width-visibleLen(line)-1) + "│")
	}

	printBottom(width)
}

func PrintSuccess(message string) {
	fmt.Println(color(green, "✓ "+message))
}

func PrintInfo(message string) {
	fmt.Println(color(cyan, "ℹ "+message))
}

func PrintWarning(message string) {
	fmt.Println(color(yellow, "⚠ "+message))
}

func PrintError(message string) {
	fmt.Println(color(red, "✕ "+message))
}

func printTop(title string, width int) {
	left := "╭─" + color(cyan, title)
	line := left + strings.Repeat("─", max(0, width-visibleLen(left)-1)) + "╮"

	fmt.Println(color(cyan, line))
}

func printBottom(width int) {
	fmt.Println(color(cyan, "╰"+strings.Repeat("─", width-2)+"╯"))
}

func printRow(label string, value string, width int) {
	coloredLabel := color(gray, padRight(label, 12))
	line := fmt.Sprintf("│ %s %s", coloredLabel, value)

	fmt.Println(line + spaces(width-visibleLen(line)-1) + "│")
}

func printTableHeader(width int) {
	header := fmt.Sprintf(
		"│ %-10s %-14s %-14s %-28s %-8s",
		"ID",
		"Status",
		"Queue",
		"Name",
		"Attempts",
	)

	fmt.Println(color(gray, header+spaces(width-visibleLen(header)-1)+"│"))
	fmt.Println(color(cyan, "├"+strings.Repeat("─", width-2)+"┤"))
}

func statusBadge(status string) string {
	switch status {
	case "completed":
		return color(green, "● completed")
	case "running":
		return color(cyan, "● running")
	case "scheduled":
		return color(yellow, "● scheduled")
	case "pending":
		return color(magenta, "● pending")
	case "cancelled":
		return color(gray, "● cancelled")
	case "dead_letter":
		return color(red, "● dead_letter")
	case "failed":
		return color(red, "● failed")
	default:
		return color(gray, "● "+status)
	}
}

func cleanEnum(value string, prefix string) string {
	value = strings.TrimPrefix(value, prefix)
	value = strings.ToLower(value)

	return value
}

func formatUnix(value int64) string {
	if value <= 0 {
		return "-"
	}

	return time.Unix(value, 0).Format(time.RFC3339)
}

func shortID(id string) string {
	if len(id) <= 8 {
		return id
	}

	return id[:8]
}

func truncate(value string, maxLength int) string {
	if visibleLen(value) <= maxLength {
		return value
	}

	runes := []rune(value)

	if len(runes) <= maxLength {
		return value
	}

	return string(runes[:maxLength-1]) + "…"
}

func padRight(value string, width int) string {
	length := visibleLen(value)

	if length >= width {
		return value
	}

	return value + spaces(width-length)
}

func spaces(count int) string {
	if count <= 0 {
		return ""
	}

	return strings.Repeat(" ", count)
}

func visibleLen(value string) int {
	clean := ansiRegex.ReplaceAllString(value, "")
	return utf8.RuneCountInString(clean)
}

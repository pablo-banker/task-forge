package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	pb "github.com/pablo-banker/taskforge/proto/taskforge/v1"
)

func (m Model) View() string {
	width := m.width
	height := m.height

	if width <= 0 {
		width = 120
	}

	if height <= 0 {
		height = 34
	}

	return appStyle.
		Width(width).
		Height(height).
		Render(m.render(width, height))
}

func (m Model) render(width int, height int) string {
	contentWidth := width - 4
	contentHeight := height - 2

	if contentWidth < 90 {
		contentWidth = 90
	}

	if contentHeight < 24 {
		contentHeight = 24
	}

	sidebarWidth := 26
	gapWidth := 2
	footerHeight := 2

	mainWidth := contentWidth - sidebarWidth - gapWidth
	mainHeight := contentHeight - footerHeight

	if mainWidth < 60 {
		mainWidth = 60
	}

	if mainHeight < 18 {
		mainHeight = 18
	}

	sidebar := sidebarStyle.
		Width(sidebarWidth).
		Height(mainHeight).
		Render(m.renderSidebarContent())

	main := panelStyle.
		Width(mainWidth).
		Height(mainHeight).
		Render(m.renderPage())

	body := lipgloss.JoinHorizontal(
		lipgloss.Top,
		sidebar,
		strings.Repeat(" ", gapWidth),
		main,
	)

	footer := footerStyle.
		Width(contentWidth).
		Render("\n" + m.renderFooter())

	return body + footer
}

func (m Model) renderSidebarContent() string {
	items := []struct {
		page  Page
		label string
	}{
		{PageDashboard, "Dashboard"},
		{PageTasks, "Tasks"},
		{PageCreate, "Create Task"},
		{PageDeadLetter, "Dead Letter"},
	}

	var b strings.Builder

	b.WriteString(headerStyle.Render("TaskForge"))
	b.WriteString("\n")
	b.WriteString(mutedStyle.Render(m.addr))
	b.WriteString("\n\n")

	for _, item := range items {
		if m.page == item.page {
			b.WriteString(activeNavItemStyle.Render("› " + item.label))
		} else {
			b.WriteString(navItemStyle.Render("  " + item.label))
		}

		b.WriteString("\n")
	}

	b.WriteString("\n")
	b.WriteString(mutedStyle.Render("Shortcuts"))
	b.WriteString("\n")

	if m.page == PageCreate {
		b.WriteString(mutedStyle.Render("tab next field"))
		b.WriteString("\n")
		b.WriteString(mutedStyle.Render("enter next/create"))
		b.WriteString("\n")
		b.WriteString(mutedStyle.Render("ctrl+s create"))
		b.WriteString("\n")
		b.WriteString(mutedStyle.Render("esc back"))
	} else {
		b.WriteString(mutedStyle.Render("↑/↓ or j/k move"))
		b.WriteString("\n")
		b.WriteString(mutedStyle.Render("enter open"))
		b.WriteString("\n")
		b.WriteString(mutedStyle.Render("r refresh"))
	}

	return b.String()
}

func (m Model) renderFooter() string {
	if m.page == PageCreate {
		return "q quit · esc back · tab/enter next field · ctrl+s create"
	}

	return "q quit · tab next · 1 dashboard · 2 tasks · 3 create · 4 dead letter · r refresh"
}

func (m Model) renderPage() string {
	if m.err != "" {
		return errorStyle.Render("Error: "+m.err) + "\n\n" +
			mutedStyle.Render("Start the server with: make run-server")
	}

	switch m.page {
	case PageDashboard:
		return m.renderDashboard()

	case PageTasks:
		return m.renderTasks("Tasks", m.CurrentTasks())

	case PageCreate:
		return m.renderCreate()

	case PageDetails:
		return m.renderDetails()

	case PageDeadLetter:
		return m.renderTasks("Dead Letter", m.CurrentTasks())

	default:
		return "Unknown page"
	}
}

func (m Model) renderDashboard() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("Dashboard"))
	b.WriteString("\n\n")

	if m.loading {
		b.WriteString(mutedStyle.Render("Loading tasks..."))
		return b.String()
	}

	if m.notice != "" {
		b.WriteString(successStyle.Render(m.notice))
		b.WriteString("\n\n")
	}

	rows := []string{
		fmt.Sprintf("Total        %d", len(m.tasks)),
		fmt.Sprintf("Pending      %d", countStatus(m.tasks, pb.TaskStatus_TASK_STATUS_PENDING)),
		fmt.Sprintf("Scheduled    %d", countStatus(m.tasks, pb.TaskStatus_TASK_STATUS_SCHEDULED)),
		fmt.Sprintf("Running      %d", countStatus(m.tasks, pb.TaskStatus_TASK_STATUS_RUNNING)),
		fmt.Sprintf("Completed    %d", countStatus(m.tasks, pb.TaskStatus_TASK_STATUS_COMPLETED)),
		fmt.Sprintf("Cancelled    %d", countStatus(m.tasks, pb.TaskStatus_TASK_STATUS_CANCELLED)),
		fmt.Sprintf("Dead Letter  %d", countStatus(m.tasks, pb.TaskStatus_TASK_STATUS_DEAD_LETTER)),
	}

	for _, row := range rows {
		b.WriteString(row)
		b.WriteString("\n")
	}

	b.WriteString("\n")
	b.WriteString(mutedStyle.Render("Press 2 to browse tasks or 3 to create a new task."))

	return b.String()
}

func (m Model) renderTasks(title string, tasks []*pb.Task) string {
	var b strings.Builder

	b.WriteString(titleStyle.Render(title))
	b.WriteString("\n\n")

	if m.loading {
		b.WriteString(mutedStyle.Render("Loading tasks..."))
		return b.String()
	}

	if m.notice != "" {
		b.WriteString(successStyle.Render(m.notice))
		b.WriteString("\n\n")
	}

	if len(tasks) == 0 {
		b.WriteString(mutedStyle.Render("No tasks found."))
		return b.String()
	}

	b.WriteString(mutedStyle.Render("  ID         Status             Queue        Name                         Attempts"))
	b.WriteString("\n")
	b.WriteString(mutedStyle.Render("  -----------------------------------------------------------------------------"))
	b.WriteString("\n")

	for i, task := range tasks {
		line := taskRow(i, m.cursor, task)

		if i == m.cursor {
			b.WriteString(accentStyle.Render(line))
		} else {
			b.WriteString(line)
		}

		b.WriteString("\n")
	}

	b.WriteString("\n")
	b.WriteString(mutedStyle.Render("enter details · c create · x cancel selected · r refresh"))

	return b.String()
}

func (m Model) renderDetails() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("Task Details"))
	b.WriteString("\n\n")

	if m.selected == nil {
		b.WriteString(mutedStyle.Render("No task selected."))
		return b.String()
	}

	task := m.selected

	rows := []string{
		fmt.Sprintf("ID           %s", task.GetId()),
		fmt.Sprintf("Queue        %s", task.GetQueue()),
		fmt.Sprintf("Name         %s", task.GetName()),
		fmt.Sprintf("Payload      %s", task.GetPayload()),
		fmt.Sprintf("Status       %s", statusStyled(task.GetStatus())),
		fmt.Sprintf("Priority     %s", priorityText(task.GetPriority())),
		fmt.Sprintf("Attempts     %d/%d", task.GetAttempts(), task.GetMaxAttempts()),
		fmt.Sprintf("Run at       %s", formatUnix(task.GetRunAtUnix())),
		fmt.Sprintf("Created      %s", formatUnix(task.GetCreatedAtUnix())),
		fmt.Sprintf("Updated      %s", formatUnix(task.GetUpdatedAtUnix())),
	}

	for _, row := range rows {
		b.WriteString(row)
		b.WriteString("\n")
	}

	if task.GetLastError() != "" {
		b.WriteString("\n")
		b.WriteString(errorStyle.Render("Last error: " + task.GetLastError()))
		b.WriteString("\n")
	}

	b.WriteString("\n")
	b.WriteString(mutedStyle.Render("b back · x cancel · r refresh"))

	return b.String()
}

func (m Model) renderCreate() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("Create Task"))
	b.WriteString("\n\n")

	if m.notice != "" {
		b.WriteString(successStyle.Render(m.notice))
		b.WriteString("\n\n")
	}

	for i, field := range m.form.fields {
		style := inputStyle
		cursor := " "

		if i == m.form.index {
			style = activeInputStyle
			cursor = accentStyle.Render("›")
		}

		value := field.Value
		if value == "" {
			value = mutedStyle.Render(field.Placeholder)
		}

		b.WriteString(cursor)
		b.WriteString(" ")
		b.WriteString(field.Label)
		b.WriteString("\n")
		b.WriteString(style.Width(52).Render(value))
		b.WriteString("\n\n")
	}

	b.WriteString(mutedStyle.Render("enter next field · ctrl+s create · esc back"))

	return b.String()
}

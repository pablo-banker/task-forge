package tui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		tea.ClearScreen,
		loadTasksCmd(m.client),
	)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		return m, nil

	case tasksLoadedMsg:
		m.loading = false

		if msg.err != nil {
			m.err = msg.err.Error()
			return m, tea.ClearScreen
		}

		m.err = ""
		m.tasks = msg.tasks
		m.fixCursor()

		return m, tea.ClearScreen

	case taskLoadedMsg:
		m.loading = false

		if msg.err != nil {
			m.err = msg.err.Error()
			return m, tea.ClearScreen
		}

		m.err = ""
		m.selected = msg.task
		m.page = PageDetails

		return m, tea.ClearScreen

	case taskCreatedMsg:
		m.loading = false

		if msg.err != nil {
			m.err = msg.err.Error()
			return m, tea.ClearScreen
		}

		m.err = ""
		m.notice = "Task created successfully."
		m.selected = msg.task
		m.form = NewCreateForm()
		m.page = PageDetails

		return m, tea.Batch(
			tea.ClearScreen,
			loadTasksCmd(m.client),
		)

	case taskCancelledMsg:
		m.loading = false

		if msg.err != nil {
			m.err = msg.err.Error()
			return m, tea.ClearScreen
		}

		m.err = ""
		m.notice = "Task cancelled."
		m.selected = msg.task

		return m, tea.Batch(
			tea.ClearScreen,
			loadTasksCmd(m.client),
		)

	case tea.KeyMsg:
		return m.handleKey(msg)
	}

	return m, nil
}

func (m Model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	key := msg.String()

	if m.page == PageCreate {
		return m.handleCreateKey(key)
	}

	switch key {
	case "ctrl+c", "q":
		return m, tea.Quit

	case "tab":
		m.nextPage()
		return m, tea.ClearScreen

	case "1":
		m.page = PageDashboard
		m.cursor = 0
		m.err = ""
		return m, tea.ClearScreen

	case "2":
		m.page = PageTasks
		m.cursor = 0
		m.err = ""
		return m, tea.ClearScreen

	case "3", "c":
		m.page = PageCreate
		m.notice = ""
		m.err = ""
		return m, tea.ClearScreen

	case "4":
		m.page = PageDeadLetter
		m.cursor = 0
		m.err = ""
		return m, tea.ClearScreen

	case "r":
		m.loading = true
		m.err = ""
		return m, tea.Batch(
			tea.ClearScreen,
			loadTasksCmd(m.client),
		)

	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
		}

		return m, nil

	case "down", "j":
		if m.cursor < len(m.CurrentTasks())-1 {
			m.cursor++
		}

		return m, nil

	case "enter":
		task := m.CurrentTask()
		if task == nil {
			return m, nil
		}

		m.loading = true
		m.selected = task

		return m, tea.Batch(
			tea.ClearScreen,
			getTaskCmd(m.client, task.GetId()),
		)

	case "b":
		m.page = PageTasks
		m.err = ""
		return m, tea.ClearScreen

	case "x":
		task := m.CurrentTask()
		if m.page == PageDetails && m.selected != nil {
			task = m.selected
		}

		if task == nil {
			return m, nil
		}

		m.loading = true

		return m, tea.Batch(
			tea.ClearScreen,
			cancelTaskCmd(m.client, task.GetId()),
		)
	}

	return m, nil
}

func (m Model) handleCreateKey(key string) (tea.Model, tea.Cmd) {
	switch key {
	case "ctrl+c", "q":
		return m, tea.Quit

	case "esc":
		m.page = PageTasks
		m.err = ""
		m.notice = ""
		m.cursor = 0

		return m, tea.Batch(
			tea.ClearScreen,
			loadTasksCmd(m.client),
		)

	case "ctrl+s":
		return m.submitCreateForm()

	case "enter":
		if m.form.index < len(m.form.fields)-1 {
			m.form.index++
			return m, tea.ClearScreen
		}

		return m.submitCreateForm()

	case "shift+tab", "up":
		if m.form.index > 0 {
			m.form.index--
		}

		return m, nil

	case "tab", "down":
		if m.form.index < len(m.form.fields)-1 {
			m.form.index++
		}

		return m, nil

	case "backspace":
		m.removeLastRune()
		return m, nil

	case "space":
		m.appendToField(" ")
		return m, nil
	}

	if len([]rune(key)) == 1 {
		m.appendToField(key)
	}

	return m, nil
}

func (m Model) submitCreateForm() (tea.Model, tea.Cmd) {
	if strings.TrimSpace(m.form.Name()) == "" {
		m.err = "Task name is required."
		return m, tea.ClearScreen
	}

	m.loading = true
	m.err = ""
	m.notice = ""

	return m, tea.Batch(
		tea.ClearScreen,
		createTaskCmd(m.client, m.form),
	)
}

func (m *Model) appendToField(value string) {
	if m.form.index < 0 || m.form.index >= len(m.form.fields) {
		return
	}

	m.form.fields[m.form.index].Value += value
}

func (m *Model) removeLastRune() {
	if m.form.index < 0 || m.form.index >= len(m.form.fields) {
		return
	}

	value := []rune(m.form.fields[m.form.index].Value)
	if len(value) == 0 {
		return
	}

	m.form.fields[m.form.index].Value = string(value[:len(value)-1])
}

func (m *Model) nextPage() {
	switch m.page {
	case PageDashboard:
		m.page = PageTasks
	case PageTasks:
		m.page = PageCreate
	case PageCreate:
		m.page = PageDeadLetter
	case PageDeadLetter:
		m.page = PageDashboard
	case PageDetails:
		m.page = PageTasks
	default:
		m.page = PageDashboard
	}

	m.err = ""
	m.cursor = 0
}

func (m *Model) fixCursor() {
	tasks := m.CurrentTasks()

	if len(tasks) == 0 {
		m.cursor = 0
		return
	}

	if m.cursor >= len(tasks) {
		m.cursor = len(tasks) - 1
	}

	if m.cursor < 0 {
		m.cursor = 0
	}
}

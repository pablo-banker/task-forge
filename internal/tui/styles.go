package tui

import "github.com/charmbracelet/lipgloss"

var (
	colorBg      = lipgloss.Color("#09090B")
	colorPanel   = lipgloss.Color("#111113")
	colorText    = lipgloss.Color("#FAFAFA")
	colorMuted   = lipgloss.Color("#9CA3AF")
	colorAccent  = lipgloss.Color("#22D3EE")
	colorGreen   = lipgloss.Color("#22C55E")
	colorYellow  = lipgloss.Color("#EAB308")
	colorRed     = lipgloss.Color("#EF4444")
	colorMagenta = lipgloss.Color("#A855F7")

	appStyle = lipgloss.NewStyle().
			Background(colorBg).
			Foreground(colorText)

	headerStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorAccent)

	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorText)

	mutedStyle = lipgloss.NewStyle().
			Foreground(colorMuted)

	accentStyle = lipgloss.NewStyle().
			Foreground(colorAccent)

	errorStyle = lipgloss.NewStyle().
			Foreground(colorRed)

	successStyle = lipgloss.NewStyle().
			Foreground(colorGreen)

	sidebarStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorAccent).
			Padding(1, 1)

	panelStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorAccent).
			Padding(1, 2)

	footerStyle = lipgloss.NewStyle().
			Foreground(colorMuted)

	navItemStyle = lipgloss.NewStyle().
			Foreground(colorMuted).
			PaddingLeft(1)

	activeNavItemStyle = lipgloss.NewStyle().
				Foreground(colorAccent).
				Bold(true).
				PaddingLeft(1)

	inputStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorMuted).
			Padding(0, 1)

	activeInputStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(colorAccent).
				Padding(0, 1)
)

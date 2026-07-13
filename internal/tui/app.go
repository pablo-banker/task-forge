package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

func Run(addr string) error {
	client, err := NewClient(addr)
	if err != nil {
		return fmt.Errorf("failed to create TaskForge client: %w", err)
	}
	defer client.Close()

	program := tea.NewProgram(
		NewModel(addr, client),
		tea.WithAltScreen(),
	)

	if _, err := program.Run(); err != nil {
		return fmt.Errorf("tui failed: %w", err)
	}

	return nil
}

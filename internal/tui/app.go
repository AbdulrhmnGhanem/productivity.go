package tui

import (
"fmt"

tea "github.com/charmbracelet/bubbletea"
"productivity.go/internal/readings"
)

func Start(service *readings.Service) error {
	model, err := InitTUI(service)
	if err != nil {
		return fmt.Errorf("failed to initialize TUI: %w", err)
	}

	p := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		return fmt.Errorf("failed to run TUI: %w", err)
	}

	return nil
}

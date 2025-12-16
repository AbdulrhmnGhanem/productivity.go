package setup

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"productivity.go/internal/config"
)

type model struct {
	step   int
	inputs []textinput.Model
	err    error
	done   bool
}

func InitialModel() model {
	tiKey := textinput.New()
	tiKey.Placeholder = "Notion API Key"
	tiKey.Focus()
	tiKey.CharLimit = 100
	tiKey.Width = 50
	tiKey.EchoMode = textinput.EchoPassword

	tiDB := textinput.New()
	tiDB.Placeholder = "Notion Database ID"
	tiDB.CharLimit = 100
	tiDB.Width = 50

	return model{
		step:   0,
		inputs: []textinput.Model{tiKey, tiDB},
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit
		case "enter":
			if m.step < len(m.inputs)-1 {
				m.inputs[m.step].Blur()
				m.step++
				m.inputs[m.step].Focus()
				return m, textinput.Blink
			}
			// Finished
			apiKey := m.inputs[0].Value()
			dbID := config.CleanDatabaseID(m.inputs[1].Value())
			if err := config.Save(apiKey, dbID); err != nil {
				m.err = err
				return m, tea.Quit
			}
			m.done = true
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.inputs[m.step], cmd = m.inputs[m.step].Update(msg)
	return m, cmd
}

func (m model) View() string {
	if m.err != nil {
		return fmt.Sprintf("Error: %v\n", m.err)
	}
	if m.done {
		return "Configuration saved successfully!\n"
	}

	s := "Readings CLI Setup\n\n"

	if m.step == 0 {
		s += "Enter your Notion API Key:\n"
		s += m.inputs[0].View()
	} else if m.step == 1 {
		s += "Enter your Notion Database ID:\n"
		s += m.inputs[1].View()
	}

	s += "\n\n(esc to quit)\n"
	return s
}

func Run() error {
	p := tea.NewProgram(InitialModel())
	if _, err := p.Run(); err != nil {
		return err
	}
	return nil
}

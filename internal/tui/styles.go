package tui

import "github.com/charmbracelet/lipgloss"

type Styles struct {
	Title        lipgloss.Style
	Item         lipgloss.Style
	SelectedItem lipgloss.Style
	FilterTitle  lipgloss.Style
	FilterItem   lipgloss.Style
	DetailTitle  lipgloss.Style
	DetailInfo   lipgloss.Style
	HelpBar      lipgloss.Style
	HelpKey      lipgloss.Style
	HelpDesc     lipgloss.Style
}

func DefaultStyles() Styles {
	return Styles{
		Title: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).
			Padding(0, 1),
		Item: lipgloss.NewStyle().
			PaddingLeft(2),
		SelectedItem: lipgloss.NewStyle().
			PaddingLeft(2).
			Foreground(lipgloss.Color("205")),
		FilterTitle: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#04B575")).
			Padding(0, 1),
		FilterItem: lipgloss.NewStyle().
			PaddingLeft(2),
		DetailTitle: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("205")).
			MarginBottom(1),
		DetailInfo: lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")),
		HelpBar: lipgloss.NewStyle().
			Width(100). // Will be updated dynamically if needed, or just let it flow
			Foreground(lipgloss.Color("#A8A8A8")).
			Background(lipgloss.Color("#303030")).
			Padding(0, 1),
		HelpKey: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FAFAFA")),
		HelpDesc: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#A8A8A8")).
			PaddingRight(1),
	}
}

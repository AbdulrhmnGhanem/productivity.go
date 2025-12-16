package tui

import (
"fmt"
"strings"

"github.com/charmbracelet/lipgloss"
)

func (m Model) View() string {
	styles := DefaultStyles()

	var content string
	switch m.view {
	case ViewList:
		content = m.viewList(styles)
	case ViewDetail:
		content = m.viewDetail(styles)
	case ViewFilter:
		content = m.viewFilter(styles)
	default:
		content = "Unknown view"
	}

	help := m.helpView(styles)
	
	return lipgloss.JoinVertical(lipgloss.Left, content, help)
}

func (m Model) viewList(styles Styles) string {
	var b strings.Builder

	b.WriteString(styles.Title.Render("Readings"))
	b.WriteString("\n\n")

	if len(m.filteredArticles) == 0 {
		b.WriteString(styles.Item.Render("No articles found."))
		// Fill remaining height
		linesUsed := 3 // Title + 2 newlines + 1 item
		for i := 0; i < m.height-linesUsed-1; i++ { // -1 for help bar
			b.WriteString("\n")
		}
		return b.String()
	}

	// Calculate visible range
	// Header (2) + Help (1) = 3 lines reserved.
	availableHeight := m.height - 3
	if availableHeight < 1 {
		availableHeight = 1
	}

	start := m.scrollOffset
	end := start + availableHeight
	if end > len(m.filteredArticles) {
		end = len(m.filteredArticles)
	}
	if start > end {
		start = end
	}

	for i := start; i < end; i++ {
		article := m.filteredArticles[i]
		title := article.Title
		if title == "" {
			title = "Untitled"
		}

		if i == m.cursor {
			b.WriteString(styles.SelectedItem.Render("> " + title))
		} else {
			b.WriteString(styles.Item.Render(title))
		}
		b.WriteString("\n")
	}

	// Fill empty space if list is shorter than height
	linesRendered := end - start
	linesToFill := availableHeight - linesRendered
	for i := 0; i < linesToFill; i++ {
		b.WriteString("\n")
	}

	return b.String()
}

func (m Model) viewDetail(styles Styles) string {
	if m.cursor >= len(m.filteredArticles) {
		return "No article selected"
	}
	article := m.filteredArticles[m.cursor]

	var b strings.Builder
	b.WriteString(styles.DetailTitle.Render(article.Title))
	b.WriteString("\n\n")
	b.WriteString(styles.DetailInfo.Render(fmt.Sprintf("URL: %s", article.URL)))
	b.WriteString("\n")
	b.WriteString(styles.DetailInfo.Render(fmt.Sprintf("Fetched: %s", article.FetchedAt.Format("2006-01-02 15:04"))))
	
	content := b.String()
	return lipgloss.Place(m.width, m.height-1, lipgloss.Top, lipgloss.Left, content)
}

func (m Model) viewFilter(styles Styles) string {
	var b strings.Builder

	b.WriteString(styles.FilterTitle.Render("Filter by Tags"))
	b.WriteString("\n\n")

	if len(m.tags) == 0 {
		b.WriteString(styles.FilterItem.Render("No tags available."))
		return lipgloss.Place(m.width, m.height-1, lipgloss.Top, lipgloss.Left, b.String())
	}

	// Calculate visible range
	availableHeight := m.height - 3
	if availableHeight < 1 {
		availableHeight = 1
	}

	start := m.scrollOffset
	end := start + availableHeight
	if end > len(m.tags) {
		end = len(m.tags)
	}
	if start > end {
		start = end
	}

	for i := start; i < end; i++ {
		tag := m.tags[i]
		selected := m.selectedTags[tag]
		prefix := "[ ] "
		if selected {
			prefix = "[x] "
		}

		if i == m.cursor {
			b.WriteString(styles.SelectedItem.Render("> " + prefix + tag))
		} else {
			b.WriteString(styles.FilterItem.Render(prefix + tag))
		}
		b.WriteString("\n")
	}
	
	// Fill empty space
	linesRendered := end - start
	linesToFill := availableHeight - linesRendered
	for i := 0; i < linesToFill; i++ {
		b.WriteString("\n")
	}

	return b.String()
}

func (m Model) helpView(styles Styles) string {
	var keys []string
	switch m.view {
	case ViewList:
		keys = []string{"j/k", "nav", "/", "filter", "enter", "details", "q", "quit"}
	case ViewDetail:
		keys = []string{"enter", "open url", "esc", "back", "q", "back"}
	case ViewFilter:
		keys = []string{"j/k", "nav", "space", "toggle", "right", "all", "enter", "apply", "esc", "cancel"}
	}

	var b strings.Builder
	for i := 0; i < len(keys); i += 2 {
		b.WriteString(styles.HelpKey.Render(keys[i]))
		b.WriteString(" ")
		b.WriteString(styles.HelpDesc.Render(keys[i+1]))
	}
	
	// Append input buffer if present
	if m.inputBuffer != "" {
		b.WriteString(styles.HelpKey.Render(fmt.Sprintf("  %s", m.inputBuffer)))
	}

	return styles.HelpBar.Width(m.width).Render(b.String())
}

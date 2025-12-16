package tui

import (
	"context"
	"fmt"
	"os/exec"
	"runtime"
	"strconv"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"productivity.go/internal/readings"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			if m.view != ViewFilter { 
				return m, tea.Quit
			}
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case StatusMsg:
		m.statusMessage = string(msg)
		return m, tea.Tick(2*time.Second, func(_ time.Time) tea.Msg {
			return ClearStatusMsg{}
		})
	case ClearStatusMsg:
		m.statusMessage = ""
		return m, nil
	}

	switch m.view {
	case ViewList:
		return m.updateList(msg)
	case ViewDetail:
		return m.updateDetail(msg)
	case ViewFilter:
		return m.updateFilter(msg)
	}

	return m, nil
}

func (m Model) updateList(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Handle numeric input
		if len(msg.String()) == 1 && msg.String() >= "0" && msg.String() <= "9" {
			m.inputBuffer += msg.String()
			return m, nil
		}

		count := 1
		if m.inputBuffer != "" {
			if c, err := strconv.Atoi(m.inputBuffer); err == nil && c > 0 {
				count = c
			}
		}

		switch msg.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor -= count
				if m.cursor < 0 {
					m.cursor = 0
				}
				if m.cursor < m.scrollOffset {
					m.scrollOffset = m.cursor
				}
			}
			m.inputBuffer = ""
		case "down", "j":
			if m.cursor < len(m.filteredArticles)-1 {
				m.cursor += count
				if m.cursor >= len(m.filteredArticles) {
					m.cursor = len(m.filteredArticles) - 1
				}
				if m.cursor >= m.scrollOffset+(m.height-4) {
					m.scrollOffset = m.cursor - (m.height - 4) + 1
					if m.scrollOffset < 0 {
						m.scrollOffset = 0
					}
				}
			}
			m.inputBuffer = ""
		case "g":
			if m.inputBuffer == "" {
				// Wait for second 'g' or handle 'gg' logic if we had a key buffer for commands
				// But here we only buffer numbers. 
				// Standard Vim 'gg' is tricky without a command buffer.
				// Let's assume 'g' goes to top if no number, or line N if number?
// Actually 'gg' is go to top. 'G' is go to bottom.
// '10G' is go to line 10.
// '10gg' is go to line 10.

// Since we don't have a command buffer for 'gg', let's just implement 'G' for now.
// Or we can treat 'g' as 'home' if no buffer?
m.cursor = 0
m.scrollOffset = 0
} else {
// If number provided, go to that line (1-based)
target := count - 1
if target < 0 { target = 0 }
if target >= len(m.filteredArticles) { target = len(m.filteredArticles) - 1 }
m.cursor = target
// Adjust scroll
if m.cursor < m.scrollOffset || m.cursor >= m.scrollOffset+(m.height-4) {
m.scrollOffset = m.cursor - (m.height-4)/2
if m.scrollOffset < 0 { m.scrollOffset = 0 }
}
}
m.inputBuffer = ""
case "G":
if m.inputBuffer == "" {
// Go to bottom
m.cursor = len(m.filteredArticles) - 1
} else {
// Go to line N
target := count - 1
if target < 0 { target = 0 }
if target >= len(m.filteredArticles) { target = len(m.filteredArticles) - 1 }
m.cursor = target
}
// Adjust scroll
if m.cursor < m.scrollOffset || m.cursor >= m.scrollOffset+(m.height-4) {
m.scrollOffset = m.cursor - (m.height-4)/2
if m.scrollOffset < 0 { m.scrollOffset = 0 }
}
m.inputBuffer = ""

case "enter":
			if len(m.filteredArticles) > 0 {
				article := m.filteredArticles[m.cursor]
				m.inputBuffer = ""
				return m, func() tea.Msg {
					added, err := m.svc.ToggleReadingInCurrentWeek(context.Background(), article.ID)
					if err != nil {
						return StatusMsg(fmt.Sprintf("Error: %v", err))
					}
					if added {
						return StatusMsg("Added to reading list")
					}
					return StatusMsg("Removed from reading list")
				}
			}
			m.inputBuffer = ""
		case "/":
// Enter filter mode
m.view = ViewFilter
m.backupSelectedTags = make(map[string]bool)
for k, v := range m.selectedTags {
m.backupSelectedTags[k] = v
}
m.cursor = 0
m.scrollOffset = 0
m.inputBuffer = ""
default:
// Clear buffer on unknown key
m.inputBuffer = ""
}
}
return m, nil
}

func (m Model) updateDetail(msg tea.Msg) (tea.Model, tea.Cmd) {
switch msg := msg.(type) {
case tea.KeyMsg:
switch msg.String() {
case "esc", "q":
m.view = ViewList
case "enter":
if m.cursor < len(m.filteredArticles) {
article := m.filteredArticles[m.cursor]
return m, openUrl(article.URL)
}
}
}
return m, nil
}

func (m Model) updateFilter(msg tea.Msg) (tea.Model, tea.Cmd) {
switch msg := msg.(type) {
case tea.KeyMsg:
switch msg.String() {
case "up", "k":
if m.cursor > 0 {
m.cursor--
if m.cursor < m.scrollOffset {
m.scrollOffset = m.cursor
}
}
case "down", "j":
if m.cursor < len(m.tags)-1 {
m.cursor++
if m.cursor >= m.scrollOffset+(m.height-4) {
m.scrollOffset++
}
}
case " ":
if m.cursor < len(m.tags) {
tag := m.tags[m.cursor]
if m.selectedTags[tag] {
delete(m.selectedTags, tag)
} else {
m.selectedTags[tag] = true
}
}
case "right":
// Select all
for _, t := range m.tags {
m.selectedTags[t] = true
}
case "enter":
// Apply filter
m.applyFilter()
m.view = ViewList
m.cursor = 0
m.scrollOffset = 0
case "esc":
// Cancel
m.selectedTags = m.backupSelectedTags
m.view = ViewList
m.cursor = 0
m.scrollOffset = 0
}
}
return m, nil
}

func (m *Model) applyFilter() {
if len(m.selectedTags) == 0 {
m.filteredArticles = m.articles
return
}

var filtered []readings.Article
for _, a := range m.articles {
match := false
for _, t := range a.Tags {
if m.selectedTags[t] {
match = true
break
}
}
if match {
filtered = append(filtered, a)
}
}
m.filteredArticles = filtered
}

func openUrl(url string) tea.Cmd {
return func() tea.Msg {
var cmd string
var args []string

switch runtime.GOOS {
case "windows":
cmd = "cmd"
args = []string{"/c", "start"}
case "darwin":
cmd = "open"
default: // "linux", "freebsd", "openbsd", "netbsd"
cmd = "xdg-open"
}
args = append(args, url)
c := exec.Command(cmd, args...)
return c.Start()
}
}

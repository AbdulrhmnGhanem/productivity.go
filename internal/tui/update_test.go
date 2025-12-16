package tui

import (
"testing"

tea "github.com/charmbracelet/bubbletea"
"github.com/stretchr/testify/assert"
"productivity.go/internal/readings"
)

func TestUpdate_Quit(t *testing.T) {
	m := Model{}
	tests := []struct {
		key string
	}{
		{"q"},
		{"ctrl+c"},
	}

	for _, tt := range tests {
		t.Run(tt.key, func(t *testing.T) {
_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(tt.key)})
if tt.key == "ctrl+c" {
_, cmd = m.Update(tea.KeyMsg{Type: tea.KeyCtrlC})
}
assert.NotNil(t, cmd)
})
	}
}

func TestUpdate_WindowSize(t *testing.T) {
	m := Model{}
	msg := tea.WindowSizeMsg{Width: 100, Height: 50}
	newM, _ := m.Update(msg)
	newModel := newM.(Model)
	assert.Equal(t, 100, newModel.width)
	assert.Equal(t, 50, newModel.height)
}

func TestUpdate_Filter(t *testing.T) {
	articles := []readings.Article{
		{Title: "Go Article", Tags: []string{"go"}},
		{Title: "Rust Article", Tags: []string{"rust"}},
		{Title: "Both Article", Tags: []string{"go", "rust"}},
	}
	tags := []string{"go", "rust"}
	m := Model{
		articles:         articles,
		filteredArticles: articles,
		tags:             tags,
		selectedTags:     make(map[string]bool),
		view:             ViewList,
	}

	// 1. Enter Filter Mode
	newM, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("/")})
	model := newM.(Model)
	assert.Equal(t, ViewFilter, model.view)
	assert.NotNil(t, model.backupSelectedTags)

	// 2. Toggle "go" (cursor at 0)
	newM, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(" ")})
	model = newM.(Model)
	assert.True(t, model.selectedTags["go"])

	// 3. Move down to "rust"
	newM, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")})
	model = newM.(Model)
	assert.Equal(t, 1, model.cursor)

	// 4. Toggle "rust"
	newM, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(" ")})
	model = newM.(Model)
	assert.True(t, model.selectedTags["rust"])

	// 5. Apply Filter
	newM, _ = model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	model = newM.(Model)
	assert.Equal(t, ViewList, model.view)
	// Should have all 3 articles since they all have either go or rust
	assert.Equal(t, 3, len(model.filteredArticles))

	// Test filtering subset
	m.selectedTags = map[string]bool{"go": true}
	m.view = ViewFilter
	newM, _ = m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	model = newM.(Model)
	// Should have "Go Article" and "Both Article"
	assert.Equal(t, 2, len(model.filteredArticles))
	assert.Equal(t, "Go Article", model.filteredArticles[0].Title)
	assert.Equal(t, "Both Article", model.filteredArticles[1].Title)

	// Test Cancel
	m.selectedTags = map[string]bool{"go": true}
	m.backupSelectedTags = map[string]bool{"go": true}
	m.view = ViewFilter
	// Toggle rust (so now go+rust)
	m.cursor = 1 // rust
	newM, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(" ")})
	model = newM.(Model)
	assert.True(t, model.selectedTags["rust"])
	// Cancel
	newM, _ = model.Update(tea.KeyMsg{Type: tea.KeyEsc})
	model = newM.(Model)
	assert.Equal(t, ViewList, model.view)
	assert.False(t, model.selectedTags["rust"]) // Should be reverted
	assert.True(t, model.selectedTags["go"])    // Should be kept
}

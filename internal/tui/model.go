package tui

import (
	"context"
	"math/rand"
	"sort"

	tea "github.com/charmbracelet/bubbletea"
	"productivity.go/internal/readings"
)

// ViewState represents the current active view in the TUI.
type ViewState int

const (
ViewList ViewState = iota
ViewDetail
ViewFilter
)

// Model holds the application state.
type Model struct {
	// Data
	articles           []readings.Article
	filteredArticles   []readings.Article
	tags               []string
	selectedTags       map[string]bool
	backupSelectedTags map[string]bool // To restore on Cancel

	// State
	view         ViewState
	cursor       int // Index of selected item in the current list
	scrollOffset int // For scrolling long lists
	width        int
	height       int
	inputBuffer  string // For Vim-style numeric commands
	statusMessage string

	// Services
	svc *readings.Service
}

type ClearStatusMsg struct{}
type StatusMsg string


// InitTUI initializes the TUI model with data.
func InitTUI(svc *readings.Service) (Model, error) {
	articles, err := svc.GetAll(context.Background())
	if err != nil {
		return Model{}, err
	}

	// Shuffle articles
	rand.Shuffle(len(articles), func(i, j int) {
		articles[i], articles[j] = articles[j], articles[i]
	})

	// Extract unique tags
	tagMap := make(map[string]bool)
	for _, a := range articles {
		for _, t := range a.Tags {
			tagMap[t] = true
		}
	}
	var tags []string
	for t := range tagMap {
		tags = append(tags, t)
	}
	sort.Strings(tags)

	return Model{
		articles:         articles,
		filteredArticles: articles, // Initially show all
		tags:             tags,
		selectedTags:     make(map[string]bool),
		view:             ViewList,
		svc:              svc,
	}, nil
}

func (m Model) Init() tea.Cmd {
	return nil
}

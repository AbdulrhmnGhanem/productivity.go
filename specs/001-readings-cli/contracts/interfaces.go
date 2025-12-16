package contracts

import (
	"context"
	"time"
)

// Article represents a reading item.
type Article struct {
	ID        string
	Title     string
	URL       string
	Tags      []string
	FetchedAt time.Time
}

// Repository defines the interface for local storage (DuckDB).
type Repository interface {
	// SaveUpsert saves articles to the local cache, updating existing ones.
	SaveUpsert(ctx context.Context, articles []Article) error

	// GetRandom returns 'count' random articles, optionally filtered by tag.
	GetRandom(ctx context.Context, count int, tag string) ([]Article, error)

	// GetAll returns all articles (fallback if count < requested).
	GetAll(ctx context.Context) ([]Article, error)

	// Close closes the storage connection.
	Close() error
}

// NotionClient defines the interface for fetching data from Notion.
type NotionClient interface {
	// FetchArticles returns all non-done articles from the configured database.
	FetchArticles(ctx context.Context) ([]Article, error)
}

// Syncer defines the interface for background synchronization.
type Syncer interface {
	// Sync performs the synchronization process (fetch -> save).
	Sync(ctx context.Context) error

	// TriggerBackgroundSync spawns a detached process to run Sync.
	TriggerBackgroundSync() error
}

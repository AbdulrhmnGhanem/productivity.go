package readings

import (
	"context"
	"time"
)

// Article represents a reading item.
type Article struct {
	ID        string    `db:"id"`
	Title     string    `db:"title"`
	URL       string    `db:"url"`
	Tags      []string  `db:"-"` // Handled via custom scanner/valuer or JSON string in DB
	FetchedAt time.Time `db:"fetched_at"`
}

// Repository defines the interface for local storage.
type Repository interface {
	// SaveUpsert saves articles to the local cache, updating existing ones.
	SaveUpsert(ctx context.Context, articles []Article) error

	// GetRandom returns 'count' random articles, optionally filtered by tag.
	GetRandom(ctx context.Context, count int, tag string) ([]Article, error)

	// GetAll returns all articles.
	GetAll(ctx context.Context) ([]Article, error)

	// Close closes the storage connection.
	Close() error
}

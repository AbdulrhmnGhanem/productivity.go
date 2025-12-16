package storage

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
	"productivity.go/internal/readings"
)

const (
	DBFileName = "readings.sqlite"
	DBDirName  = "productivity.go"
)

type SQLite struct {
	db *sql.DB
}

// NewSQLite creates a new SQLite storage instance.
func NewSQLite() (*SQLite, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get user home dir: %w", err)
	}

	dbDir := filepath.Join(home, ".config", DBDirName)
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create db directory: %w", err)
	}

	dbPath := filepath.Join(dbDir, DBFileName)
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open sqlite: %w", err)
	}

	store := &SQLite{db: db}
	if err := store.migrate(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to migrate db: %w", err)
	}

	return store, nil
}

func (s *SQLite) migrate() error {
	query := `
	CREATE TABLE IF NOT EXISTS articles (
		id TEXT PRIMARY KEY,
		title TEXT NOT NULL,
		url TEXT NOT NULL,
		tags TEXT, -- Stored as JSON string
		fetched_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	`
	_, err := s.db.Exec(query)
	return err
}

func (s *SQLite) Close() error {
	return s.db.Close()
}

func (s *SQLite) SaveUpsert(ctx context.Context, articles []readings.Article) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, `
		INSERT INTO articles (id, title, url, tags, fetched_at)
		VALUES (?, ?, ?, ?, ?)
		ON CONFLICT (id) DO UPDATE SET
			title = excluded.title,
			url = excluded.url,
			tags = excluded.tags,
			fetched_at = excluded.fetched_at
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, a := range articles {
		tagsJSON, err := json.Marshal(a.Tags)
		if err != nil {
			return fmt.Errorf("failed to marshal tags for article %s: %w", a.ID, err)
		}

		_, err = stmt.ExecContext(ctx, a.ID, a.Title, a.URL, string(tagsJSON), a.FetchedAt)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (s *SQLite) GetRandom(ctx context.Context, count int, tag string) ([]readings.Article, error) {
	var query string
	var args []interface{}

	if tag != "" {
		// SQLite doesn't have ILIKE by default, but modernc might support it or we use LIKE with upper/lower.
		// modernc/sqlite supports LIKE which is case-insensitive for ASCII by default.
		query = `SELECT id, title, url, tags, fetched_at FROM articles WHERE tags LIKE ? ORDER BY random() LIMIT ?`
		args = []interface{}{"%" + tag + "%", count}
	} else {
		query = `SELECT id, title, url, tags, fetched_at FROM articles ORDER BY random() LIMIT ?`
		args = []interface{}{count}
	}

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanArticles(rows)
}

func (s *SQLite) GetAll(ctx context.Context) ([]readings.Article, error) {
	query := `SELECT id, title, url, tags, fetched_at FROM articles`
	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanArticles(rows)
}

func scanArticles(rows *sql.Rows) ([]readings.Article, error) {
	var articles []readings.Article
	for rows.Next() {
		var a readings.Article
		var tagsJSON string
		if err := rows.Scan(&a.ID, &a.Title, &a.URL, &tagsJSON, &a.FetchedAt); err != nil {
			return nil, err
		}

		if tagsJSON != "" {
			if err := json.Unmarshal([]byte(tagsJSON), &a.Tags); err != nil {
				// Log error but continue? Or fail?
				// For now, empty tags
				a.Tags = []string{}
			}
		}
		articles = append(articles, a)
	}
	return articles, rows.Err()
}

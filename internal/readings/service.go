package readings

import (
	"context"
	"fmt"
)

// NotionClient defines the interface for fetching articles from Notion.
type NotionClient interface {
	FetchArticles(ctx context.Context) ([]Article, error)
	FetchCurrentWeek(ctx context.Context) (*Week, error)
	UpdateWeekReadingList(ctx context.Context, weekPageID string, readingPageIDs []string) error
}

type Service struct {
	repo        Repository
	notion      NotionClient
	currentWeek *Week
}

func NewService(repo Repository, notion NotionClient) *Service {
	return &Service{
		repo:   repo,
		notion: notion,
	}
}

func (s *Service) GetReadings(ctx context.Context, count int, tag string) ([]Article, error) {
	articles, err := s.repo.GetRandom(ctx, count, tag)
	if err != nil {
		return nil, fmt.Errorf("failed to get readings: %w", err)
	}

	// If we have enough articles, return them
	if len(articles) >= count {
		return articles, nil
	}

	// If we don't have enough (or zero), and it's a fresh start, we might want to sync.
	// For now, if we have 0, we sync.
	if len(articles) == 0 {
		// Try to sync
		if err := s.Sync(ctx); err != nil {
			return nil, fmt.Errorf("failed to sync: %w", err)
		}
		// Try again
		return s.repo.GetRandom(ctx, count, tag)
	}

	return articles, nil
}

func (s *Service) Sync(ctx context.Context) error {
	articles, err := s.notion.FetchArticles(ctx)
	if err != nil {
		return err
	}

	if err := s.repo.SaveUpsert(ctx, articles); err != nil {
		return err
	}

	return nil
}

func (s *Service) GetAll(ctx context.Context) ([]Article, error) {
	return s.repo.GetAll(ctx)
}

func (s *Service) ToggleReadingInCurrentWeek(ctx context.Context, articleID string) (bool, error) {
	if s.currentWeek == nil {
		week, err := s.notion.FetchCurrentWeek(ctx)
		if err != nil {
			return false, err
		}
		s.currentWeek = week
	}

	// Check if article is already in the list
	exists := false
	var newIDs []string
	for _, id := range s.currentWeek.ReadingListIDs {
		if id == articleID {
			exists = true
		} else {
			newIDs = append(newIDs, id)
		}
	}

	added := false
	if !exists {
		newIDs = append(newIDs, articleID)
		added = true
	}

	if err := s.notion.UpdateWeekReadingList(ctx, s.currentWeek.ID, newIDs); err != nil {
		return false, err
	}

	// Update cache
	s.currentWeek.ReadingListIDs = newIDs

	return added, nil
}

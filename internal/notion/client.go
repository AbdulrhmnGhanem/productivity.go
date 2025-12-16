package notion

import (
	"context"
	"fmt"
	"time"

	"github.com/jomei/notionapi"
	"productivity.go/internal/readings"
)

type Client struct {
	api        *notionapi.Client
	databaseID notionapi.DatabaseID
	weeksDBID  notionapi.DatabaseID
}

func NewClient(apiKey, databaseID, weeksDBID string) *Client {
	return &Client{
		api:        notionapi.NewClient(notionapi.Token(apiKey)),
		databaseID: notionapi.DatabaseID(databaseID),
		weeksDBID:  notionapi.DatabaseID(weeksDBID),
	}
}

func (c *Client) FetchArticles(ctx context.Context) ([]readings.Article, error) {
	var articles []readings.Article
	var cursor notionapi.Cursor

	for {
		req := &notionapi.DatabaseQueryRequest{
			Filter: &notionapi.PropertyFilter{
				Property: "Done",
				Checkbox: &notionapi.CheckboxFilterCondition{
					DoesNotEqual: true,
				},
			},
			StartCursor: cursor,
		}

		resp, err := c.api.Database.Query(ctx, c.databaseID, req)
		if err != nil {
			return nil, fmt.Errorf("failed to query notion database: %w", err)
		}

		for _, page := range resp.Results {
			article, err := parsePage(page)
			if err != nil {
				// Log error but continue? For now, let's skip malformed pages
				continue
			}
			articles = append(articles, article)
		}

		if !resp.HasMore {
			break
		}
		cursor = resp.NextCursor
	}

	return articles, nil
}

func parsePage(page notionapi.Page) (readings.Article, error) {
	var title string
	if prop, ok := page.Properties["Name"].(*notionapi.TitleProperty); ok {
		for _, t := range prop.Title {
			title += t.PlainText
		}
	}

	var url string
	if prop, ok := page.Properties["URL"].(*notionapi.URLProperty); ok {
		url = prop.URL
	}

	var tags []string
	if prop, ok := page.Properties["Tags"].(*notionapi.MultiSelectProperty); ok {
		for _, option := range prop.MultiSelect {
			tags = append(tags, option.Name)
		}
	}

	// If URL is empty, maybe use the page URL?
	if url == "" {
		url = page.URL
	}

	return readings.Article{
		ID:        page.ID.String(),
		Title:     title,
		URL:       url,
		Tags:      tags,
		FetchedAt: time.Now(),
	}, nil
}

func (c *Client) FetchCurrentWeek(ctx context.Context) (*readings.Week, error) {
	now := time.Now()
	// We sort by Name descending to get recent weeks first (assuming naming convention follows date).
	req := &notionapi.DatabaseQueryRequest{
		Sorts: []notionapi.SortObject{
			{
				Property:  "Name",
				Direction: notionapi.SortOrderDESC,
			},
		},
		PageSize: 20,
	}

	resp, err := c.api.Database.Query(ctx, c.weeksDBID, req)
	if err != nil {
		return nil, fmt.Errorf("failed to query weeks database: %w", err)
	}

	for _, page := range resp.Results {
		if prop, ok := page.Properties["🗓️ Span"].(*notionapi.DateProperty); ok {
			if prop.Date.Start != nil {
				start := time.Time(*prop.Date.Start)
				var end time.Time
				if prop.Date.End != nil {
					end = time.Time(*prop.Date.End)
				} else {
					// If no end date, assume it covers the start day
					end = start.Add(24 * time.Hour)
				}

				// Adjust end to be end of the day if it's 00:00:00
				if end.Hour() == 0 && end.Minute() == 0 && end.Second() == 0 {
					end = end.Add(24 * time.Hour).Add(-1 * time.Second)
				}

				if !now.Before(start) && !now.After(end) {
					return parseWeek(page)
				}
			}
		}
	}

	return nil, fmt.Errorf("no current week found in Notion")
}

func getKeys(m map[string]notionapi.Property) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func (c *Client) UpdateWeekReadingList(ctx context.Context, weekPageID string, readingPageIDs []string) error {
	relations := make([]notionapi.Relation, len(readingPageIDs))
	for i, id := range readingPageIDs {
		relations[i] = notionapi.Relation{ID: notionapi.PageID(id)}
	}

	params := &notionapi.PageUpdateRequest{
		Properties: notionapi.Properties{
			"📑 Reading List": notionapi.RelationProperty{
				Relation: relations,
			},
		},
	}

	_, err := c.api.Page.Update(ctx, notionapi.PageID(weekPageID), params)
	if err != nil {
		return fmt.Errorf("failed to update week reading list: %w", err)
	}
	return nil
}

func parseWeek(page notionapi.Page) (*readings.Week, error) {
	var readingListIDs []string
	if prop, ok := page.Properties["📑 Reading List"].(*notionapi.RelationProperty); ok {
		for _, rel := range prop.Relation {
			readingListIDs = append(readingListIDs, rel.ID.String())
		}
	}

	return &readings.Week{
		ID:             page.ID.String(),
		ReadingListIDs: readingListIDs,
	}, nil
}

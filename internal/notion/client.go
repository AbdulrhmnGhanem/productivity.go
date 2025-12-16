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
}

func NewClient(apiKey, databaseID string) *Client {
	return &Client{
		api:        notionapi.NewClient(notionapi.Token(apiKey)),
		databaseID: notionapi.DatabaseID(databaseID),
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

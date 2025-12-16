# Data Model: Readings CLI

## Entities

### Article

Represents a single reading item fetched from Notion.

| Field        | Type                | Description                                 |
| ------------ | ------------------- | ------------------------------------------- |
| `id`         | String (UUID)       | Unique identifier from Notion. Primary Key. |
| `title`      | String              | Title of the article.                       |
| `url`        | String              | URL of the article.                         |
| `tags`       | String (JSON Array) | List of tags associated with the article.   |
| `fetched_at` | Timestamp           | When the article was last synced.           |

## Storage Schema (DuckDB)

```sql
CREATE TABLE IF NOT EXISTS articles (
    id VARCHAR PRIMARY KEY,
    title VARCHAR NOT NULL,
    url VARCHAR NOT NULL,
    tags VARCHAR, -- Stored as JSON string e.g. "['tech', 'productivity']"
    fetched_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

## Go Struct

```go
type Article struct {
    ID        string    `db:"id"`
    Title     string    `db:"title"`
    URL       string    `db:"url"`
    Tags      []string  `db:"-"` // Handled via custom scanner/valuer
    FetchedAt time.Time `db:"fetched_at"`
}
```

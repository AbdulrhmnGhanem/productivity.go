# Implementation Plan - Link Reading to Week

## Proposed Changes

### Configuration

- **File**: `internal/config/config.go`
- **Change**: Add `NotionWeeksDBID` to `Config` struct.
- **Change**: Update `loadViperConfig` to read `notion_weeks_db_id` from the config file.

### Notion Client

- **File**: `internal/notion/client.go`
- **Change**: Add `FetchCurrentWeek(ctx context.Context) (*Week, error)` method.
    - Query the Weeks database using `NotionWeeksDBID`.
    - Filter by "Span" date property containing the current date.
    - Return a `Week` struct (or similar) containing the Page ID and current "Reading List" IDs.
- **Change**: Add `UpdateWeekReadingList(ctx context.Context, weekPageID string, readingPageIDs []string) error` method.
    - Update the "Reading List" relation property of the specified week page.

### Service Layer

- **File**: `internal/readings/service.go`
- **Change**: Update `NotionClient` interface to include new methods.
- **Change**: Add `ToggleReadingInCurrentWeek(ctx context.Context, articleID string) (bool, error)` method.
    - Fetch current week.
    - Check if `articleID` exists in the week's reading list.
    - Toggle presence (add if missing, remove if present).
    - Call `UpdateWeekReadingList`.
    - Return `true` if added, `false` if removed.

### TUI

- **File**: `internal/tui/model.go`
- **Change**: Add `statusMessage` string and `statusTimer` *time.Timer (or similar mechanism) to `Model`.
- **Change**: Define `ClearStatusMsg` tea.Msg.

- **File**: `internal/tui/update.go`
- **Change**: Handle `Enter` key in `ViewList`.
    - Trigger `ToggleReadingInCurrentWeek` (async via `tea.Cmd`).
    - Update `statusMessage` on success/failure.
    - Schedule `ClearStatusMsg` after 2 seconds.
- **Change**: Handle `ClearStatusMsg` to clear the status line.

- **File**: `internal/tui/view.go`
- **Change**: Render `statusMessage` at the bottom of the view (ephemeral status line).

## Verification Plan

### Automated Tests
- **Unit Tests**:
    - `internal/readings/service_test.go`: Mock `NotionClient` and test `ToggleReadingInCurrentWeek` logic (add vs remove).
    - `internal/tui/update_test.go`: Test `Enter` key triggers the correct command and state change.

### Manual Verification
- **Setup**: Configure `notion_weeks_db_id` in `productivity.go.toml`.
- **Scenario 1**: Select an article, press `Enter`. Verify "Added to reading list" message. Check Notion.
- **Scenario 2**: Press `Enter` again on the same article. Verify "Removed from reading list" message. Check Notion.
- **Scenario 3**: Press `Enter` on a different article. Verify it's added.

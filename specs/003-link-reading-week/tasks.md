# Tasks: Link Reading to Week

**Feature Branch**: `003-link-reading-week`
**Status**: Completed

## Phase 1: Setup

- [x] T001 Update `Config` struct in [internal/config/config.go](internal/config/config.go) to include `NotionWeeksDBID`
- [x] T002 Update `loadViperConfig` in [internal/config/config.go](internal/config/config.go) to read `notion_weeks_db_id`

## Phase 2: Foundational

- [x] T003 [P] Implement `FetchCurrentWeek` in [internal/notion/client.go](internal/notion/client.go) to query Weeks DB
- [x] T004 [P] Implement `UpdateWeekReadingList` in [internal/notion/client.go](internal/notion/client.go) to update page relations
- [x] T005 Update `NotionClient` interface in [internal/readings/service.go](internal/readings/service.go) to include new methods
- [x] T006 [US1] Implement `ToggleReadingInCurrentWeek` in [internal/readings/service.go](internal/readings/service.go)

## Phase 3: User Story 1 - Link Article to Current Week

**Goal**: Enable users to add/remove articles from the current week's reading list using the Enter key.

- [x] T007 [US1] Add `statusMessage` field and `ClearStatusMsg` type to `Model` in [internal/tui/model.go](internal/tui/model.go)
- [x] T008 [US1] Update `Update` function in [internal/tui/update.go](internal/tui/update.go) to handle `Enter` key in `ViewList`
- [x] T009 [US1] Implement `ToggleReadingInCurrentWeek` call and status message handling in [internal/tui/update.go](internal/tui/update.go)
- [x] T010 [US1] Update `View` in [internal/tui/view.go](internal/tui/view.go) to render the ephemeral status line

## Phase 4: Polish & Cross-Cutting

- [x] T011 Verify configuration loading and error handling for missing DB ID
- [x] T012 Manual verification of the full flow with Notion

## Dependencies

- Phase 2 depends on Phase 1
- Phase 3 depends on Phase 2

## Implementation Strategy

1.  **Configuration**: Start by ensuring we can read the new config value.
2.  **Notion Logic**: Implement the core logic for interacting with the Weeks database.
3.  **Service Logic**: Wrap the Notion calls in a convenient service method that handles the toggle logic.
4.  **UI Integration**: Finally, wire it up to the TUI with the new keybinding and status feedback.

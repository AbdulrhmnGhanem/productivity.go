# Feature Specification: Readings CLI with Notion Sync

**Feature Branch**: `001-readings-cli`  
**Created**: 2025-12-16  
**Status**: Draft  
**Input**: User description: "the first cli should be named `readings` it will fetch articles from my reading list (which is located at in a notion db https://www.notion.so/abdulrhmnghanem/a0e3e448792a4aa59f0d4576333457e9?v=f291b0e4b2f64b7d818fe996318ecdf1) and it should let me select the 7 articles for the week. By default, it selects 7 randomly, but I can select on of the tags, so it limits the selection to those. To avoid delays we have to keep a local version of the db and run a sync process in the background so the next time I ask for articles it has them fetched"

## User Scenarios & Testing _(mandatory)_

<!--
  IMPORTANT: User stories should be PRIORITIZED as user journeys ordered by importance.
  Each user story/journey must be INDEPENDENTLY TESTABLE - meaning if you implement just ONE of them,
  you should still have a viable MVP (Minimum Viable Product) that delivers value.

  Assign priorities (P1, P2, P3, etc.) to each story, where P1 is the most critical.
  Think of each story as a standalone slice of functionality that can be:
  - Developed independently
  - Tested independently
  - Deployed independently
  - Demonstrated to users independently
-->

### User Story 1 - Get Weekly Readings (Priority: P1)

As a user, I want to get a list of 7 random articles from my reading list so that I can decide what to read this week without manual selection.

**Why this priority**: This is the core value proposition of the tool.

**Independent Test**: Can be fully tested by running the command and verifying 7 articles are returned from the local cache.

**Acceptance Scenarios**:

1. **Given** the local cache contains more than 7 articles, **When** the user runs `readings`, **Then** the system displays a list of 7 randomly selected articles.
2. **Given** the local cache contains fewer than 7 articles, **When** the user runs `readings`, **Then** the system displays all available articles.
3. **Given** the local cache is empty, **When** the user runs `readings`, **Then** the system attempts to fetch data from Notion (or informs the user to sync).

---

### User Story 2 - Filter by Tag (Priority: P2)

As a user, I want to limit the selection of articles to a specific tag so that I can focus on a particular topic.

**Why this priority**: Allows for focused reading sessions.

**Independent Test**: Run the command with a tag argument and verify all returned articles have that tag.

**Acceptance Scenarios**:

1. **Given** the local cache has articles with the tag "tech", **When** the user runs `readings --tag "tech"`, **Then** the system displays 7 randomly selected articles that have the "tech" tag.
2. **Given** no articles match the provided tag, **When** the user runs `readings --tag "invalid"`, **Then** the system displays a message indicating no articles were found for that tag.

---

### User Story 3 - Background Sync (Priority: P3)

As a user, I want the system to synchronize with Notion in the background so that my local list is always up-to-date without making me wait during command execution.

**Why this priority**: Ensures a smooth and fast user experience by avoiding network latency during the main interaction.

**Independent Test**: Add an article to Notion, trigger the background sync, and verify the article appears in the local cache.

**Acceptance Scenarios**:

1. **Given** the user runs a `readings` command, **When** the command completes, **Then** a background synchronization process is triggered to update the local cache from Notion.
2. **Given** new articles have been added to the Notion database, **When** the background sync completes, **Then** the new articles are available in the local cache for the next execution.

---

### User Story 4 - Initial Setup (Priority: P0)

As a user, I want to easily configure the tool with my Notion credentials so that I can start using it without manually editing config files.

**Why this priority**: Prerequisite for using the tool.

**Independent Test**: Run the setup command on a fresh install and verify files are created correctly.

**Acceptance Scenarios**:

1. **Given** the tool is not configured, **When** the user runs `readings setup`, **Then** the system prompts for the Notion API Key and Database ID.
2. **Given** valid inputs, **When** the setup completes, **Then** the API Key is saved to `.netrc` and the Database ID is saved to `productivity.go.toml`.
3. **Given** the configuration exists, **When** the user runs `readings`, **Then** the system uses the stored credentials to fetch data.

---

### Edge Cases

- What happens when the Notion API is unreachable during sync? (System should log the error and keep using the existing cache).
- What happens when the Notion database structure changes? (System should handle missing fields gracefully).
- How does the system handle duplicate articles? (System should identify unique articles by ID or URL).
- What happens if the user has no internet connection? (System should function using the local cache).

## Requirements _(mandatory)_

### Functional Requirements

- **FR-001**: System MUST provide a CLI command `readings` that outputs a list of articles.
- **FR-002**: System MUST fetch article data (Title, URL, Tags) from the specified Notion Database, filtering for articles where the "Done" property is false/unchecked.
- **FR-003**: System MUST store article data in a local cache/database to ensure instant retrieval.
- **FR-004**: System MUST select a fresh set of 7 random articles from the cache on each execution (no persistence of weekly selection).
- **FR-005**: System MUST allow the user to filter the selection by a specific tag using a flag (e.g., `--tag`).
- **FR-006**: System MUST perform synchronization with Notion by spawning a detached background process after the main command execution, ensuring the user is not blocked.
- **FR-007**: System MUST return all matching articles if fewer than 7 are available.
- **FR-008**: System MUST provide an interactive setup command (e.g., `readings setup`) to prompt the user for configuration details.
- **FR-009**: System MUST save sensitive credentials (Notion API Key) to the user's `.netrc` file.
- **FR-010**: System MUST save non-sensitive configuration (Notion Database ID) to a `productivity.go.toml` file.
- **FR-011**: System MUST treat the Notion "Done" status as read-only; the CLI MUST NOT modify article status in Notion.

### Key Entities _(include if feature involves data)_

- **Article**: Represents a reading item. Attributes: Title, URL, Tags, Notion ID.
- **Tag**: A category label associated with an article.

## Success Criteria _(mandatory)_

### Measurable Outcomes

- **SC-001**: The `readings` command returns results in under 1 second when using the local cache.
- **SC-002**: Background synchronization successfully updates the local cache with new articles from Notion within 5 minutes of execution (dependent on network).
- **SC-003**: Users can successfully filter articles by any existing tag found in the Notion database.

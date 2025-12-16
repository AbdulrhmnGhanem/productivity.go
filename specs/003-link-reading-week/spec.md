# Feature Specification: Link Reading to Week

**Feature Branch**: `003-link-reading-week`
**Created**: 2025-12-16
**Status**: Draft
**Input**: User description: "I have another notion database that tracks my weekly work. I want to update the readings cli that when I press enter on an article it adds it to the latest week reading list instead of openning the article details. now we have two db configured in the cli so we have to track them seperately. Alos my weeks database changes every 12 months so I will need to set the db id"

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Link Article to Current Week (Priority: P1)

As a user, I want to quickly add a reading item to my current week's plan by pressing Enter, so that I can organize my reading list without leaving the CLI or navigating Notion manually.

**Why this priority**: This is the core functionality requested to streamline the workflow.

**Independent Test**: Can be tested by mocking the Notion API responses for the Weeks database and verifying the update call when Enter is pressed on an article.

**Acceptance Scenarios**:

1. **Given** the CLI is open and a list of articles is displayed, **When** I navigate to an article and press `Enter`, **Then** the system identifies the current week from the Weeks database.
2. **Given** the article is NOT in the current week's "Reading List", **When** I press `Enter`, **Then** the system adds it and displays a success message "Added to reading list".
3. **Given** the article IS already in the current week's "Reading List", **When** I press `Enter`, **Then** the system removes it and displays a success message "Removed from reading list".
4. **Given** no matching week is found for the current date, **When** I press `Enter`, **Then** an error message is displayed.
5. **Given** the CLI is open, **When** I press `Enter`, **Then** the article details view does NOT open.

---

### User Story 2 - Configure Weeks Database (Priority: P2)

As a user, I want to configure the Notion Database ID for my Weeks database, so that I can update it when I switch to a new yearly database.

**Why this priority**: Essential for the feature to work across different years/databases.

**Independent Test**: Can be tested by changing the configuration and verifying the application uses the new ID for queries.

**Acceptance Scenarios**:

1. **Given** a new Weeks database ID, **When** I update the configuration (e.g., config file or environment variable), **Then** the CLI uses this new ID for subsequent operations.

### Edge Cases

- **No Current Week Found**: If no week entry covers the current date, the system displays an error message "No current week found in Notion".
- **Toggle Behavior**: If the article is already in the "Reading List", pressing `Enter` removes it. If it is not present, it adds it.
- **Network Failure**: If the Notion API call fails, the system displays an error message and retains the UI state.
- **Missing Configuration**: If `NOTION_WEEKS_DB_ID` is not set, the system prompts the user to configure it or displays an error.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST allow configuration of the Weeks Notion Database ID (e.g., via `NOTION_WEEKS_DB_ID`).
- **FR-002**: System MUST identify the "current week" entry in the Weeks database.
    - *Assumption*: "Current week" is defined as the entry where the current date falls within the entry's date range property.
- **FR-003**: System MUST check if the selected article is already in the "Reading List" of the identified Week.
- **FR-004**: System MUST toggle the article's presence in the "Reading List":
    - If present: Remove it.
    - If absent: Add it.
- **FR-005**: System MUST bind the `Enter` key in the article list view to trigger this toggle action.
- **FR-006**: System MUST disable the previous behavior of opening article details on `Enter`.
- **FR-007**: System MUST provide visual feedback via an ephemeral status line at the bottom of the screen that disappears automatically.

### Key Entities *(include if feature involves data)*

- **Week**: Represents a weekly planning entry in Notion.
    - **Date Range**: A date property defining the start and end of the week.
    - **Reading List**: A relation property linking to items in the Readings database.
- **Reading**: Represents an article or book in the Readings database.

## Success Criteria *(mandatory)*

- **Efficiency**: Users can add an article to the current week with a single keystroke (Enter).
- **Accuracy**: The system correctly identifies the week containing the current date 100% of the time when a valid week entry exists.
- **Configurability**: The Weeks Database ID can be updated by the user without code changes.
- **Feedback**: Users receive immediate confirmation (within < 2 seconds, network permitting) that the action succeeded.

## Assumptions

- The "Weeks" database has a date property named "Span" that covers the current date.
- The "Weeks" database has a "Reading List" property that is a Relation to the Readings database.
- The user has appropriate permissions to read/write to the Weeks database.

## Clarifications

### Session 2025-12-16
- Q: How should the system handle articles already in the reading list? → A: Toggle (remove if present, add if absent).

- Q: Are the property names for the Weeks database fixed or configurable? → A: Fixed names ("Span" for date, "Reading List" for relation).
- Q: How should success/failure feedback be presented? → A: Ephemeral Status Line (disappears automatically).
- Q: Should the user be able to select a different week than the current one? → A: No (Current week only).

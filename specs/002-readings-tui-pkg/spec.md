# Feature Specification: Readings TUI & Packaging

**Feature Branch**: `002-readings-tui-pkg`
**Created**: 2025-12-16
**Status**: Draft
**Input**: User description: "currently the readings cli doesn't use a bubble tea tui to display the articles also I want to be able to select the tags interactively before filtering for selection use space to select one or right arrow to select all | also I need a github ci job to build a packages that can be installed on manjaro and ubuntu add the installation instructions to the readme (keep the readme as close as possible to man pages don't abuse emojis)"

## Clarifications

### Session 2025-12-16

- Q: Key binding to enter filter mode? → A: Press '/' (Standard search key)
- Q: Action when pressing Enter on an article? → A: Display article metadata (title, tags, URL) in a TUI detail view
- Q: Key binding to go back or exit? → A: Press `q` or `Esc` to return to the list (or exit if at root)
- Q: Action to confirm tag selection? → A: Press `Enter` to apply the filter and return to the list
- Q: Action to open URL? → A: Press 'Enter' in the detail view

## User Scenarios & Testing _(mandatory)_

### User Story 1 - Interactive Article Browsing (Priority: P1)

As a user, I want to browse my reading list in a terminal user interface (TUI) so that I can navigate through articles more efficiently than using standard CLI output.

**Why this priority**: This is the core visual upgrade requested, moving from static text to an interactive interface.

**Independent Test**: Can be fully tested by launching the application and verifying the TUI renders the list of articles and allows navigation.

**Acceptance Scenarios**:

1. **Given** the application is launched without arguments, **When** the TUI loads, **Then** a list of articles is displayed.
2. **Given** the article list is displayed, **When** I use up/down arrow keys, **Then** the selection highlight moves accordingly.
3. **Given** an article is selected, **When** I press Enter, **Then** the article details are shown.

---

### User Story 2 - Interactive Tag Filtering (Priority: P2)

As a user, I want to filter articles by selecting tags from an interactive list so that I can quickly find relevant content without typing tag names manually.

**Why this priority**: Enhances the usability of the TUI by adding powerful filtering capabilities.

**Independent Test**: Can be tested by entering the tag selection mode and verifying filtering behavior.

**Acceptance Scenarios**:

1. **Given** I am in the TUI, **When** I trigger the filter mode, **Then** a list of available tags is displayed.
2. **Given** the tag list is displayed, **When** I press Space on a tag, **Then** that tag is toggled (selected/deselected).
3. **Given** the tag list is displayed, **When** I press Right Arrow, **Then** all visible tags are selected.
4. **Given** tags are selected, **When** I confirm the selection, **Then** the main article list updates to show only articles matching the selected tags.

---

### User Story 3 - Linux Package Installation (Priority: P3)

As a Linux user (Ubuntu or Manjaro), I want to install the application using native package managers so that I can easily manage the software and its updates.

**Why this priority**: Simplifies distribution and adoption for Linux users.

**Independent Test**: Can be tested by downloading the generated artifacts from CI and attempting installation on the respective OS.

**Acceptance Scenarios**:

1. **Given** a successful CI build, **When** I check the artifacts, **Then** a `.deb` file (for Ubuntu) and a package file for Manjaro are present.
2. **Given** the `.deb` file, **When** I run the installation command on Ubuntu, **Then** the application is installed correctly.
3. **Given** the Manjaro package, **When** I run the installation command on Manjaro, **Then** the application is installed correctly.
4. **Given** the README file, **When** I read the installation section, **Then** it provides clear, text-based instructions for both platforms without excessive emojis.

### Edge Cases

- What happens when the terminal window is too small? (Should display a warning or resize gracefully)
- What happens if there are no articles or tags? (Should display an empty state message)
- What happens if the CI build fails due to dependency issues? (Pipeline should report failure)

## Requirements _(mandatory)_

### Functional Requirements

- **FR-001**: The system MUST provide a TUI mode for browsing articles.
- **FR-002**: The TUI MUST display a scrollable list of reading articles.
- **FR-003**: The system MUST allow users to enter a "Tag Selection" mode by pressing the `/` key.
- **FR-004**: In Tag Selection mode, the user MUST be able to toggle individual tags using the `Space` key.
- **FR-005**: In Tag Selection mode, the user MUST be able to select all tags using the `Right Arrow` key.
- **FR-006**: The system MUST apply the selected tag filter and return to the article list when the user presses `Enter`.
- **FR-007**: The project MUST include a CI/CD workflow to build the application.
- **FR-008**: The CI workflow MUST build an installable package for Ubuntu (e.g., `.deb`).
- **FR-009**: The CI workflow MUST build an installable package for Manjaro (e.g., `.pkg.tar.zst`).
- **FR-010**: The README MUST be updated with installation instructions for the generated packages.
- **FR-011**: The README documentation style MUST be formal ("man-page" style) and avoid the use of emojis.
- **FR-012**: When an article is selected (Enter key), the system MUST display a detail view showing the title, URL, and tags within the TUI.
- **FR-013**: The system MUST allow users to return to the previous view or exit the application by pressing `q` or `Esc`.
- **FR-014**: The system MUST open the article URL in the default web browser when the user presses `Enter` while in the detail view.

### Success Criteria

- **SC-001**: Users can navigate the article list using keyboard controls.
- **SC-002**: Users can filter articles by selecting multiple tags via the TUI.
- **SC-003**: CI pipeline successfully produces valid installable packages for Ubuntu and Manjaro.
- **SC-004**: Installation instructions are verified to work on fresh instances of the target operating systems.

### Key Entities

- **Article**: Represents a reading item with title, url, and tags.
- **Tag**: A label associated with an article used for filtering.

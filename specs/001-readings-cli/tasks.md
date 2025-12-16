# Tasks: Readings CLI with Notion Sync

**Input**: Design documents from `/specs/001-readings-cli/`
**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Project initialization and basic structure

- [x] T001 Create project structure (`cmd/`, `internal/` packages)
- [x] T002 Initialize Go module and install dependencies (`cobra`, `viper`, `bubbletea`, `notionapi`, `go-duckdb`)

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core infrastructure that MUST be complete before ANY user story can be implemented

- [x] T003 Implement Configuration (`internal/config`): Load `.netrc` and `productivity.go.toml`
- [x] T004 Implement DuckDB Storage (`internal/storage`): Schema migration, `SaveUpsert`, `GetRandom`, `GetAll`
- [x] T005 Implement Notion Client (`internal/notion`): `FetchArticles` with "Done" checkbox filter

**Checkpoint**: Foundation ready - user story implementation can now begin

## Phase 3: User Story 4 - Initial Setup (Priority: P0)

**Goal**: Allow user to configure credentials interactively

**Independent Test**: Run `readings setup` and verify config files are created

- [x] T006 [US4] Implement Setup Logic (`internal/setup`): Interactive prompt using Bubble Tea
- [x] T007 [US4] Implement Config Saving: Write to `.netrc` and `productivity.go.toml`
- [x] T008 [US4] Create `setup` command in `cmd/readings/setup.go`

**Checkpoint**: User can configure the tool

## Phase 4: User Story 1 - Get Weekly Readings (Priority: P1) 🎯 MVP

**Goal**: Display 7 random articles from cache

**Independent Test**: Run `readings` and see 7 articles

- [x] T009 [US1] Implement Readings Service (`internal/readings`): `GetReadings` logic (random selection)
- [x] T010 [US1] Implement Root Command (`cmd/readings/root.go`): Display list of articles
- [x] T011 [US1] Implement Auto-Fetch fallback: If DB empty, trigger sync (foreground)

**Checkpoint**: Core value proposition delivered

## Phase 5: User Story 2 - Filter by Tag (Priority: P2)

**Goal**: Filter articles by tag

**Independent Test**: Run `readings --tag tech`

- [x] T012 [US2] Update Readings Service: Add tag filtering to `GetReadings` and Storage
- [x] T013 [US2] Update Root Command: Add `--tag` flag and pass to service

**Checkpoint**: Filtering capability added

## Phase 6: User Story 3 - Background Sync (Priority: P3)

**Goal**: Keep cache up to date without blocking

**Independent Test**: Run command, wait, check DB for new articles

- [x] T014 [US3] Implement Sync Service (`internal/sync`): Orchestrate Notion -> DB sync
- [x] T015 [US3] Create hidden `sync` command in `cmd/readings/sync.go`
- [x] T016 [US3] Implement Detached Process Spawning: Trigger `readings sync` after root command finishes

**Checkpoint**: Feature complete

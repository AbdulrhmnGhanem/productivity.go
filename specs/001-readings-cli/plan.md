# Implementation Plan: Readings CLI with Notion Sync

**Branch**: `001-readings-cli` | **Date**: 2025-12-16 | **Spec**: [specs/001-readings-cli/spec.md](../001-readings-cli/spec.md)
**Input**: Feature specification from `/specs/001-readings-cli/spec.md`

**Note**: This template is filled in by the `/speckit.plan` command. See `.specify/templates/commands/plan.md` for the execution workflow.

## Summary

A CLI tool `readings` to fetch articles from a Notion database, cache them locally in DuckDB, and display 7 random articles for the week. Includes background synchronization and interactive setup.

## Technical Context

**Language/Version**: Go 1.23 (Latest Stable)
**Primary Dependencies**:

- `github.com/jomei/notionapi` (Notion Integration)
- `github.com/marcboeker/go-duckdb` (Local Persistence) [NEEDS CLARIFICATION: Static linking support?]
- `github.com/spf13/cobra` (CLI Framework)
- `github.com/charmbracelet/bubbletea` (TUI/Interactive Setup)
- `github.com/spf13/viper` (Configuration)
  **Storage**: DuckDB (Local file)
  **Testing**: Go `testing` package, `testify`
  **Target Platform**: Linux (CLI)
  **Project Type**: Single CLI application
  **Performance Goals**: <1s response time from cache
  **Constraints**: Offline capability, Background sync (detached process)

## Constitution Check

_GATE: Must pass before Phase 0 research. Re-check after Phase 1 design._

- [x] **Simplicity & Reliability**: CLI focus, local cache for reliability.
- [x] **Elegant TUI**: Bubble Tea for setup and list display.
- [x] **Seamless Integration**: Notion API, background sync.
- [?] **Go & Static Linking**: DuckDB driver uses CGO. [NEEDS CLARIFICATION: Can we statically link DuckDB?]
- [x] **Test-Driven Reliability**: Unit/Integration tests planned.

## Project Structure

### Documentation (this feature)

```text
specs/001-readings-cli/
├── plan.md              # This file
├── research.md          # Phase 0 output
├── data-model.md        # Phase 1 output
├── quickstart.md        # Phase 1 output
├── contracts/           # Phase 1 output
└── tasks.md             # Phase 2 output
```

### Source Code (repository root)

```text
cmd/
└── readings/           # Main entry point

internal/
├── config/             # Configuration loading (viper, .netrc)
├── notion/             # Notion API client wrapper
├── storage/            # DuckDB implementation
├── sync/               # Background sync logic
├── tui/                # Bubble Tea models
└── readings/           # Core business logic (selection, filtering)
```

**Structure Decision**: Standard Go CLI layout with `cmd/` for entry points and `internal/` for private application code.

## Complexity Tracking

> **Fill ONLY if Constitution Check has violations that must be justified**

- **DuckDB CGO**: DuckDB requires CGO, which complicates static linking. Justification: User explicitly requested DuckDB. We will attempt to statically link the C library or accept dynamic linking if strictly necessary (requires amendment or exception).

| Violation                  | Why Needed         | Simpler Alternative Rejected Because |
| -------------------------- | ------------------ | ------------------------------------ |
| [e.g., 4th project]        | [current need]     | [why 3 projects insufficient]        |
| [e.g., Repository pattern] | [specific problem] | [why direct DB access insufficient]  |

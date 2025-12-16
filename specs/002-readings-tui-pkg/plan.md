# Implementation Plan - Readings TUI & Packaging

## 1. Architecture & Design

### TUI Architecture (Bubble Tea)
We will implement a Terminal User Interface using the [Bubble Tea](https://github.com/charmbracelet/bubbletea) framework. The application will follow the Model-Update-View (ELM) architecture.

**State Management (`internal/tui/model.go`)**:
The main `Model` will hold:
- `articles`: All loaded articles.
- `filteredArticles`: Articles currently displayed.
- `tags`: List of all unique tags.
- `selectedTags`: Set of tags currently selected for filtering.
- `view`: Current view state (`ViewList`, `ViewDetail`, `ViewFilter`).
- `cursor`: Index of the currently selected item (article or tag).
- `scrollOffset`: For scrolling long lists.

**Views**:
1.  **Article List View**: Displays a scrollable list of articles.
2.  **Detail View**: Shows metadata for a selected article.
3.  **Filter View**: Overlay or separate screen to toggle tags.

**Key Interactions**:
- **Navigation**: Up/Down arrows.
- **Selection**: Enter to view details (List) or toggle tag (Filter).
- **Mode Switching**: `/` for Filter, `Esc`/`q` for Back/Exit.
- **Actions**: `Enter` in Detail view opens URL.

### Packaging & CI/CD
We will use GitHub Actions to automate the build and packaging process.
- **Tooling**: `goreleaser` will be used to automate builds and packaging.
- **Build**: Binaries MUST be statically linked (`CGO_ENABLED=0`).
- **Targets**:
    - `linux/amd64` (Ubuntu/Debian -> `.deb` via nfpm)
    - `linux/amd64` (Manjaro/Arch -> `.tar.gz` + PKGBUILD or `.pkg.tar.zst`)

## 2. Proposed Changes

### Codebase Structure
```text
internal/tui/
├── model.go       # State definitions and Init()
├── update.go      # Message handling and state transitions
├── view.go        # UI rendering logic
├── styles.go      # Lipgloss definitions
└── app.go         # Entry point to Run() the program
```

### Component Details

#### `internal/tui`
- **`InitTUI(service *readings.Service)`**: Factory to create the initial model, loading articles via the service.
- **`Update(msg tea.Msg)`**:
    - Handle `tea.KeyMsg`.
    - Switch on `m.view`.
    - Implement filtering logic when returning from Filter view.
- **`View()`**:
    - Render header/footer.
    - Render content based on `m.view`.

#### `cmd/readings`
- Modify `root.go` (or `main.go` logic) to detect if TUI mode is requested (or make it default if no args provided).
- Initialize `readings.Service`.
- Call `tui.Start(service)`.

#### `.github/workflows`
- Create `ci.yml` or `release.yml`.
- Steps:
    - Checkout code.
    - Setup Go.
    - Build binary.
    - Package for Debian (`dpkg-deb` or `nfpm`).
    - Package for Arch (manual `makepkg` or `nfpm`).
    - Upload artifacts.

#### `README.md`
- Add "Installation" section.
- Format as man-page style (headers, clear sections, no emojis).

## 3. Verification Plan

### Automated Tests
- **Unit Tests**: Comprehensive tests for TUI update logic and state transitions in `internal/tui/update_test.go` (Required by Constitution).
- **Integration**: CI pipeline success is the verification for packaging.

### Manual Verification
- Run `go run ./cmd/readings` and verify TUI launches.
- Test navigation, filtering, and opening links.
- Download CI artifacts and install on VM/Container (Ubuntu/Manjaro).

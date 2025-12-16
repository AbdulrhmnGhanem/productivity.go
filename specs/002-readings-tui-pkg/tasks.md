# Tasks: Readings TUI & Packaging

## Phase 1: Setup & Infrastructure
*Goal: Initialize the TUI structure and ensure dependencies are ready.*

- [x] T001 Create  package structure (model.go, update.go, view.go, styles.go)
- [x] T002 Define  struct in  with necessary fields (articles, cursor, state)
- [x] T003 Define  function to load articles from 

## Phase 2: Foundational TUI Implementation
*Goal: Get a basic TUI running that can display static text.*

- [x] T004 Implement  function in  to launch Bubble Tea program
- [x] T005 Implement basic  loop in  (handle Quit key)
- [x] T006 Implement basic  in  (placeholder text)
- [x] T007 Define Lipgloss styles in  (list items, selected item, title)
- [x] T008 Modify  to launch the TUI when no arguments are provided

## Phase 3: User Story 1 - Interactive Article Browsing
*Goal: Enable users to browse and view article details.*

- [x] T009 [US1] Implement  logic for  in  to render the article list
- [x] T010 [US1] Implement navigation logic (Up/Down arrows) in  to move cursor
- [x] T011 [US1] Implement  logic for  in  to show article metadata
- [x] T012 [US1] Implement state transition to  on  and back to  on /
- [x] T013 [US1] Implement "Open URL" action when pressing  inside 

## Phase 4: User Story 2 - Interactive Tag Filtering
*Goal: Enable users to filter articles by tags.*

- [x] T014 [US2] Implement logic to extract and store unique tags in  upon initialization
- [x] T015 [US2] Implement  logic for  in  to render tag list
- [x] T016 [US2] Implement mode switching: Press  to enter  from 
- [x] T017 [US2] Implement tag selection logic:  to toggle,  to select all in 
- [x] T018 [US2] Implement filter application: Press  to filter  into  and return to 

## Phase 5: User Story 3 - Linux Package Installation
*Goal: Automate packaging for Linux distributions.*

- [x] T019 [US3] Create  configuration for statically linked binaries ()
- [x] T020 [US3] Configure  to generate  packages (nfpm)
- [x] T021 [US3] Configure  to generate Arch Linux packages (PKGBUILD or archive)
- [x] T022 [US3] Create  to run Goreleaser on tag/push
- [x] T023 [US3] Update  with formal, man-page style installation instructions

## Phase 6: Polish & Cross-Cutting
*Goal: Refine UX and handle edge cases.*

- [x] T024 Handle empty list states in  (e.g., "No articles found")
- [x] T025 Ensure TUI handles terminal resizing gracefully
- [x] T026 Verify all styles are consistent and readable

# Research Findings: Readings CLI

## 1. DuckDB Static Linking in Go

**Decision**: Use `CGO_ENABLED=1` with standard build for development, and attempt fully static build for release if needed.
**Rationale**: DuckDB requires CGO. While a fully static binary is possible, it requires linking against static system libraries (`glibc`, `libstdc++`).
**Build Command**:

```bash
# Standard (Dynamic libc)
CGO_ENABLED=1 go build -o readings cmd/readings/main.go

# Fully Static (Optional, for portability)
CGO_ENABLED=1 go build -ldflags '-linkmode external -extldflags "-static"' -o readings cmd/readings/main.go
```

**Alternatives**: SQLite (pure Go drivers exist), but user requested DuckDB.

## 2. Notion API Filtering

**Decision**: Use `github.com/jomei/notionapi` with `PropertyFilter`.
**Rationale**: Library provides structured types for filters.
**Implementation Pattern**:

```go
query := &notionapi.DatabaseQueryRequest{
    Filter: notionapi.PropertyFilter{
        Property: "Done",
        Checkbox: &notionapi.CheckboxFilterCondition{
            Equals: &unchecked, // *bool pointing to false
        },
    },
}
```

## 3. Background Synchronization

**Decision**: Spawn a detached child process using `syscall.Setsid`.
**Rationale**: Ensures the sync process survives parent exit and doesn't block the terminal.
**Implementation Pattern**:

```go
cmd := exec.Command(os.Args[0], "sync") // Call self with sync command
cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
cmd.Stdin, cmd.Stdout, cmd.Stderr = nil, nil, nil
cmd.Start()
cmd.Process.Release()
```

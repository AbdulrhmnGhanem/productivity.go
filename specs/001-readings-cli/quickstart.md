# Quickstart: Readings CLI

## Prerequisites

- Go 1.23+
- Notion Integration Token (API Key)
- Notion Database ID (with "Done" checkbox property)

## Build

```bash
# Build the binary
CGO_ENABLED=1 go build -o readings cmd/readings/main.go
```

## Setup

1. Run the interactive setup:
   ```bash
   ./readings setup
   ```
   Enter your Notion API Key and Database ID when prompted.

## Usage

```bash
# Get 7 random articles
./readings

# Get 7 random articles with tag "tech"
./readings --tag tech

# Force a sync (foreground)
./readings sync
```

## Development

```bash
# Run tests
go test ./...
```

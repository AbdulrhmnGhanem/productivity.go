# Productivity Tools

This repository is a collection of productivity tools written in Go.

## Tools

### Readings

**readings** is a terminal user interface (TUI) for browsing and managing your reading list stored in Notion. It allows you to view articles, filter them by tags, and open them in your default browser.

#### Installation

**Debian/Ubuntu**

```bash
sudo dpkg -i readings_*.deb
```

**Arch Linux**

```bash
sudo pacman -U readings_*.pkg.tar.zst
```

**Binary**
Download the binary for your architecture from the Releases page and place it in your PATH.

#### Keybindings

**List View**

- **j / Down**: Move cursor down
- **k / Up**: Move cursor up
- **Enter**: View article details
- **/ (Slash)**: Open filter view
- **q / Ctrl+C**: Quit

**Detail View**

- **Enter**: Open article URL in browser
- **Esc / q**: Return to list view

**Filter View**

- **j / Down**: Move cursor down
- **k / Up**: Move cursor up
- **Space**: Toggle tag selection
- **Right**: Select all tags
- **Enter**: Apply filter
- **Esc**: Cancel filter

#### Configuration

The application requires a Notion API key and Database ID. These can be configured via environment variables or a config file.

# Productivity Tools

This repository is a collection of productivity tools written in Go.

## Tools

### Readings

**readings** is a terminal user interface (TUI) for browsing and managing my reading list stored in Notion. It allows me to view articles, filter them by tags, and open them in my default browser.

<video src="docs/assets/readings-demo.mp4" controls title="Readings Demo"></video>

#### Installation

**Debian/Ubuntu**

```bash
sudo dpkg -i readings_*.deb
```

**Arch Linux / Manjaro**

```bash
# Install from local file
sudo pacman -U readings_*.pkg.tar.zst

# Or install directly from GitHub Releases (replace version as needed)
curl -LO https://github.com/AbdulrhmnGhanem/productivity.go/releases/download/v1.0.0/readings_1.0.0_linux_amd64.pkg.tar.zst
sudo pacman -U readings_1.0.0_linux_amd64.pkg.tar.zst
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

The application requires a Notion API key and Database ID. These can be configured via environment variables or a config file using the `readings setup` command.

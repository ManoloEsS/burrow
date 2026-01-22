# Burrow

## Project Title
A terminal-based HTTP client with TUI interface for managing and executing HTTP requests.

## Overview
Burrow provides a rich terminal user interface for creating, saving, and managing HTTP requests. Built with Go and tview, it offers a keyboard-driven experience for API testing and development directly from your terminal.

## Installation

### From Source
```bash
git clone https://github.com/ManoloEsS/burrow.git
cd burrow
go build -o burrow cmd/burrow/main.go
```

### Binary Release
Download the latest binary from the [releases page](https://github.com/ManoloEsS/burrow/releases).

### Go Install
```bash
go install github.com/ManoloEsS/burrow@latest
```

## Quick Start
```bash
burrow
```

Once running, use Ctrl+N to create a new request, enter the URL, select HTTP method, and press Enter to execute.

## Usage

### Basic Commands
- Ctrl+N - Create new request
- Ctrl+S - Save request
- Ctrl+O - Open saved request
- Ctrl+R - Execute request
- Ctrl+Q - Quit

### Configuration
Burrow uses XDG Base Directory Specification for file storage and supports YAML configuration files. The application works out-of-the-box with sensible defaults.

#### Configuration Priority
1. Config file: `~/.config/burrow/config.yaml` (if exists)
2. Environment variables: `DEFAULT_PORT`, `DB_FILE`, `GOOSE_MIGRATIONS_DIR`
3. Sensible defaults: Port 8080, XDG paths for storage

#### File Locations
- **Config**: `~/.config/burrow/config.yaml`
- **Database**: `~/.local/share/burrow/burrow.db`
- **Logs**: `~/.local/state/burrow/burrow_log`
- **Server Cache**: `~/.cache/burrow/servers/`

#### Example Configuration
Copy `config.example.yaml` to `~/.config/burrow/config.yaml` and customize:

```yaml
app:
  default_port: "3000"

database:
  migrations_dir: "sql/migrations"
  # path: ""  # Empty uses XDG default
```

#### Environment Variables
You can override settings with environment variables:
```bash
DEFAULT_PORT=3000 DB_FILE=/tmp/my.db burrow
```

## Features
- Interactive TUI for building HTTP requests
- Support for GET, POST, PUT, DELETE and other HTTP methods
- Request/response history tracking
- SQLite database for persistent storage
- Custom headers, parameters, and request bodies
- Keyboard-driven navigation

## Requirements
- Go 1.25.1 or later
- SQLite3
- Terminal with color support

## Development

### Building
```bash
git clone https://github.com/ManoloEsS/burrow.git
cd burrow
go mod tidy
go build
```

### Testing
```bash
go test ./...
```

## Contributing
Contributions are welcome! Please ensure all tests pass and follow the existing code style.

## License
MIT License - see LICENSE file for details.

## FAQ

**Q: Where are requests stored?**
A: Requests are stored locally in a SQLite database file.

**Q: Can I use burrow without a mouse?**
A: Yes, burrow is designed to be fully keyboard-driven.

## Support
Please report issues through GitHub Issues.

## Changelog
See [releases page](https://github.com/ManoloEsS/burrow/releases) for version history.
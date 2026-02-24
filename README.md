# Burrow

[![CI](https://github.com/ManoloEsS/burrow/actions/workflows/ci.yaml/badge.svg)](https://github.com/ManoloEsS/burrow/actions/workflows/ci.yaml)

**Terminal-based HTTP Client and Go Server Manager**

Burrow is a keyboard-driven TUI application for building and sending HTTP requests while running and monitoring Go servers — all from your terminal.

It is designed as a lightweight, terminal-native alternative to tools like Postman or Insomnia, built specifically for Go developers.
![Burrow](assets/burrow.gif)
## Overview

When developing APIs in Go, you often:

- Run a local server
- Switch to another tool to test endpoints
- Jump between multiple applications repeatedly

Burrow keeps everything in one place:

- Build and send HTTP requests
- Start and stop Go servers
- Monitor server health
- Save and reuse requests

All without leaving the terminal.

## Installation

Install via `go install`:

```bash
go install github.com/ManoloEsS/burrow/cmd/burrow@latest
```

## Quick Start

1. Install Burrow.
2. Copy `test_go_server.go` from the `test_server` directory into a working folder.
3. Navigate into that folder and run:

   ```bash
   cd your/test/folder
   burrow
   ```

4. Press **Ctrl-G** and enter:

   ```
   test_go_server.go
   ```

5. Press **Ctrl-R** to start the server (ensure port `8080` is free).
6. Press **Ctrl-S** to send a request (leave fields empty).
7. View the response in the **Response** panel.

You are now running and testing a Go server from the same terminal interface.

## Features

- Interactive terminal UI built with `tview`
- Support for `GET`, `POST`, `PUT`, `DELETE`, `HEAD`
- Save requests to embedded SQLite database
- Start and stop Go server files
- Automatic server health checking
- YAML configuration support
- XDG Base Directory compliant storage
- Fully keyboard-driven (mouse optional)

## HTTP Client Usage

### Default Behavior

If the URL field is empty, Burrow sends requests to:

- `http://localhost:8080`

You can modify the default port via configuration or environment variables.

### URL Handling

Supported formats:

- `somewebsite.com`
- `http://somewebsite.com`
- `https://somewebsite.com`

If no protocol is provided, `https://` is automatically added.

### Local Development Shortcuts

Send to local port:

- `:3030`

Resolves to:

- `http://localhost:3030`

Send to specific endpoint:

- `:3030/foobar`

Resolves to:

- `http://localhost:3030/foobar`

If using default port:

- `/foo`

Resolves to:

- `http://localhost:8080/foo`

## Server Management

Burrow runs Go server files directly from your working directory.

If launched in:

- `home/app`

Entering:

- `server.go`

Will execute:

- `./server.go`

To run a file elsewhere, use the absolute path.

## Health Checker

When a server starts, Burrow launches a background goroutine that sends a `GET` request to:

- `/health`

every 5 seconds.

For this to function properly, your server must expose a `/health` endpoint.

Currently, Burrow supports Go servers only.

## Configuration

Burrow follows the XDG Base Directory Specification and supports YAML configuration.

It works out-of-the-box with sensible defaults.

### Configuration Priority

1. `~/.config/burrow/config.yaml`
2. Environment variables
3. Default values

### Default File Locations

- Config: `~/.config/burrow/config.yaml`
- Database: `~/.local/share/burrow/burrow.db`
- Logs: `~/.local/state/burrow/burrow_log`
- Server Cache: `~/.cache/burrow/servers/`

### Example Configuration

```yaml
app:
  default_port: "8080"

database:
  path: ""
```

If the database path is empty, Burrow uses the default XDG path.

### Environment Variables

Override configuration using:

```bash
DEFAULT_PORT=3000 burrow
DB_FILE=/tmp/mydb.db burrow
```

Available variables:

- `DEFAULT_PORT`
- `DB_FILE`

## Keybindings

### Request Form

- **Ctrl-F** – Focus form
- **Ctrl-S** – Send request
- **Ctrl-A** – Save request
- **Ctrl-U** – Clear form
- **Ctrl-N / Ctrl-P** – Navigate fields

### Response View

- **Ctrl-T** – Focus response
- **J / K** – Scroll

### Saved Requests

- **Ctrl-L** – Focus list
- **Ctrl-O** – Load request
- **Ctrl-D** – Delete request
- **J / K** – Navigate list

### Server Controls

- **Ctrl-G** – Focus server path
- **Ctrl-R** – Start server
- **Ctrl-X** – Stop server

### Exit

- **Ctrl-C**

## Technical Highlights

- Written in Go
- Embedded SQLite database
- Background health-check goroutine
- Structured configuration with YAML
- XDG-compliant file management
- Modular internal architecture
- TUI built with `tview`

## Design Decisions

### Terminal-First Interface

Burrow is intentionally built as a keyboard-driven TUI rather than a GUI application.

Reasons:

- Many Go developers work primarily in the terminal.
- Reduces context switching between tools.
- Encourages a fast, focused API development workflow.

### Go-Only Server Execution

Burrow currently supports running Go server files directly.

Reasons:

- Optimized for Go developers.
- Tight integration with the Go toolchain.
- Simplified execution and monitoring logic.

Future iterations may introduce support for Python or Node.js servers.

### Background Health Checker

When a server starts, Burrow launches a background goroutine that periodically checks:

- `GET /health`

Reasons:

- Encourages explicit health endpoints.
- Demonstrates safe concurrent design.
- Provides immediate developer feedback if the server crashes.

The health checker runs independently to avoid blocking UI interactions.

### XDG Compliance

Burrow follows the XDG Base Directory Specification for configuration and storage.

Reasons:

- Keeps the system predictable and Linux-friendly.
- Avoids polluting the home directory.
- Respects modern filesystem conventions.

### Embedded SQLite Database

Saved requests are persisted using SQLite.

Reasons:

- Zero external dependencies.
- Lightweight and fast.
- Durable local storage.
- Structured querying over flat-file storage.

### Configuration Hierarchy

Configuration priority:

1. YAML config
2. Environment variables
3. Defaults

Reasons:

- Mirrors production configuration patterns.
- Enables environment-specific overrides.
- Maintains predictable fallback behavior.

### Explicit URL Normalization

Burrow automatically:

- Adds `https://` if protocol is missing.
- Interprets `:PORT` as `http://localhost:PORT`.
- Supports relative endpoint paths like `/foo`.

Reasons:

- Optimizes for local API development.
- Reduces typing friction.
- Improves workflow efficiency.

### Separation of Concerns

The project separates:

- UI logic
- HTTP request handling
- Server lifecycle management
- Storage
- Configuration

Reasons:

- Improves maintainability.
- Simplifies future expansion.
- Reduces tight coupling between components.

### Concurrency Considerations

Concurrency is used intentionally:

- Health checker runs in a separate goroutine.
- Server lifecycle management avoids blocking the UI.
- Context cancellation ensures controlled shutdowns.

The goal is safe, predictable concurrency.

## Requirements

- Go `1.25.1` or later
- CGO enabled (SQLite dependency)
- Terminal with color support

## Development

```bash
git clone https://github.com/ManoloEsS/burrow.git

cd burrow
go mod tidy
go build -o burrow cmd/burrow/main.go
```

## Roadmap

- Support for additional body types (form-data, multipart, etc.)
- Multi-language server execution (Python, Node.js)
- Enhanced request history filtering
- Improved response formatting

## Contributing

Contributions, issues, and suggestions are welcome.

## License

MIT License

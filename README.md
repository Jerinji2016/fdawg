# FDAWG - Flutter Development Assistant with Go

FDAWG is a CLI tool designed to help manage Flutter projects. It provides various commands to streamline Flutter development workflows and includes a web interface for project management.

## Installation

```bash
# Clone the repository
git clone https://github.com/Jerinji2016/fdawg.git

# Build the binary
cd fdawg
go build -o fdawg ./cmd/flutter-manager

# Optional: Move to a directory in your PATH
mv fdawg /usr/local/bin/
```

## Usage

```bash
fdawg [command] [options]
```

## Available Commands

| Command | Description | Options |
|---------|-------------|---------|
| `serve` | Start a web server for project management | `--port` - Port to run the server on (default: 8080) |

## Examples

Start the web server:

```bash
fdawg serve
```

Start the web server on a custom port:

```bash
fdawg serve --port 3000
```

## Project Structure

```text
fdawg/
├── cmd/                  # Application entry points
│   └── flutter-manager/  # Main CLI application
├── internal/             # Private application code
│   ├── commands/         # CLI command implementations
│   └── server/           # Web server implementation
├── pkg/                  # Reusable packages
│   └── utils/            # Utility functions
└── web/                  # Web assets
    ├── static/           # Static assets (CSS, JS)
    └── templates/        # HTML templates
```

## License

[MIT License](LICENSE)

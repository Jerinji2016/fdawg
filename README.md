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
| `serve [directory]` | Start a web server for project management | `--port` - Port to run the server on (default: 8080) |
| `init [directory]`  | Check if the specified directory is a Flutter project | Optional directory path (defaults to current directory) |

## Examples

Start the web server for the current directory:
```bash
fdawg serve
```

Start the web server for a specific Flutter project:
```bash
fdawg serve /path/to/my/flutter/project
```

Start the web server on a custom port:
```bash
fdawg serve --port 3000
# or
fdawg serve -p 3000
```

Start the web server for a specific project on a custom port:
```bash
fdawg serve --port 3000 /path/to/my/flutter/project
# or
fdawg serve -p 3000 /path/to/my/flutter/project

# NOTE: The following will not take the port details
fdawg serve /path/to/my/flutter/project --port 3000
# or
fdawg serve /path/to/my/flutter/project -p 3000
```

Check if the current directory is a Flutter project:
```bash
fdawg init
```

Check if a specific directory is a Flutter project:
```bash
fdawg init /path/to/my/flutter/project
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

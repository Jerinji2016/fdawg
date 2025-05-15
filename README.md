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
│       ├── web/          # Web assets
│       │   ├── static/   # Static assets
│       │   │   ├── css/  # Compiled CSS
│       │   │   ├── js/   # JavaScript files
│       │   │   └── scss/ # SASS source files
│       │   └── templates/# HTML templates
├── pkg/                  # Reusable packages
│   └── utils/            # Utility functions
```

## SASS Structure

The project uses SASS for styling with a modular structure:

```text
scss/
├── _variables.scss      # Variables for colors, fonts, sizes, etc.
├── _base.scss           # Base styles and typography
├── _layout.scss         # Layout components (header, sidebar, content)
├── _components.scss     # UI components (buttons, cards, etc.)
├── _responsive.scss     # Media queries for responsive design
└── main.scss            # Main file that imports all partials
```

## Development

### SASS Compilation

The project includes npm scripts for SASS compilation:

```bash
# Install dependencies
npm install

# Compile SASS to CSS
npm run sass

# Watch SASS files for changes
npm run sass:watch
```

### Makefile Commands

```bash
# Build the project
make build

# Compile SASS and build
make all

# Compile SASS
make sass

# Watch SASS files for changes
make sass-watch

# Run the application (compiles SASS first)
make run

# Run in development mode (watches SASS files)
make dev
```

## License

[MIT License](LICENSE)

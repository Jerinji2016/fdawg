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

| Command | Description |
|---------|-------------|
| `serve [directory]` | Start a web server for project management |
| `init [directory]`  | Check if the specified directory is a Flutter project |
| `env list` | List all environment files in the `.environment` directory |
| `env show <env-name>` | Show all variables in a specific environment file |
| `env create <env-name>` | Create a new environment file |
| `env add <key> <value>` | Add or update a variable in an environment file |
| `env delete <env-name>` | Delete an environment file |
| `env remove <key>` | Remove a variable from an environment file |
| `env generate-dart` | Generate a Dart environment file with all environment variables |
| `asset add <asset-path>` | Add an asset to the project |
| `asset remove <asset-name>` | Remove an asset from the project |
| `asset list` | List all assets in the project |
| `asset generate-dart` | Generate a Dart asset file with all assets |
| `asset migrate` | Migrate assets to organized folders by type and clean up empty directories |

### Environment Command Options

- `env show <env-name>`:
  - `<env-name>`: Name of the environment file

- `env create <env-name>`:
  - `<env-name>`: Name of the environment file
  - `--copy` or `-c`: Copy variables from an existing environment file

- `env add <key> <value>`:
  - `<key>`: Variable key (must start with a letter or underscore)
  - `<value>`: Variable value
  - `--env` or `-e`: Environment file to add the variable to (default: development)

- `env delete <env-name>`:
  - `<env-name>`: Name of the environment file to delete

- `env remove <key>`:
  - `<key>`: Variable key to remove
  - `--env` or `-e`: Environment file to remove the variable from (default: development)

### Environment Command Examples

List all environment files:

```bash
fdawg env list
```

Show variables in the development environment:

```bash
fdawg env show development
```

Create a new environment file:

```bash
fdawg env create production
```

Create a new environment file by copying from an existing one:

```bash
fdawg env create staging --copy development
```

Add a variable to the default (development) environment:

```bash
fdawg env add API_URL https://api.example.com
```

Add a variable to a specific environment:

```bash
fdawg env add DEBUG_MODE false --env production
```

Remove a variable from an environment:

```bash
fdawg env remove API_URL --env staging
```

Delete an environment file:

```bash
fdawg env delete staging
```

Generate the Dart environment file:

```bash
fdawg env generate-dart
```

### Asset Command Options

- `asset add <asset-path>`:
  - `--type, -t`: Type of asset (images, animations, audio, videos, json, svgs, misc)
  - If no type is specified, it will be determined from the file extension

- `asset remove <asset-name>`:
  - `--type, -t`: Type of asset (images, animations, audio, videos, json, svgs, misc)
  - If no type is specified, it will search for the asset in all directories

- `asset migrate`:
  - Migrates assets to organized folders by type (images, animations, audio, etc.)
  - Creates a backup of all assets in assets.backup directory
  - Automatically categorizes assets based on file type and content
  - Recursively processes all files in nested directories
  - Removes empty directories after migration

### Asset Command Examples

Add an image asset:

```bash
fdawg asset add path/to/image.png
```

Add an asset with a specific type:

```bash
fdawg asset add path/to/file.json --type animations
```

Remove an asset:

```bash
fdawg asset remove image.png
```

Remove an asset with a specific type:

```bash
fdawg asset remove font.ttf --type fonts
```

List all assets:

```bash
fdawg asset list
```

Generate the Dart asset file:

```bash
fdawg asset generate-dart
```

Migrate assets to organized folders:

```bash
fdawg asset migrate
```

## Server and Init Command Examples

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

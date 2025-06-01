---
layout: default
title: Command Reference
nav_order: 3
has_children: true
description: "Complete reference for all FDAWG commands"
permalink: /commands/
---

# Command Reference

FDAWG provides a comprehensive set of CLI commands for managing Flutter projects. All commands follow the pattern:

```bash
fdawg [command] [subcommand] [options] [arguments]
```

## Available Commands

### Core Commands

| Command | Description | Documentation |
|---------|-------------|---------------|
| [`serve`](server.html) | Start web server for project management | [Server Commands](server.html) |
| [`init`](server.html#init) | Validate Flutter project structure | [Server Commands](server.html) |

### Project Management Commands

| Command | Description | Documentation |
|---------|-------------|---------------|
| [`env`](environment.html) | Environment variable management | [Environment Commands](environment.html) |
| [`asset`](assets.html) | Project asset management | [Asset Commands](assets.html) |
| [`lang`](localization.html) | Localization and translation management | [Localization Commands](localization.html) |
| [`namer`](namer.html) | Cross-platform app naming | [App Namer Commands](namer.html) |

## Quick Examples

### Environment Management
```bash
# List all environment files
fdawg env list

# Create a new environment
fdawg env create production

# Add a variable
fdawg env add API_URL https://api.example.com --env production
```

### Asset Management
```bash
# Add an asset
fdawg asset add path/to/image.png

# Migrate assets to organized folders
fdawg asset migrate

# Generate Dart asset file
fdawg asset generate-dart
```

### Localization
```bash
# Initialize localization
fdawg lang init

# Add a new language
fdawg lang add es

# Add a translation key
fdawg lang insert app.welcome
```

### App Naming
```bash
# Get current app names
fdawg namer list

# Set app name for all platforms
fdawg namer set --value "My Awesome App"

# Set platform-specific names
fdawg namer set --android "Android App" --ios "iOS App"
```

### Web Interface
```bash
# Start web server on default port (8080)
fdawg serve

# Start on custom port
fdawg serve --port 3000

# Serve specific project
fdawg serve /path/to/flutter/project
```

## Global Options

Most commands support these global options:

- `--help, -h`: Show help for the command
- `--version, -v`: Show version information (root command only)

## Command Structure

### Environment Commands (`env`)
- `list` - List all environment files
- `show <env-name>` - Show variables in an environment
- `create <env-name>` - Create a new environment file
- `add <key> <value>` - Add/update a variable
- `remove <key>` - Remove a variable
- `delete <env-name>` - Delete an environment file
- `generate-dart` - Generate Dart environment file

### Asset Commands (`asset`)
- `list` - List all project assets
- `add <asset-path>` - Add an asset to the project
- `remove <asset-name>` - Remove an asset
- `migrate` - Organize assets into folders by type
- `generate-dart` - Generate Dart asset file

### Localization Commands (`lang`)
- `list` - List supported languages
- `init` - Initialize localization
- `add <language-code>` - Add language support
- `remove <language-code>` - Remove language support
- `insert <key>` - Add translation key
- `delete <key>` - Delete translation key

### App Namer Commands (`namer`)
- `list` / `get` - Get current app names
- `set` - Set app names (universal or platform-specific)

## Error Handling

FDAWG includes comprehensive error handling:

- **Validation**: Commands validate Flutter project structure
- **Backups**: Destructive operations create automatic backups
- **Rollback**: Failed operations automatically restore from backups
- **Clear messages**: Detailed error messages with suggested solutions

## Tips and Best Practices

1. **Always run commands from your Flutter project root** - FDAWG validates the project structure
2. **Use the web interface for complex operations** - Visual feedback and validation
3. **Create backups before major changes** - Some commands create automatic backups
4. **Test in development first** - Try commands in a test project before production use

---

Choose a command category below to explore detailed documentation and examples:

- [🔧 Environment Commands](environment.html) - Manage environment variables
- [📦 Asset Commands](assets.html) - Handle project assets  
- [🌍 Localization Commands](localization.html) - Translation management
- [🏷️ App Namer Commands](namer.html) - Cross-platform app naming
- [🌐 Server Commands](server.html) - Web interface and validation

---
layout: home
title: FDAWG Documentation
---

# FDAWG - Flutter Development Assistant with Go

Welcome to the comprehensive documentation for FDAWG, a powerful CLI tool and web interface designed to streamline Flutter development workflows.

## What is FDAWG?

FDAWG is a comprehensive Flutter project management tool that provides:

- **Modern Web Interface** - Intuitive dashboard for visual project management
- **Environment Management** - Handle multiple environment configurations with ease
- **Asset Management** - Organize and manage project assets efficiently
- **Localization Support** - Multi-language support with translation management
- **Cross-Platform Naming** - Manage app names across all platforms
- **CLI & Web Interface** - Use via command line or modern web dashboard

## Quick Navigation

### 🚀 Getting Started
- [Installation Guide](installation.html) - Get FDAWG up and running
- [Quick Start](installation.html#quick-start) - Basic usage examples

### 📖 Command Reference
- [Environment Commands](commands/environment.html) - Manage environment variables
- [Asset Commands](commands/assets.html) - Handle project assets
- [Localization Commands](commands/localization.html) - Translation management
- [App Namer Commands](commands/namer.html) - Cross-platform app naming
- [Server Commands](commands/server.html) - Web interface and project validation

### 🌐 Web Interface
- [Web Interface Guide](web-interface.html) - Complete web dashboard documentation
- [Features Overview](web-interface.html#features) - All available web features

### 👨‍💻 Development
- [Development Setup](development.html) - Contributing to FDAWG
- [Project Structure](development.html#project-structure) - Understanding the codebase
- [Building from Source](development.html#building) - Compilation instructions

## Key Features

### CLI Commands
FDAWG provides a comprehensive set of CLI commands for all aspects of Flutter project management:

```bash
# Environment management
fdawg env list
fdawg env create production
fdawg env add API_URL https://api.example.com

# Asset management
fdawg asset add path/to/image.png
fdawg asset migrate
fdawg asset generate-dart

# Localization
fdawg lang init
fdawg lang add es
fdawg lang insert app.welcome

# App naming
fdawg namer set --value "My App"
fdawg namer get --platforms android,ios
```

### Web Interface
Launch a modern, responsive web interface for visual project management:

```bash
fdawg serve
```

The web interface includes:
- Project overview and statistics
- Visual environment variable management
- Drag-and-drop asset management
- Translation management with Google Translate integration
- Cross-platform app naming with live preview

## Support

- **Issues**: [GitHub Issues](https://github.com/Jerinji2016/fdawg/issues)
- **Discussions**: [GitHub Discussions](https://github.com/Jerinji2016/fdawg/discussions)
- **License**: [MIT License](https://github.com/Jerinji2016/fdawg/blob/main/LICENSE)

---

Ready to get started? Check out our [Installation Guide](installation.html) to begin using FDAWG in your Flutter projects.

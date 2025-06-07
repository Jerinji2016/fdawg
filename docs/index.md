---
layout: default
title: FDAWG - Flutter Development Assistant with Go
description: "FDAWG - Flutter Development Assistant with Go"
---

# FDAWG - Flutter Development Assistant with Go

<div align="center">
  <img src="{{ '/assets/images/fdawg_logo.png' | relative_url }}" alt="FDAWG Logo" width="200" height="200">
</div>

Welcome to the comprehensive documentation for FDAWG, a powerful CLI tool and web interface designed to streamline Flutter development workflows.

## What is FDAWG?

FDAWG is a comprehensive Flutter project management tool that provides:

- 🌐 **Modern Web Interface** - Intuitive dashboard for visual project management
- 🔧 **Environment Management** - Handle multiple environment configurations with ease
- 📦 **Asset Management** - Organize and manage project assets efficiently
- 🌍 **Localization Support** - Multi-language support with translation management
- 🏷️ **Cross-Platform Naming** - Manage app names across all platforms
- 🚀 **CLI & Web Interface** - Use via command line or modern web dashboard

{: .callout .callout-info }
> **Quick Start Tip**
>
> New to FDAWG? Start with our [Installation Guide](installation/) to get up and running in minutes!

## Quick Navigation

### 🚀 Getting Started
- [Installation Guide](installation/) - Get FDAWG up and running
- [Quick Start](installation/#quick-start) - Basic usage examples

### 📖 Command Reference
- [Environment Commands](commands/environment/) - Manage environment variables
- [Asset Commands](commands/assets/) - Handle project assets
- [Localization Commands](commands/localization/) - Translation management
- [App Namer Commands](commands/namer/) - Cross-platform app naming
- [Bundle ID Commands](commands/bundler/) - Bundle identifier management
- [Build Commands](commands/build/) - Comprehensive build system with multi-platform support
- [Server Commands](commands/server/) - Web interface and project validation

### 🌐 Web Interface
- [Web Interface Guide](web-interface/) - Complete web dashboard documentation
- [Features Overview](web-interface/#features) - All available web features

### 👨‍💻 Development
- [Development Setup](development/) - Contributing to FDAWG
- [Project Structure](development/#project-structure) - Understanding the codebase
- [Building from Source](development/#building) - Compilation instructions

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

# Bundle ID management
fdawg bundler set --universal com.company.app
fdawg bundler get

# Build management
fdawg build setup
fdawg build run --platforms android,ios
fdawg build status
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
- Bundle ID management with validation
- Comprehensive build management with real-time progress

## Support

- **Issues**: [GitHub Issues](https://github.com/Jerinji2016/fdawg/issues)
- **Discussions**: [GitHub Discussions](https://github.com/Jerinji2016/fdawg/discussions)
- **License**: [MIT License](https://github.com/Jerinji2016/fdawg/blob/main/LICENSE)

---

Ready to get started? Check out our [Installation Guide](installation/) to begin using FDAWG in your Flutter projects.

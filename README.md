# FDAWG - Flutter Development Assistant with Go

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.23.2+-blue.svg)](https://golang.org)
[![Flutter](https://img.shields.io/badge/Flutter-Compatible-blue.svg)](https://flutter.dev)
[![Deploy Documentation](https://github.com/Jerinji2016/fdawg/actions/workflows/docs.yml/badge.svg)](https://github.com/Jerinji2016/fdawg/actions/workflows/docs.yml)

FDAWG is a powerful CLI tool and web interface designed to streamline Flutter development workflows. It provides comprehensive project management capabilities including asset management, environment configuration, localization, cross-platform app naming, and much more.

## ✨ Features

- 🌐 **Modern Web Interface** - Intuitive project management dashboard
- 🔧 **Environment Management** - Handle multiple environment configurations
- 📦 **Asset Management** - Organize and manage project assets efficiently
- 🌍 **Localization Support** - Multi-language support with easy translation management
- 🏷️ **Cross-Platform Naming** - Manage app names across all platforms
- 🆔 **Bundle ID Management** - Configure bundle identifiers for all platforms
- 🚀 **CLI & Web Interface** - Use via command line or modern web dashboard

## 🚀 Quick Start

### Installation

```bash
# Clone the repository
git clone https://github.com/Jerinji2016/fdawg.git
cd fdawg

# Build the binary
make build

# Optional: Add to PATH
mv build/fdawg /usr/local/bin/
```

### Basic Usage

```bash
# Start web interface
fdawg serve

# Or use CLI commands
fdawg env list
fdawg asset list
fdawg lang list
```

## 📋 Available Commands

| Command | Description |
|---------|-------------|
| `serve [directory]` | Start web server for project management |
| `init [directory]` | Validate Flutter project structure |
| `env` | Environment variable management |
| `asset` | Project asset management |
| `lang` | Localization and translation management |
| `namer` | Cross-platform app naming |
| `bundler` | Bundle ID management for all platforms |

For detailed command documentation and examples, visit our [comprehensive documentation](https://jerinji2016.github.io/fdawg/).

## 🌐 Web Interface

Launch the modern web interface for visual project management:

```bash
fdawg serve
```

The web interface provides:
- 📊 Project overview and statistics
- 🔧 Visual environment management
- 📦 Drag-and-drop asset management
- 🌍 Translation management with Google Translate integration
- 🏷️ Cross-platform app naming with live preview
- 🆔 Bundle ID management with validation and platform-specific configuration

## 📚 Documentation

- [📖 Complete Documentation](https://jerinji2016.github.io/fdawg/)
- [⚡ Quick Start Guide](https://jerinji2016.github.io/fdawg/installation)
- [🔧 Command Reference](https://jerinji2016.github.io/fdawg/commands/)
- [🌐 Web Interface Guide](https://jerinji2016.github.io/fdawg/web-interface)
- [👨‍💻 Development Setup](https://jerinji2016.github.io/fdawg/development)

## 🤝 Contributing

We welcome contributions! Please see our [Development Guide](https://jerinji2016.github.io/fdawg/development) for details on:
- Setting up the development environment
- Building and testing
- Code style guidelines
- Submitting pull requests

## 📄 License

This project is licensed under the [MIT License](LICENSE).

---

<div align="center">
  <strong>Made with ❤️ for the Flutter community</strong>
</div>

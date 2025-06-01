---
layout: page
title: Installation Guide
permalink: /installation/
---

# Installation Guide

This guide will help you install and set up FDAWG on your system.

## Prerequisites

- **Go 1.23.2+** - Required for building from source
- **Node.js** - Required for SASS compilation (development only)
- **Flutter SDK** - For the projects you'll be managing
- **Git** - For cloning the repository

## Installation Methods

### Method 1: Build from Source (Recommended)

1. **Clone the repository**:
   ```bash
   git clone https://github.com/Jerinji2016/fdawg.git
   cd fdawg
   ```

2. **Build the binary**:
   ```bash
   make build
   ```

3. **Install globally** (optional):
   ```bash
   # On macOS/Linux
   sudo mv build/fdawg /usr/local/bin/

   # Or add to your PATH
   export PATH=$PATH:$(pwd)/build
   ```

### Method 2: Direct Go Build

If you prefer to use Go directly:

```bash
git clone https://github.com/Jerinji2016/fdawg.git
cd fdawg
go build -o fdawg ./cmd/flutter-manager
```

## Verification

Verify your installation by running:

```bash
fdawg --version
```

You should see the FDAWG version information.

## Quick Start

### 1. Navigate to Your Flutter Project

```bash
cd /path/to/your/flutter/project
```

### 2. Validate Your Project

Check if your directory is a valid Flutter project:

```bash
fdawg init
```

### 3. Start the Web Interface

Launch the modern web interface:

```bash
fdawg serve
```

This will start a web server (default port 8080) and open your project management dashboard.

### 4. Try CLI Commands

Explore the available commands:

```bash
# List environment files
fdawg env list

# List project assets
fdawg asset list

# List supported languages
fdawg lang list

# Get current app names
fdawg namer list
```

## Configuration

### Environment Setup

FDAWG works with Flutter projects and doesn't require additional configuration. However, you can customize:

- **Web server port**: Use `--port` or `-p` flag with the `serve` command
- **Project directory**: Specify a different directory as an argument

### Development Setup

If you plan to contribute to FDAWG development:

1. **Install Node.js dependencies**:
   ```bash
   npm install
   ```

2. **Compile SASS** (for web interface styling):
   ```bash
   npm run sass
   ```

3. **Watch for changes** (development mode):
   ```bash
   make dev
   ```

## Troubleshooting

### Common Issues

**Issue**: `fdawg: command not found`
- **Solution**: Ensure the binary is in your PATH or use the full path to the executable

**Issue**: Web interface doesn't load
- **Solution**: Check if the port is already in use, try a different port with `--port` flag

**Issue**: Commands fail with "not a Flutter project"
- **Solution**: Ensure you're in a valid Flutter project directory with a `pubspec.yaml` file

**Issue**: Build fails
- **Solution**: Ensure you have Go 1.23.2+ installed and run `make clean` before building

### Getting Help

If you encounter issues:

1. Check the [GitHub Issues](https://github.com/Jerinji2016/fdawg/issues) for known problems
2. Create a new issue with:
   - Your operating system
   - Go version (`go version`)
   - Flutter version (`flutter --version`)
   - Complete error message
   - Steps to reproduce

## Next Steps

Now that FDAWG is installed:

- Explore the [Command Reference](commands/) for detailed command documentation
- Learn about the [Web Interface](web-interface.html) features
- Check out the [Development Guide](development.html) if you want to contribute

---

Ready to start managing your Flutter projects more efficiently? Head over to the [Command Reference](commands/) to learn about all available features!

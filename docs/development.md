---
layout: default
title: Development Guide
nav_order: 5
description: "Contributing to FDAWG development"
permalink: /development/
---

# Development Guide

This guide covers everything you need to know about contributing to FDAWG, including setting up the development environment, understanding the codebase, and following our development practices.

## Prerequisites

Before you start developing FDAWG, ensure you have:

- **Go 1.23.2+** - [Download Go](https://golang.org/dl/)
- **Node.js 16+** - [Download Node.js](https://nodejs.org/)
- **Git** - [Download Git](https://git-scm.com/)
- **Flutter SDK** - [Install Flutter](https://flutter.dev/docs/get-started/install) (for testing)
- **Make** - Usually pre-installed on macOS/Linux, [install on Windows](https://gnuwin32.sourceforge.net/packages/make.htm)

## Getting Started

### 1. Fork and Clone

```bash
# Fork the repository on GitHub, then clone your fork
git clone https://github.com/YOUR_USERNAME/fdawg.git
cd fdawg

# Add upstream remote
git remote add upstream https://github.com/Jerinji2016/fdawg.git
```

### 2. Install Dependencies

```bash
# Install Go dependencies
go mod download

# Install Node.js dependencies for SASS compilation
npm install
```

### 3. Build and Test

```bash
# Clean and build
make clean
make build

# Test the binary
./build/fdawg --version
```

### 4. Development Workflow

```bash
# Start development mode (watches SASS files)
make dev

# In another terminal, test your changes
./build/fdawg serve /path/to/test/flutter/project
```

## Project Structure

```
fdawg/
├── cmd/                    # Application entry points
│   └── flutter-manager/    # Main CLI application
│       └── main.go         # Application entry point
├── internal/               # Private application code
│   ├── commands/           # CLI command implementations
│   │   ├── asset.go        # Asset management commands
│   │   ├── environment.go  # Environment commands
│   │   ├── init.go         # Project validation
│   │   ├── localization.go # Localization commands
│   │   ├── namer.go        # App naming commands
│   │   └── serve.go        # Web server command
│   └── server/             # Web server implementation
│       ├── handlers/       # HTTP request handlers
│       ├── web/            # Web assets and templates
│       │   ├── static/     # Static assets (CSS, JS, images)
│       │   │   ├── css/    # Compiled CSS files
│       │   │   ├── js/     # JavaScript files
│       │   │   └── scss/   # SASS source files
│       │   └── templates/  # HTML templates
│       └── server.go       # Server setup and routing
├── pkg/                    # Reusable packages
│   ├── asset/              # Asset management functionality
│   ├── config/             # Configuration management
│   ├── environment/        # Environment variables management
│   ├── flutter/            # Flutter project utilities
│   ├── localization/       # Localization management
│   ├── namer/              # App name management
│   ├── translate/          # Translation services
│   └── utils/              # Utility functions
├── docs/                   # Documentation (GitHub Pages)
├── dwag_tests/             # Test Flutter project
├── Makefile                # Build automation
├── package.json            # Node.js dependencies
└── go.mod                  # Go module definition
```

## Architecture Overview

### Command Structure

FDAWG uses the [urfave/cli](https://github.com/urfave/cli) library for command-line interface:

```go
// cmd/flutter-manager/main.go
app := &cli.App{
    Name:  "fdawg",
    Usage: "Flutter Development Assistant with Go",
    Commands: []*cli.Command{
        commands.ServeCommand,
        commands.InitCommand,
        commands.EnvironmentCommand,
        commands.AssetCommand,
        commands.LocalizationCommand,
        commands.NamerCommand,
    },
}
```

### Package Organization

**`pkg/` packages** are designed to be reusable and independent:
- Each package handles a specific domain (assets, environment, etc.)
- Packages don't depend on CLI-specific code
- Can be imported by other Go projects

**`internal/` packages** contain application-specific code:
- CLI command implementations
- Web server and handlers
- Application-specific logic

### Web Server Architecture

The web server uses Go's standard `net/http` package with:
- **Static file serving** for CSS, JS, and images
- **Template rendering** for HTML pages
- **JSON API endpoints** for dynamic functionality
- **WebSocket support** for real-time updates (planned)

## Development Practices

### Code Style

We follow standard Go conventions:

```bash
# Format code
go fmt ./...

# Lint code (install golangci-lint first)
golangci-lint run

# Vet code
go vet ./...
```

### Testing

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests for specific package
go test ./pkg/asset/...
```

### SASS Development

For web interface styling:

```bash
# Compile SASS once
npm run sass

# Watch for changes (development)
npm run sass:watch

# Or use Make
make sass-watch
```

## Building

### Development Build

```bash
# Clean previous builds
make clean

# Build with SASS compilation
make all

# Or build without SASS
make build
```

### Release Build

```bash
# Build optimized binary
make clean
make all

# The binary will be in build/fdawg
```

### Cross-Platform Building

```bash
# Build for different platforms
GOOS=linux GOARCH=amd64 go build -o build/fdawg-linux ./cmd/flutter-manager
GOOS=windows GOARCH=amd64 go build -o build/fdawg-windows.exe ./cmd/flutter-manager
GOOS=darwin GOARCH=amd64 go build -o build/fdawg-macos ./cmd/flutter-manager
```

## Testing

### Test Flutter Project

Use the included test project for development:

```bash
# Navigate to test project
cd dwag_tests

# Test FDAWG commands
../build/fdawg init
../build/fdawg serve
../build/fdawg env list
```

### Unit Tests

Write tests for new functionality:

```go
// pkg/asset/asset_test.go
func TestAddAsset(t *testing.T) {
    // Test implementation
}
```

### Integration Tests

Test complete workflows:

```bash
# Test complete environment workflow
fdawg env create test
fdawg env add TEST_VAR test_value --env test
fdawg env show test
fdawg env delete test
```

## Contributing

### Workflow

1. **Create a feature branch:**
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. **Make your changes:**
   - Follow coding standards
   - Add tests for new functionality
   - Update documentation

3. **Test your changes:**
   ```bash
   make clean
   make all
   go test ./...
   ```

4. **Commit your changes:**
   ```bash
   git add .
   git commit -m "feat: add new feature description"
   ```

5. **Push and create PR:**
   ```bash
   git push origin feature/your-feature-name
   # Create pull request on GitHub
   ```

### Commit Message Format

We use conventional commits:

- `feat:` - New features
- `fix:` - Bug fixes
- `docs:` - Documentation changes
- `style:` - Code style changes
- `refactor:` - Code refactoring
- `test:` - Test additions or changes
- `chore:` - Build process or auxiliary tool changes

### Pull Request Guidelines

- **Clear description** of changes
- **Link to related issues** if applicable
- **Include tests** for new functionality
- **Update documentation** as needed
- **Ensure CI passes** before requesting review

## Adding New Features

### Adding a New Command

1. **Create command file:**
   ```go
   // internal/commands/newcommand.go
   package commands
   
   import "github.com/urfave/cli/v2"
   
   var NewCommand = &cli.Command{
       Name:  "newcommand",
       Usage: "Description of new command",
       Action: func(c *cli.Context) error {
           // Implementation
           return nil
       },
   }
   ```

2. **Add to main.go:**
   ```go
   // cmd/flutter-manager/main.go
   Commands: []*cli.Command{
       // ... existing commands
       commands.NewCommand,
   },
   ```

3. **Create package if needed:**
   ```go
   // pkg/newfeature/newfeature.go
   package newfeature
   
   // Implementation
   ```

### Adding Web Interface Features

1. **Create HTML template:**
   ```html
   <!-- internal/server/web/templates/newfeature.html -->
   {{define "content"}}
   <div class="project-info">
       <!-- Your HTML content -->
   </div>
   {{end}}
   ```

2. **Add route handler:**
   ```go
   // internal/server/handlers/newfeature.go
   func NewFeatureHandler(w http.ResponseWriter, r *http.Request) {
       // Handler implementation
   }
   ```

3. **Add JavaScript if needed:**
   ```javascript
   // internal/server/web/static/js/newfeature.js
   // JavaScript functionality
   ```

4. **Update navigation:**
   ```html
   <!-- internal/server/web/templates/layout.html -->
   <li>
       <a href="/newfeature">
           <i class="fas fa-icon"></i>
           <span>New Feature</span>
       </a>
   </li>
   ```

## Debugging

### Debug Mode

```bash
# Enable verbose logging
export FDAWG_DEBUG=true
./build/fdawg serve

# Or use Go's built-in debugging
go run -race ./cmd/flutter-manager serve
```

### Common Issues

**Build failures:**
- Check Go version compatibility
- Ensure all dependencies are installed
- Run `go mod tidy`

**SASS compilation issues:**
- Check Node.js version
- Run `npm install`
- Verify SASS files syntax

**Test failures:**
- Ensure test Flutter project is valid
- Check file permissions
- Verify Go module dependencies

## Release Process

### Version Management

1. **Update version** in relevant files
2. **Create git tag:**
   ```bash
   git tag -a v1.0.0 -m "Release version 1.0.0"
   git push origin v1.0.0
   ```

3. **Build release binaries:**
   ```bash
   make clean
   make all
   # Build for multiple platforms
   ```

4. **Create GitHub release** with binaries

### Documentation Updates

- Update README.md if needed
- Update documentation in `docs/`
- Ensure GitHub Pages builds correctly

## Getting Help

### Resources

- **GitHub Issues**: [Report bugs or request features](https://github.com/Jerinji2016/fdawg/issues)
- **GitHub Discussions**: [Ask questions or discuss ideas](https://github.com/Jerinji2016/fdawg/discussions)
- **Documentation**: [Complete documentation](https://jerinji2016.github.io/fdawg/)

### Development Questions

When asking for help:
1. Describe what you're trying to achieve
2. Include relevant code snippets
3. Provide error messages if any
4. Mention your development environment

---

Thank you for contributing to FDAWG! Your contributions help make Flutter development more efficient for everyone.

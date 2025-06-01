# Contributing to FDAWG

Thank you for your interest in contributing to FDAWG! We welcome contributions from the community and are excited to see what you'll bring to the project.

## 📋 Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Setup](#development-setup)
- [How to Contribute](#how-to-contribute)
- [Pull Request Process](#pull-request-process)
- [Coding Standards](#coding-standards)
- [Testing Guidelines](#testing-guidelines)
- [Documentation](#documentation)
- [Community](#community)

## Code of Conduct

This project and everyone participating in it is governed by our Code of Conduct. By participating, you are expected to uphold this code. Please report unacceptable behavior to the project maintainers.

### Our Standards

- **Be respectful** and inclusive of different viewpoints and experiences
- **Be collaborative** and help others learn and grow
- **Be constructive** in feedback and discussions
- **Be patient** with newcomers and those learning
- **Focus on what's best** for the community and project

## Getting Started

### Prerequisites

Before contributing, ensure you have:

- **Go 1.23.2+** - [Download Go](https://golang.org/dl/)
- **Node.js 16+** - [Download Node.js](https://nodejs.org/)
- **Git** - [Download Git](https://git-scm.com/)
- **Flutter SDK** - [Install Flutter](https://flutter.dev/docs/get-started/install)
- **Make** - Usually pre-installed on macOS/Linux

### First Steps

1. **Read the documentation**: [jerinji2016.github.io/fdawg](https://jerinji2016.github.io/fdawg/)
2. **Explore the codebase**: Understand the project structure
3. **Run FDAWG locally**: Follow the development setup guide
4. **Look for good first issues**: Check issues labeled `good first issue`

## Development Setup

### 1. Fork and Clone

```bash
# Fork the repository on GitHub, then:
git clone https://github.com/YOUR_USERNAME/fdawg.git
cd fdawg
git remote add upstream https://github.com/Jerinji2016/fdawg.git
```

### 2. Install Dependencies

```bash
# Go dependencies
go mod download

# Node.js dependencies (for SASS)
npm install
```

### 3. Build and Test

```bash
# Build the project
make clean
make build

# Test the binary
./build/fdawg --version

# Run tests
go test ./...
```

### 4. Development Workflow

```bash
# Start development mode (watches SASS files)
make dev

# In another terminal, test your changes
./build/fdawg serve /path/to/test/flutter/project
```

## How to Contribute

### Types of Contributions

We welcome various types of contributions:

- **🐛 Bug fixes** - Fix issues and improve stability
- **✨ New features** - Add functionality that benefits users
- **📚 Documentation** - Improve guides, examples, and API docs
- **🧪 Tests** - Add test coverage and improve test quality
- **🎨 UI/UX** - Enhance the web interface design
- **⚡ Performance** - Optimize code and improve efficiency
- **🔧 Tooling** - Improve build process and development tools

### Finding Work

1. **Check existing issues**: Look for issues labeled:
   - `good first issue` - Great for newcomers
   - `help wanted` - Community help needed
   - `bug` - Bug fixes needed
   - `enhancement` - New features

2. **Create new issues**: If you find bugs or have ideas:
   - Use our issue templates
   - Provide detailed descriptions
   - Include reproduction steps for bugs

3. **Join discussions**: Participate in GitHub Discussions

## Pull Request Process

### Before You Start

1. **Check for existing work**: Search issues and PRs
2. **Discuss major changes**: Open an issue for significant features
3. **Create a branch**: Use descriptive branch names

```bash
git checkout -b feature/your-feature-name
# or
git checkout -b fix/bug-description
```

### Making Changes

1. **Follow coding standards**: See [Coding Standards](#coding-standards)
2. **Write tests**: Add tests for new functionality
3. **Update documentation**: Keep docs current
4. **Test thoroughly**: Verify your changes work

### Submitting Your PR

1. **Commit your changes**:
   ```bash
   git add .
   git commit -m "feat: add new feature description"
   ```

2. **Push to your fork**:
   ```bash
   git push origin feature/your-feature-name
   ```

3. **Create a Pull Request**:
   - Use our PR template
   - Provide clear description
   - Link related issues
   - Add screenshots if applicable

### PR Review Process

1. **Automated checks**: Ensure CI passes
2. **Code review**: Maintainers will review your code
3. **Address feedback**: Make requested changes
4. **Final approval**: PR will be merged when approved

## Coding Standards

### Go Code Style

- Follow standard Go conventions
- Use `gofmt` for formatting
- Run `go vet` to catch issues
- Use meaningful variable and function names
- Add comments for exported functions

```bash
# Format code
go fmt ./...

# Vet code
go vet ./...

# Lint (if golangci-lint is installed)
golangci-lint run
```

### Commit Messages

Use conventional commit format:

```
type(scope): description

[optional body]

[optional footer]
```

**Types:**
- `feat:` - New features
- `fix:` - Bug fixes
- `docs:` - Documentation changes
- `style:` - Code style changes
- `refactor:` - Code refactoring
- `test:` - Test additions/changes
- `chore:` - Build/tool changes

**Examples:**
```
feat(asset): add drag-and-drop upload support
fix(env): resolve variable validation issue
docs(readme): update installation instructions
```

### Code Organization

- Keep functions small and focused
- Use packages to organize related functionality
- Separate CLI logic from core functionality
- Follow the existing project structure

## Testing Guidelines

### Writing Tests

- Write tests for new functionality
- Include edge cases and error conditions
- Use table-driven tests when appropriate
- Mock external dependencies

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific package tests
go test ./pkg/asset/...

# Run tests with verbose output
go test -v ./...
```

### Test Structure

```go
func TestFunctionName(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected string
        wantErr  bool
    }{
        {
            name:     "valid input",
            input:    "test",
            expected: "expected",
            wantErr:  false,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result, err := FunctionName(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("FunctionName() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if result != tt.expected {
                t.Errorf("FunctionName() = %v, want %v", result, tt.expected)
            }
        })
    }
}
```

## Documentation

### Types of Documentation

- **Code comments**: Document exported functions and complex logic
- **README updates**: Keep the main README current
- **API documentation**: Document public APIs
- **User guides**: Help users understand features
- **Developer docs**: Help contributors understand the codebase

### Documentation Standards

- Use clear, concise language
- Provide examples where helpful
- Keep documentation up-to-date with code changes
- Use proper markdown formatting

## Community

### Communication Channels

- **GitHub Issues**: Bug reports and feature requests
- **GitHub Discussions**: Questions and community discussions
- **Pull Requests**: Code review and collaboration

### Getting Help

- **Documentation**: [jerinji2016.github.io/fdawg](https://jerinji2016.github.io/fdawg/)
- **GitHub Discussions**: Ask questions and get help
- **Issues**: Report bugs or request features

### Recognition

Contributors are recognized through:
- GitHub contributor graphs
- Release notes acknowledgments
- Community appreciation

## Questions?

If you have questions about contributing:

1. Check the [documentation](https://jerinji2016.github.io/fdawg/)
2. Search existing [issues](https://github.com/Jerinji2016/fdawg/issues)
3. Start a [discussion](https://github.com/Jerinji2016/fdawg/discussions)
4. Open a new issue with the question template

Thank you for contributing to FDAWG! 🎉

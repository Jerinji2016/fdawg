---
layout: default
title: Build Management
parent: Command Reference
nav_order: 6
description: Comprehensive build system with pre-build setup and artifact organization
permalink: /commands/build/
---

# Build Management

Comprehensive build system for Flutter applications with pre-build setup, artifact organization, and multi-platform support.

---

## Overview

The `build` command provides a comprehensive build management system for Flutter projects. It supports multi-platform builds, pre-build setup steps, environment integration, and organized artifact management with streaming build output.

### Key Features

- **Multi-Platform Support**: Build for Android, iOS, Web, macOS, Linux, and Windows
- **Pre-Build Setup**: Automatic detection and execution of common build tools
- **Environment Integration**: Use environment files with `--dart-define-from-file`
- **Artifact Organization**: Organized output with date-based folders and naming patterns
- **Build Configuration**: Interactive setup wizard with customizable options
- **Streaming Output**: Real-time build progress and output streaming
- **Web Interface**: Modern UI for visual build management and monitoring
- **Parallel Builds**: Optional parallel execution for faster builds (experimental)

---

## Command Usage

### Setup Commands

```bash
# Interactive build configuration wizard
fdawg build setup

# Quick setup with default configuration
fdawg build setup --default

# Force overwrite existing configuration
fdawg build setup --force
```

### Build Execution

```bash
# Build for specific platforms
fdawg build run --platforms android,ios

# Build all available platforms
fdawg build run --platforms all

# Build with environment
fdawg build run --platforms android --env production

# Show build plan without executing
fdawg build run --platforms android --dry-run

# Skip pre-build steps
fdawg build run --platforms web --skip-pre-build

# Continue building other platforms if one fails
fdawg build run --platforms all --continue-on-error

# Parallel builds (experimental)
fdawg build run --platforms android,ios --parallel
```

### Status and Information

```bash
# Show build status and available artifacts
fdawg build status

# Show available platforms
fdawg build run --platforms help
```

---

## Build Configuration

### Configuration File

Build configuration is stored in `.fdawg/build.yaml` and includes:

- **Metadata**: App name and version sources
- **Pre-Build**: Global and platform-specific setup steps
- **Platforms**: Platform-specific build configurations
- **Artifacts**: Output organization and naming
- **Execution**: Build execution options

### Interactive Setup

The setup wizard guides you through configuration:

1. **Metadata Configuration**: Choose app name and version sources
2. **Pre-Build Detection**: Automatically detect common build tools
3. **Platform Configuration**: Enable/disable platforms and configure build types
4. **Artifact Organization**: Configure output structure and naming
5. **Execution Options**: Set parallel builds and error handling

### Default Configuration

Quick setup creates a default configuration with:
- App name from namer configuration or pubspec.yaml
- Version from pubspec.yaml
- Common pre-build steps (build_runner, flutter_launcher_icons)
- Release builds for all available platforms
- Date-organized artifacts with descriptive naming

---

## Pre-Build Steps

### Automatic Detection

The build system automatically detects and configures:

- **build_runner**: Code generation (`dart run build_runner build`)
- **flutter_launcher_icons**: Icon generation
- **flutter_native_splash**: Splash screen generation
- **Custom scripts**: Scripts in the `scripts/` directory

### Custom Pre-Build Steps

Add custom pre-build steps for:
- Code generation
- Asset processing
- Environment setup
- Dependency installation
- Custom build preparation

### Global vs Platform-Specific

- **Global steps**: Run once before all platform builds
- **Platform-specific steps**: Run before each platform build

---

## Platform Support

### Android

**Build Types:**
- APK (release/debug)
- AAB (Android App Bundle)
- Split APKs by ABI

**Configuration Options:**
- Build mode (release, debug, profile)
- Split per ABI
- Obfuscation and debug info
- Custom arguments

### iOS

**Build Types:**
- Archive (for App Store)
- IPA (for distribution)
- Simulator builds

**Configuration Options:**
- Export method (app-store, development, enterprise)
- Code signing options
- Custom arguments

### Web

**Build Types:**
- Standard web build
- PWA (Progressive Web App)

**Configuration Options:**
- Build mode
- PWA features
- Custom arguments

### Desktop (macOS, Linux, Windows)

**Build Types:**
- Native executables
- Distribution packages

**Configuration Options:**
- Build mode
- Custom arguments
- Platform-specific options

---

## Artifact Management

### Organization Structure

```
build/fdawg-outputs/
├── January-15/              # Date-based folders
│   ├── android/
│   │   ├── release_apk/
│   │   │   ├── MyApp_1.0.0_arm64-v8a.apk
│   │   │   └── MyApp_1.0.0_armeabi-v7a.apk
│   │   └── release_aab/
│   │       └── MyApp_1.0.0_universal.aab
│   ├── ios/
│   │   └── archive/
│   │       └── MyApp_1.0.0_universal.ipa
│   └── web/
│       └── release/
│           └── MyApp_1.0.0_web.zip
└── January-16/
    └── ...
```

### Naming Patterns

Artifacts are named using configurable patterns:
- `{app_name}_{version}_{arch}`: Default pattern
- App name from namer configuration or custom
- Version from pubspec.yaml or custom
- Architecture for split builds

### Cleanup

Automatic cleanup options:
- Keep last N builds
- Remove builds older than X days
- Manual cleanup commands

---

## Environment Integration

### Using Environments

```bash
# Build with specific environment
fdawg build run --platforms android --env production

# Environment file is passed to Flutter as --dart-define-from-file
flutter build apk --dart-define-from-file=.fdawg/env/production.env
```

### Environment Files

Environment files from the `env` command are automatically integrated:
- `.fdawg/env/development.env`
- `.fdawg/env/staging.env`
- `.fdawg/env/production.env`

---

## Web Interface

Access the build management interface through the web dashboard:

```bash
fdawg serve
# Navigate to http://localhost:8080/build
```

### Features

#### Build Configuration
- Visual configuration editor
- Form-based setup with validation
- Configuration preview and editing
- Platform availability detection

#### Build Execution
- Platform selection with checkboxes
- Environment dropdown selection
- Build options (dry-run, parallel, etc.)
- Real-time build progress
- Streaming build output

#### Artifact Management
- Build history and status
- Artifact download links
- File size and type information
- Build logs and details

#### Build Monitoring
- Real-time progress updates
- Platform-specific status
- Build output streaming
- Error handling and recovery

---

## Examples

### Basic Build Workflow

```bash
# 1. Set up build configuration
fdawg build setup

# 2. Build for Android
fdawg build run --platforms android

# 3. Check build status
fdawg build status

# 4. Build for production with environment
fdawg build run --platforms all --env production
```

### Advanced Build Configuration

```bash
# Set up with custom configuration
fdawg build setup

# Build with specific options
fdawg build run \
  --platforms android,ios \
  --env production \
  --continue-on-error \
  --parallel

# Dry run to see build plan
fdawg build run --platforms all --dry-run
```

### CI/CD Integration

```bash
# Set up default configuration in CI
fdawg build setup --default

# Build for release
fdawg build run --platforms android --env production

# Check for artifacts
fdawg build status
```

---

## Best Practices

### Configuration Management

1. **Version control**: Include `.fdawg/build.yaml` in version control
2. **Environment-specific configs**: Use different configs for dev/staging/prod
3. **Documentation**: Document custom build steps and requirements
4. **Testing**: Test build configuration in development first

### Build Optimization

1. **Pre-build steps**: Only include necessary pre-build steps
2. **Parallel builds**: Use parallel builds for faster execution (experimental)
3. **Platform selection**: Build only required platforms
4. **Artifact cleanup**: Configure cleanup to manage disk space

### Development Workflow

1. **Setup once**: Configure build settings early in development
2. **Test regularly**: Run builds frequently to catch issues early
3. **Environment testing**: Test with different environments
4. **Artifact management**: Organize and clean up build artifacts

---

## Troubleshooting

### Common Issues

**Build configuration not found**
```bash
# Set up configuration first
fdawg build setup
```

**Platform not available**
```bash
# Check available platforms
fdawg build run --platforms help
```

**Pre-build step failures**
```bash
# Skip pre-build steps to isolate issues
fdawg build run --platforms android --skip-pre-build
```

**Environment file not found**
```bash
# Check available environments
fdawg env list

# Create environment if needed
fdawg env create production
```

### Build Failures

**Flutter build errors**
- Check Flutter installation and project setup
- Verify platform-specific requirements
- Review build logs for specific errors

**Pre-build step failures**
- Check individual step commands
- Verify required dependencies
- Review step configuration and conditions

**Artifact collection failures**
- Check output directory permissions
- Verify build completed successfully
- Review artifact naming configuration

### Recovery

**Reset build configuration**
```bash
# Remove existing configuration
rm .fdawg/build.yaml

# Set up new configuration
fdawg build setup
```

**Clean build artifacts**
```bash
# Manual cleanup
rm -rf build/fdawg-outputs/

# Flutter clean
flutter clean
```

---

## Related Commands

- [`env`](environment.html) - Manage environment variables for builds
- [`namer`](namer.html) - Configure app names used in build artifacts
- [`serve`](server.html) - Start web interface for visual build management
- [`init`](../installation.html#project-validation) - Validate Flutter project structure

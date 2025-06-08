---
layout: default
title: Bundle ID Management
parent: Command Reference
nav_order: 5
description: Manage bundle ids across all platforms
permalink: /commands/bundler/
---

# Bundle ID Management

Manage bundle identifiers across all Flutter platforms with validation and platform-specific configuration.

---

## Overview

The `bundler` command provides comprehensive bundle ID management for Flutter projects. Bundle IDs (also known as application IDs or package names) uniquely identify your app across different platforms and app stores.

### Key Features

- **Universal Configuration**: Set the same bundle ID across all platforms
- **Platform-Specific Management**: Configure different bundle IDs per platform
- **Format Validation**: Automatic validation of bundle ID format and conventions
- **Cross-Platform Support**: Android, iOS, macOS, Linux, Windows, and Web
- **Backup & Recovery**: Automatic backups with rollback capability
- **Web Interface**: Modern UI for visual bundle ID management

---

## Command Usage

### Basic Commands

```bash
# Get current bundle IDs for all platforms
fdawg bundler get

# Set universal bundle ID for all platforms
fdawg bundler set --universal com.company.app

# Set platform-specific bundle ID
fdawg bundler set --platform android --bundle-id com.company.android

# Validate bundle ID format
fdawg bundler validate com.company.app
```

### Advanced Usage

```bash
# Set different bundle IDs for multiple platforms
fdawg bundler set \
  --platform android --bundle-id com.company.android \
  --platform ios --bundle-id com.company.ios

# Get bundle ID for specific platform
fdawg bundler get --platform android

# Set bundle ID with backup
fdawg bundler set --universal com.company.app --backup
```

---

## Bundle ID Format

### Standard Format

Bundle IDs follow reverse domain notation:
```
com.company.application
```

### Validation Rules

- **Minimum 2 segments**: Must contain at least one dot (e.g., `com.app`)
- **Alphanumeric characters**: Only letters, numbers, and underscores
- **No leading numbers**: Segments cannot start with numbers
- **Maximum length**: 255 characters total
- **Lowercase recommended**: Following platform conventions

### Examples

✅ **Valid Bundle IDs**
```
com.company.app
org.example.myapp
dev.flutter.sample_app
```

❌ **Invalid Bundle IDs**
```
company.app          # Missing domain
com.123company.app   # Segment starts with number
com..app            # Empty segment
com.company.app!     # Invalid character
```

---

## Platform Configuration

### Android

**Files Modified:**
- `android/app/build.gradle` or `android/app/build.gradle.kts`

**Properties:**
- `applicationId`: Main bundle identifier
- `namespace`: Kotlin/Java package namespace

### iOS

**Files Modified:**
- `ios/Runner.xcodeproj/project.pbxproj`

**Properties:**
- `PRODUCT_BUNDLE_IDENTIFIER`: Bundle identifier for App Store

### macOS

**Files Modified:**
- `macos/Runner/Configs/AppInfo.xcconfig`

**Properties:**
- `PRODUCT_BUNDLE_IDENTIFIER`: Bundle identifier for Mac App Store

### Linux

**Files Modified:**
- `linux/CMakeLists.txt`

**Properties:**
- `BINARY_NAME`: Application binary name

### Windows

**Files Modified:**
- `windows/CMakeLists.txt`

**Properties:**
- `project()`: Project name
- `BINARY_NAME`: Application binary name

### Web

**Files Modified:**
- `web/manifest.json`

**Properties:**
- `id`: Web app identifier for PWA

---

## Web Interface

Access the bundle ID management interface through the web dashboard:

```bash
fdawg serve
# Navigate to http://localhost:8080/bundler
```

### Features

#### Universal Bundle ID
- Set the same bundle ID for all available platforms
- Real-time format validation
- Confirmation dialog before applying changes

#### Platform-Specific Management
- Individual bundle ID configuration per platform
- Current bundle ID display with namespace information
- Edit-in-place functionality
- Platform availability detection

#### Validation & Safety
- Automatic format validation before submission
- Detailed error messages for invalid formats
- Confirmation dialogs for all changes
- Automatic backup and rollback on failure

---

## Examples

### Setting Universal Bundle ID

```bash
# Set the same bundle ID for all platforms
fdawg bundler set --universal com.mycompany.myapp
```

This will update the bundle ID for all available platforms in your Flutter project.

### Platform-Specific Configuration

```bash
# Different bundle IDs for different platforms
fdawg bundler set --platform android --bundle-id com.mycompany.android
fdawg bundler set --platform ios --bundle-id com.mycompany.ios
```

### Validation

```bash
# Validate bundle ID format
fdawg bundler validate com.mycompany.myapp
# Output: ✓ Valid bundle ID format

fdawg bundler validate invalid..bundle
# Output: ✗ Bundle ID segment is empty
```

---

## Best Practices

### Naming Conventions

1. **Use reverse domain notation**: Start with your domain in reverse
2. **Keep it simple**: Avoid complex or lengthy names
3. **Be consistent**: Use similar patterns across platforms
4. **Follow platform guidelines**: Respect platform-specific conventions

### Development Workflow

1. **Plan early**: Decide on bundle ID structure before development
2. **Use staging variants**: Different bundle IDs for dev/staging/production
3. **Document changes**: Keep track of bundle ID changes
4. **Test thoroughly**: Verify app functionality after bundle ID changes

### Security Considerations

- **Unique identifiers**: Ensure bundle IDs are unique across app stores
- **Avoid conflicts**: Check for existing apps with similar bundle IDs
- **Protect your domain**: Use domains you own and control

---

## Troubleshooting

### Common Issues

**Bundle ID already exists in app store**
```bash
# Use a different bundle ID
fdawg bundler set --universal com.yourcompany.yourapp.v2
```

**Invalid format errors**
```bash
# Validate first to see specific issues
fdawg bundler validate your.bundle.id
```

**Platform not available**
```bash
# Check which platforms are available
fdawg bundler get
```

### Recovery

If something goes wrong, FDAWG automatically creates backups:

```bash
# Backups are stored in .fdawg/backups/
# Manual restoration may be needed in extreme cases
```

---

## Related Commands

- [`namer`](namer.html) - Manage app names across platforms
- [`serve`](server.html) - Start web interface for visual management
- [`init`](../installation.html#project-validation) - Validate Flutter project structure

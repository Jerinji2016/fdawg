---
layout: default
title: Server Commands
parent: Command Reference
nav_order: 7
description: "Web interface and project validation"
permalink: /commands/server/
---

# Server Commands

The server commands provide web interface capabilities and project validation for FDAWG. These commands allow you to launch a modern web dashboard for visual project management and validate Flutter project structure.

## Overview

FDAWG includes two main server-related commands:
- `serve` - Launch the web interface for project management
- `init` - Validate Flutter project structure

## serve - Web Interface

Starts a web server that provides a modern, responsive interface for managing your Flutter project.

```bash
fdawg serve [directory] [--port <port>]
```

**Parameters:**
- `[directory]`: Optional path to Flutter project (defaults to current directory)
- `--port, -p`: Port number for the web server (default: 8080)

**Examples:**
```bash
# Start server for current directory on default port (8080)
fdawg serve

# Start server for specific project
fdawg serve /path/to/my/flutter/project

# Start server on custom port
fdawg serve --port 3000
fdawg serve -p 3000

# Start server for specific project on custom port
fdawg serve --port 3000 /path/to/my/flutter/project
fdawg serve -p 3000 /path/to/my/flutter/project
```

**Important Note:** The port flag must come before the directory argument:
```bash
# ✅ Correct
fdawg serve --port 3000 /path/to/project

# ❌ Incorrect (port will be ignored)
fdawg serve /path/to/project --port 3000
```

### Web Interface Features

The web interface provides comprehensive project management capabilities:

#### 📊 Overview Dashboard
- Project information and metadata
- Dependency overview (dependencies and dev_dependencies)
- Asset summary and statistics
- Quick project health check

#### 🔧 Environment Management
- Visual environment variable management
- Create, edit, and delete environment files
- Add, update, and remove variables
- Generate Dart environment files
- Real-time validation and feedback

#### 📦 Asset Management
- Drag-and-drop asset upload
- Visual asset organization by type
- Asset migration and cleanup tools
- Generate Dart asset files
- Preview and manage existing assets

#### 🌍 Localization Management
- Visual translation management
- Inline editing of translation keys
- Google Translate integration
- Add and remove languages
- Hierarchical translation key organization

#### 🏷️ App Namer
- Cross-platform app naming interface
- Real-time preview of changes
- Platform availability detection
- Visual backup and rollback management
- Universal and platform-specific naming

#### 🚀 Run Configurations (Coming Soon)
- Android Studio run configuration management
- VS Code launch configuration management
- Custom run configuration templates

#### ⚡ Fastlane Integration (Coming Soon)
- Fastlane configuration management
- Build and deployment automation
- Certificate and provisioning profile management

### Web Interface Benefits

**Real-time Updates:**
- All changes are applied immediately
- Live validation and error checking
- Instant feedback on operations

**User Experience:**
- Responsive design for mobile and desktop
- Toast notifications for user feedback
- Confirmation dialogs for destructive actions
- Auto-dismissing notifications

**Safety Features:**
- Automatic backups before destructive operations
- Rollback capabilities
- Project validation before changes
- Error handling with clear messages

### Accessing the Web Interface

1. **Start the server:**
   ```bash
   fdawg serve
   ```

2. **Open your browser:**
   - Default URL: `http://localhost:8080`
   - Custom port: `http://localhost:<your-port>`

3. **Navigate through the interface:**
   - Use the sidebar navigation to access different features
   - Each section provides comprehensive management tools
   - All changes are saved automatically

## init - Project Validation

Validates whether a directory contains a valid Flutter project structure.

```bash
fdawg init [directory]
```

**Parameters:**
- `[directory]`: Optional path to check (defaults to current directory)

**Examples:**
```bash
# Check current directory
fdawg init

# Check specific directory
fdawg init /path/to/my/flutter/project
```

### Validation Checks

The `init` command performs comprehensive validation:

#### ✅ Required Files
- `pubspec.yaml` - Flutter project configuration
- `lib/` directory - Dart source code
- `lib/main.dart` - Application entry point

#### ✅ Flutter Project Structure
- Valid `pubspec.yaml` format
- Flutter SDK dependency
- Proper project name and version

#### ✅ Platform Support Detection
- Android (`android/` directory)
- iOS (`ios/` directory)
- Web (`web/` directory)
- macOS (`macos/` directory)
- Linux (`linux/` directory)
- Windows (`windows/` directory)

### Example Output

**Valid Flutter Project:**
```bash
$ fdawg init
✓ Valid Flutter project detected
✓ pubspec.yaml found and valid
✓ lib/ directory exists
✓ lib/main.dart found

Project: my_flutter_app (1.0.0+1)
Flutter SDK: ^3.0.0

Supported Platforms:
✓ Android
✓ iOS  
✓ Web
✗ macOS (not configured)
✗ Linux (not configured)
✗ Windows (not configured)

Project is ready for FDAWG management!
```

**Invalid Project:**
```bash
$ fdawg init
✗ Not a Flutter project
✗ pubspec.yaml not found
✗ lib/ directory missing

This directory does not contain a valid Flutter project.
Please navigate to a Flutter project directory or create a new Flutter project:

flutter create my_new_project
cd my_new_project
fdawg init
```

## Common Workflows

### Setting up FDAWG for a new project
```bash
# Navigate to your Flutter project
cd /path/to/my/flutter/project

# Validate the project
fdawg init

# Start the web interface
fdawg serve

# Open browser to http://localhost:8080
```

### Using FDAWG with multiple projects
```bash
# Project 1
fdawg serve --port 8080 /path/to/project1 &

# Project 2  
fdawg serve --port 8081 /path/to/project2 &

# Project 3
fdawg serve --port 8082 /path/to/project3 &

# Access each project on different ports
```

### Development workflow
```bash
# Start FDAWG server
fdawg serve --port 3000

# In another terminal, run Flutter
cd /path/to/project
flutter run

# Use FDAWG web interface for project management
# Use Flutter CLI for development and testing
```

### CI/CD Integration
```bash
# Validate project in CI pipeline
fdawg init || exit 1

# Use CLI commands for automated tasks
fdawg env add BUILD_NUMBER $CI_BUILD_NUMBER --env production
fdawg asset generate-dart
fdawg lang generate-dart

# Build Flutter app
flutter build apk --release
```

## Configuration and Customization

### Server Configuration
The web server can be customized through command-line options:

```bash
# Custom port
fdawg serve --port 9000

# Specific project directory
fdawg serve /path/to/project

# Both custom port and directory
fdawg serve --port 9000 /path/to/project
```

### Network Access
By default, the server binds to `localhost` only. For network access in development:

**Security Note:** Only enable network access in trusted development environments.

### Browser Compatibility
The web interface supports modern browsers:
- Chrome 80+
- Firefox 75+
- Safari 13+
- Edge 80+

## Troubleshooting

### Common Issues

**Issue:** Port already in use
```bash
Error: Port 8080 is already in use
```
**Solution:** Use a different port with `--port` flag

**Issue:** Project not detected
```bash
Error: Not a Flutter project
```
**Solution:** Ensure you're in a valid Flutter project directory with `pubspec.yaml`

**Issue:** Web interface doesn't load
**Solution:** 
- Check if the server started successfully
- Verify the port is not blocked by firewall
- Try a different port

**Issue:** Changes not saving
**Solution:**
- Check file permissions in project directory
- Ensure project is not read-only
- Verify disk space availability

### Performance Optimization

For large projects:
- Use asset migration to organize files
- Clean up unused environment variables
- Remove unused translation keys
- Optimize asset file sizes

### Security Considerations

- The web server is intended for development use
- Don't expose the server to public networks
- Use in trusted development environments only
- Sensitive environment variables should be managed carefully

## Integration with Other Tools

### IDE Integration
The web interface complements your IDE workflow:
- Use FDAWG for project-level management
- Use your IDE for code development
- Both can work simultaneously

### Version Control
FDAWG respects your version control:
- Generated files can be committed or ignored
- Backup files are automatically ignored
- Configuration changes are tracked

### Flutter Tools
FDAWG works alongside Flutter CLI:
- Use `flutter` commands for development
- Use `fdawg` commands for project management
- Both tools can be used together

---

**Next:** Learn about the [Web Interface]({{ '/web-interface/' | relative_url }}) in detail or explore [Development Setup]({{ '/development/' | relative_url }}).

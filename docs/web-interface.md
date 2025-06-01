---
layout: page
title: Web Interface Guide
permalink: /web-interface/
---

# Web Interface Guide

FDAWG's web interface provides a modern, intuitive dashboard for managing your Flutter projects. It offers visual tools for all FDAWG features with real-time updates, drag-and-drop functionality, and comprehensive project management capabilities.

## Getting Started

### Launching the Web Interface

```bash
# Start the web server (default port 8080)
fdawg serve

# Custom port
fdawg serve --port 3000

# Specific project
fdawg serve /path/to/your/flutter/project
```

Once started, open your browser to:
- Default: `http://localhost:8080`
- Custom port: `http://localhost:<your-port>`

### System Requirements

- **Modern Browser**: Chrome 80+, Firefox 75+, Safari 13+, Edge 80+
- **JavaScript**: Must be enabled
- **Network**: Local network access (localhost)

## Interface Overview

### Navigation Structure

The web interface uses a sidebar navigation with the following sections:

- **📊 Overview** - Project dashboard and statistics
- **🔧 Environment** - Environment variable management
- **📦 Assets** - Asset management and organization
- **🌍 Localizations** - Translation management
- **🏷️ App Namer** - Cross-platform app naming
- **⚡ Fastlane** - Build automation (coming soon)
- **🚀 Run Configs** - IDE configuration management (coming soon)

### Common UI Elements

**Toast Notifications:**
- Success messages (green)
- Error messages (red)
- Warning messages (yellow)
- Info messages (blue)
- Auto-dismiss after 5 seconds

**Confirmation Dialogs:**
- Appear for destructive actions
- Require explicit confirmation
- Show details about the action

**Loading States:**
- Spinner indicators during operations
- Progress bars for file uploads
- Disabled states during processing

## Features

### 📊 Overview Dashboard

The overview provides a comprehensive project summary:

#### Project Information
- **Project Name**: From pubspec.yaml
- **Version**: Current app version
- **Description**: Project description
- **Flutter SDK**: Required Flutter version

#### Dependencies Overview
- **Dependencies**: Production dependencies with versions
- **Dev Dependencies**: Development dependencies
- **Dependency Count**: Total number of dependencies

#### Asset Summary
- **Asset Count**: Total number of project assets
- **Asset Types**: Breakdown by type (images, audio, etc.)
- **Asset Organization**: Status of asset organization

#### Quick Actions
- **Generate Dart Files**: Quick access to code generation
- **Project Validation**: Check project health
- **Open in IDE**: Launch external IDE (if configured)

### 🔧 Environment Management

Visual environment variable management with full CRUD operations:

#### Environment File Management
- **List Environments**: View all environment files
- **Create Environment**: Add new environment files
- **Copy Environment**: Duplicate existing environments
- **Delete Environment**: Remove environment files

#### Variable Management
- **Add Variables**: Create new environment variables
- **Edit Variables**: Modify existing values inline
- **Remove Variables**: Delete variables with confirmation
- **Bulk Operations**: Select multiple variables for batch actions

#### Features
- **Real-time Validation**: Immediate feedback on variable names
- **Auto-completion**: Suggestions for common variable names
- **Search and Filter**: Find variables quickly
- **Export/Import**: Backup and restore environment files

#### Code Generation
- **Generate Dart**: Create `lib/generated/environment.dart`
- **Preview Code**: See generated code before saving
- **Auto-update**: Regenerate when variables change

### 📦 Asset Management

Comprehensive asset management with visual organization:

#### Asset Upload
- **Drag-and-Drop**: Drop files directly onto the interface
- **File Browser**: Traditional file selection
- **Bulk Upload**: Multiple files at once
- **Type Detection**: Automatic asset type recognition

#### Asset Organization
- **Visual Grid**: Thumbnail view of assets
- **Type Filtering**: Filter by asset type
- **Search**: Find assets by name or type
- **Sorting**: Sort by name, size, date, or type

#### Asset Operations
- **Preview**: View images, play audio/video
- **Rename**: Change asset names
- **Move**: Reorganize into different folders
- **Delete**: Remove assets with confirmation

#### Migration Tools
- **Asset Migration**: Organize legacy assets automatically
- **Backup Creation**: Automatic backup before migration
- **Progress Tracking**: Visual progress during migration
- **Rollback**: Undo migration if needed

#### Code Generation
- **Generate Dart**: Create `lib/generated/assets.dart`
- **Constant Names**: Auto-generated constant names
- **Type Organization**: Group constants by asset type

### 🌍 Localization Management

Advanced translation management with Google Translate integration:

#### Language Management
- **Add Languages**: Support for 50+ languages
- **Remove Languages**: Delete language support
- **Language Status**: See translation completion
- **Default Language**: Manage fallback language

#### Translation Editing
- **Inline Editing**: Edit translations directly in the interface
- **Hierarchical View**: Navigate nested translation keys
- **Bulk Editing**: Edit multiple translations at once
- **Key Management**: Add, remove, and organize translation keys

#### Google Translate Integration
- **Auto-translate**: Translate missing keys automatically
- **Batch Translation**: Translate multiple keys at once
- **Language Detection**: Detect source language automatically
- **Quality Indicators**: Show translation confidence

#### Translation Tools
- **Import/Export**: JSON file import and export
- **Search**: Find translations across all languages
- **Validation**: Check for missing translations
- **Statistics**: Translation completion metrics

#### Features
- **Real-time Preview**: See changes immediately
- **Conflict Resolution**: Handle translation conflicts
- **Version History**: Track translation changes
- **Collaboration**: Multi-user editing support

### 🏷️ App Namer

Cross-platform app naming with live preview:

#### Platform Management
- **Platform Detection**: Automatically detect available platforms
- **Universal Naming**: Set same name across all platforms
- **Platform-specific**: Different names per platform
- **Mixed Strategy**: Combine universal and specific naming

#### Visual Interface
- **Live Preview**: See changes before applying
- **Platform Icons**: Visual platform identification
- **Status Indicators**: Show which platforms are configured
- **Validation**: Real-time name validation

#### Safety Features
- **Automatic Backup**: Backup files before changes
- **Rollback**: Undo changes if something goes wrong
- **Confirmation**: Confirm before applying changes
- **Error Handling**: Clear error messages and recovery

#### Naming Strategies
- **Consumer Apps**: User-friendly names
- **Enterprise Apps**: Professional naming
- **Platform Conventions**: Follow platform-specific guidelines
- **Branding**: Consistent brand representation

### ⚡ Fastlane Integration (Coming Soon)

Automated build and deployment management:

#### Configuration Management
- **Fastfile Editing**: Visual Fastfile configuration
- **Lane Management**: Create and manage build lanes
- **Plugin Integration**: Add and configure Fastlane plugins
- **Environment Setup**: Configure build environments

#### Build Automation
- **Build Triggers**: Automated build triggers
- **Deployment Pipelines**: Multi-stage deployment
- **Certificate Management**: Handle signing certificates
- **Store Integration**: App Store and Play Store integration

### 🚀 Run Configurations (Coming Soon)

IDE configuration management:

#### Android Studio
- **Run Configurations**: Manage Android Studio run configs
- **Build Variants**: Configure build variants
- **Debugging**: Debug configuration setup
- **Testing**: Test configuration management

#### VS Code
- **Launch Configurations**: Manage VS Code launch configs
- **Task Configuration**: Build and test tasks
- **Debug Settings**: Debugging configuration
- **Extension Integration**: Flutter extension settings

## User Experience Features

### Responsive Design
- **Mobile Friendly**: Works on tablets and phones
- **Desktop Optimized**: Full-featured desktop experience
- **Adaptive Layout**: Adjusts to screen size
- **Touch Support**: Touch-friendly interface elements

### Accessibility
- **Keyboard Navigation**: Full keyboard support
- **Screen Reader**: Compatible with screen readers
- **High Contrast**: Support for high contrast themes
- **Focus Management**: Clear focus indicators

### Performance
- **Fast Loading**: Optimized for quick startup
- **Efficient Updates**: Only update changed elements
- **Caching**: Smart caching for better performance
- **Lazy Loading**: Load content as needed

### Customization
- **Theme Support**: Light and dark themes
- **Layout Options**: Customizable layout preferences
- **Shortcuts**: Keyboard shortcuts for common actions
- **Preferences**: User preference management

## Best Practices

### Project Management
1. **Start with Overview**: Always check project status first
2. **Use Environment Management**: Organize configurations properly
3. **Migrate Assets**: Use migration tools for better organization
4. **Regular Backups**: Let FDAWG create automatic backups

### Workflow Integration
1. **Complement IDE**: Use alongside your preferred IDE
2. **Version Control**: Commit generated files appropriately
3. **Team Collaboration**: Share environment templates
4. **CI/CD Integration**: Use CLI commands in automation

### Performance Tips
1. **Asset Optimization**: Optimize assets before upload
2. **Translation Management**: Keep translations organized
3. **Regular Cleanup**: Remove unused assets and translations
4. **Monitor Size**: Keep project size manageable

## Troubleshooting

### Common Issues

**Interface doesn't load:**
- Check if server is running
- Verify port is not blocked
- Try different browser
- Clear browser cache

**Changes not saving:**
- Check file permissions
- Verify disk space
- Ensure project is not read-only
- Check network connectivity

**Upload failures:**
- Check file size limits
- Verify file permissions
- Ensure supported file types
- Try smaller batches

**Translation issues:**
- Check Google Translate API limits
- Verify internet connectivity
- Try manual translation
- Check language support

### Performance Issues

**Slow loading:**
- Clear browser cache
- Reduce asset count
- Optimize large files
- Check system resources

**Memory usage:**
- Close unused tabs
- Restart browser
- Reduce concurrent operations
- Check available RAM

## Security Considerations

### Development Use Only
- Intended for development environments
- Not suitable for production deployment
- Use in trusted networks only
- Avoid exposing to public internet

### Data Protection
- Environment variables may contain sensitive data
- Use secure networks for translation services
- Regular backup of important configurations
- Careful handling of API keys and secrets

### Access Control
- Server runs on localhost by default
- No built-in authentication
- Suitable for single-user development
- Consider network security in team environments

---

The web interface provides a comprehensive, user-friendly way to manage your Flutter projects. For command-line alternatives, see the [Command Reference](commands/) documentation.

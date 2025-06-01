---
layout: default
title: App Namer Commands
parent: Command Reference
nav_order: 4
description: "Manage app names across all platforms"
permalink: /commands/namer/
---

# App Namer Commands

The `namer` command group manages your Flutter app's display name across all platforms. It provides both universal naming (same name for all platforms) and platform-specific naming capabilities with automatic backup and rollback features.

## Overview

```bash
fdawg namer [subcommand] [options] [arguments]
```

The namer commands update platform-specific configuration files to set your app's display name consistently across Android, iOS, macOS, Linux, Windows, and Web platforms.

## Available Subcommands

### `get` / `list` - Get Current App Names

Retrieves the current app names from all or specific platforms.

```bash
fdawg namer get [--platforms <platform-list>]
fdawg namer list [--platforms <platform-list>]
```

**Options:**
- `--platforms, -p`: Comma-separated list of platforms (android, ios, macos, linux, windows, web)

**Examples:**
```bash
# Get names from all available platforms
fdawg namer list

# Get names from specific platforms
fdawg namer get --platforms android,ios,web
```

**Example Output:**
```
Current App Names:
✓ Android: My Flutter App
✓ iOS: My Flutter App  
✓ macOS: My Flutter App
✓ Linux: my_flutter_app
✓ Windows: my_flutter_app
✓ Web: My Flutter App

Platforms found: 6/6
```

### `set` - Set App Names

Sets app names universally across all platforms or for specific platforms.

```bash
fdawg namer set [options]
```

**Universal Options:**
- `--value, -v`: App name to set across all platforms
- `--platforms, -p`: Limit universal setting to specific platforms

**Platform-Specific Options:**
- `--android`: App name for Android only
- `--ios`: App name for iOS only  
- `--macos`: App name for macOS only
- `--linux`: App name for Linux only
- `--windows`: App name for Windows only
- `--web`: App name for Web only

**Examples:**

#### Universal Naming
```bash
# Set same name for all platforms
fdawg namer set --value "My Awesome App"

# Set same name for specific platforms only
fdawg namer set --value "Mobile App" --platforms android,ios
```

#### Platform-Specific Naming
```bash
# Set different names for different platforms
fdawg namer set --android "Android App" --ios "iOS App" --web "Web App"

# Mix universal and platform-specific
fdawg namer set --value "My App" --android "Android Version" --ios "iOS Version"
```

#### Complex Example
```bash
# Set base name for most platforms, with custom names for mobile
fdawg namer set \
  --value "Desktop App" \
  --platforms macos,linux,windows,web \
  --android "Mobile Android" \
  --ios "Mobile iOS"
```

## Platform Configuration Details

The namer command updates the following files for each platform:

### Android
**File:** `android/app/src/main/AndroidManifest.xml`
**Property:** `android:label` attribute in `<application>` tag

```xml
<application
    android:label="My Flutter App"
    android:icon="@mipmap/ic_launcher">
```

### iOS  
**File:** `ios/Runner/Info.plist`
**Properties:** `CFBundleDisplayName` and `CFBundleName`

```xml
<key>CFBundleDisplayName</key>
<string>My Flutter App</string>
<key>CFBundleName</key>
<string>My Flutter App</string>
```

### macOS
**File:** `macos/Runner/Configs/AppInfo.xcconfig`
**Property:** `PRODUCT_NAME` configuration variable

```
PRODUCT_NAME = My Flutter App
```

### Linux
**File:** `linux/CMakeLists.txt`
**Property:** `BINARY_NAME` variable

```cmake
set(BINARY_NAME "my_flutter_app")
```

### Windows
**File:** `windows/CMakeLists.txt`
**Properties:** `project()` name and `BINARY_NAME` variable

```cmake
project(my_flutter_app LANGUAGES CXX)
set(BINARY_NAME "my_flutter_app")
```

### Web
**Files:** `web/manifest.json` and `web/index.html`

**manifest.json:**
```json
{
  "name": "My Flutter App",
  "short_name": "My Flutter App"
}
```

**index.html:**
```html
<title>My Flutter App</title>
<meta name="apple-mobile-web-app-title" content="My Flutter App">
```

## Safety Features

The namer command includes comprehensive safety features:

### 1. Automatic Backups
Before making any changes, FDAWG creates backups of all files that will be modified:

```
.fdawg_backups/
├── namer_backup_20231201_143022/
│   ├── android_app_src_main_AndroidManifest.xml
│   ├── ios_Runner_Info.plist
│   ├── macos_Runner_Configs_AppInfo.xcconfig
│   ├── linux_CMakeLists.txt
│   ├── windows_CMakeLists.txt
│   ├── web_manifest.json
│   └── web_index.html
```

### 2. Automatic Rollback
If any error occurs during the naming process, FDAWG automatically restores all files from the backup:

```
Error updating iOS configuration
Rolling back all changes...
✓ Restored android/app/src/main/AndroidManifest.xml
✓ Restored ios/Runner/Info.plist
All changes have been rolled back successfully.
```

### 3. Platform Validation
FDAWG only attempts to update platforms that exist in your project:

```
Checking available platforms...
✓ Android platform found
✓ iOS platform found
✗ macOS platform not found (skipping)
✓ Web platform found

Updating 3 platforms...
```

### 4. Flutter Project Validation
Ensures you're running the command in a valid Flutter project before making any changes.

## Best Practices

### 1. Consistent Naming Strategy
```bash
# For consumer apps - use friendly names
fdawg namer set --value "My Awesome App"

# For enterprise apps - use descriptive names
fdawg namer set --value "Company Mobile Portal"

# For platform-specific branding
fdawg namer set \
  --android "Android App" \
  --ios "iOS App" \
  --web "Web Portal"
```

### 2. Naming Conventions by Platform

**Mobile Platforms (Android/iOS):**
- Use user-friendly names
- Keep under 30 characters
- Avoid special characters

**Desktop Platforms (macOS/Linux/Windows):**
- Can use longer, more descriptive names
- Consider including company name
- Follow platform conventions

**Web Platform:**
- Use SEO-friendly names
- Consider PWA installation name
- Keep short for browser tabs

### 3. Testing Changes
```bash
# Check current names before changing
fdawg namer list

# Make changes
fdawg namer set --value "New App Name"

# Verify changes
fdawg namer list

# Test app launch on each platform
flutter run -d android
flutter run -d ios
flutter run -d web
```

## Common Workflows

### Setting up app names for a new project
```bash
# Check what platforms are available
fdawg namer list

# Set universal name
fdawg namer set --value "My New Flutter App"

# Verify changes
fdawg namer list
```

### Rebranding an existing app
```bash
# Backup current state (automatic, but good to check)
fdawg namer list

# Update to new brand name
fdawg namer set --value "New Brand Name"

# Test on all platforms
flutter clean
flutter pub get
# Test builds for each platform
```

### Platform-specific naming for different markets
```bash
# Different names for different regions/platforms
fdawg namer set \
  --android "Mobile App - Android" \
  --ios "Mobile App - iOS" \
  --web "Web Portal" \
  --windows "Desktop Application"
```

### Preparing for app store releases
```bash
# Set production names
fdawg namer set --value "Production App Name"

# Verify all platforms before release
fdawg namer list

# Build release versions
flutter build apk --release
flutter build ios --release
flutter build web --release
```

## Troubleshooting

### Common Issues

**Issue:** "Platform not found" warnings
- **Solution:** Normal if your project doesn't support all platforms. FDAWG skips missing platforms automatically.

**Issue:** Permission denied errors
- **Solution:** Ensure you have write permissions to the project directory and platform-specific folders.

**Issue:** Changes not reflected in app
- **Solution:** Run `flutter clean` and `flutter pub get` after changing app names, then rebuild.

**Issue:** Backup restoration needed
- **Solution:** If something goes wrong, backups are in `.fdawg_backups/` directory. You can manually restore files if needed.

### Manual Backup Restoration
If automatic rollback fails, you can manually restore from backups:

```bash
# Find the latest backup
ls -la .fdawg_backups/

# Restore specific files manually
cp .fdawg_backups/namer_backup_*/android_app_src_main_AndroidManifest.xml android/app/src/main/AndroidManifest.xml
```

## Integration with Web Interface

The FDAWG web interface provides visual app naming:
- Real-time preview of names across platforms
- Platform availability detection
- Visual feedback for changes
- Backup and rollback management

Access via: `fdawg serve` → App Namer tab

## Advanced Usage

### Scripting and Automation
```bash
# Use in CI/CD pipelines
fdawg namer set --value "$APP_NAME_ENV_VAR"

# Conditional naming based on environment
if [ "$ENVIRONMENT" = "production" ]; then
  fdawg namer set --value "Production App"
else
  fdawg namer set --value "Development App"
fi
```

### Batch Operations
```bash
# Set different names for development vs production
fdawg namer set --value "Dev App" --platforms android,ios
fdawg namer set --value "Development Portal" --platforms web,windows,macos
```

---

**Next:** Learn about [Server Commands]({{ '/commands/server/' | relative_url }}) for web interface and project validation.

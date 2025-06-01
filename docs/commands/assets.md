---
layout: page
title: Asset Commands
permalink: /commands/assets/
---

# Asset Commands

The `asset` command group provides comprehensive asset management for Flutter projects. FDAWG helps organize, manage, and generate Dart code for easy access to your project assets.

## Overview

```bash
fdawg asset [subcommand] [options] [arguments]
```

FDAWG automatically categorizes assets by type and updates your `pubspec.yaml` file accordingly. It supports images, animations, audio files, videos, JSON files, SVGs, fonts, and miscellaneous assets.

## Available Subcommands

### `list` - List Project Assets

Lists all assets currently registered in your Flutter project.

```bash
fdawg asset list
```

**Example Output:**
```
Project Assets:
Images (5):
  - assets/images/logo.png
  - assets/images/background.jpg
  - assets/images/icons/home.png

Animations (2):
  - assets/animations/loading.json
  - assets/animations/success.json

Audio (1):
  - assets/audio/notification.mp3

Total: 8 assets
```

### `add` - Add Asset to Project

Adds an asset to your Flutter project and updates the `pubspec.yaml` file.

```bash
fdawg asset add <asset-path> [--type <asset-type>]
```

**Parameters:**
- `<asset-path>`: Path to the asset file (relative or absolute)
- `--type, -t`: Asset type (images, animations, audio, videos, json, svgs, fonts, misc)

**Asset Types:**
- `images`: PNG, JPG, JPEG, GIF, BMP, WEBP
- `animations`: JSON (Lottie), GIF
- `audio`: MP3, WAV, AAC, OGG
- `videos`: MP4, AVI, MOV, WMV
- `json`: JSON files (non-animation)
- `svgs`: SVG vector graphics
- `fonts`: TTF, OTF font files
- `misc`: Other file types

**Examples:**
```bash
# Add image (type auto-detected)
fdawg asset add path/to/logo.png

# Add with specific type
fdawg asset add path/to/data.json --type json

# Add animation
fdawg asset add path/to/loading.json --type animations

# Add font
fdawg asset add path/to/custom-font.ttf --type fonts
```

### `remove` - Remove Asset

Removes an asset from your project and updates the `pubspec.yaml` file.

```bash
fdawg asset remove <asset-name> [--type <asset-type>]
```

**Parameters:**
- `<asset-name>`: Name of the asset file to remove
- `--type, -t`: Asset type (optional, searches all types if not specified)

**Examples:**
```bash
# Remove asset (searches all types)
fdawg asset remove logo.png

# Remove with specific type
fdawg asset remove custom-font.ttf --type fonts

# Remove animation
fdawg asset remove loading.json --type animations
```

### `migrate` - Organize Assets

Migrates and organizes existing assets into structured folders by type. This is useful for cleaning up legacy projects or reorganizing assets.

```bash
fdawg asset migrate
```

**What it does:**
- Creates organized folder structure (`assets/images/`, `assets/audio/`, etc.)
- Moves assets to appropriate folders based on file type
- Creates backup in `assets.backup/` directory
- Updates `pubspec.yaml` with new asset paths
- Removes empty directories after migration
- Excludes `translations/` directory to avoid conflicts with localization

**Folder Structure Created:**
```
assets/
├── images/          # PNG, JPG, JPEG, GIF, BMP, WEBP
├── animations/      # JSON (Lottie), animated GIFs
├── audio/           # MP3, WAV, AAC, OGG
├── videos/          # MP4, AVI, MOV, WMV
├── json/            # JSON data files
├── svgs/            # SVG vector graphics
├── fonts/           # TTF, OTF font files
└── misc/            # Other file types
```

**Example:**
```bash
fdawg asset migrate
```

**Output:**
```
Creating backup in assets.backup/...
Migrating assets to organized folders...
  - Moved logo.png to assets/images/
  - Moved loading.json to assets/animations/
  - Moved notification.mp3 to assets/audio/
Updated pubspec.yaml with new asset paths
Migration completed successfully!
```

### `generate-dart` - Generate Dart Asset File

Generates a Dart file with constants for all project assets, making it easy to reference assets in your code.

```bash
fdawg asset generate-dart
```

This creates `lib/generated/assets.dart` with a class containing constants for all assets.

**Generated Code Example:**
```dart
class Assets {
  // Images
  static const String logoImage = 'assets/images/logo.png';
  static const String backgroundImage = 'assets/images/background.jpg';
  static const String homeIcon = 'assets/images/icons/home.png';
  
  // Animations
  static const String loadingAnimation = 'assets/animations/loading.json';
  static const String successAnimation = 'assets/animations/success.json';
  
  // Audio
  static const String notificationSound = 'assets/audio/notification.mp3';
  
  // Fonts
  static const String customFont = 'assets/fonts/custom-font.ttf';
}
```

## Asset Type Detection

FDAWG automatically detects asset types based on file extensions:

| Type | Extensions |
|------|------------|
| **Images** | `.png`, `.jpg`, `.jpeg`, `.gif`, `.bmp`, `.webp` |
| **Animations** | `.json` (Lottie), `.gif` (animated) |
| **Audio** | `.mp3`, `.wav`, `.aac`, `.ogg`, `.m4a` |
| **Videos** | `.mp4`, `.avi`, `.mov`, `.wmv`, `.mkv` |
| **JSON** | `.json` (non-animation data) |
| **SVGs** | `.svg` |
| **Fonts** | `.ttf`, `.otf`, `.woff`, `.woff2` |
| **Misc** | All other file types |

## Best Practices

### 1. Organized Asset Structure
```bash
# Start with migration for existing projects
fdawg asset migrate

# Add new assets with proper organization
fdawg asset add new-logo.png --type images
fdawg asset add loading-spinner.json --type animations
```

### 2. Asset Naming Conventions
- Use descriptive names: `user_profile_placeholder.png` instead of `img1.png`
- Use snake_case or kebab-case consistently
- Include size indicators for images: `logo_small.png`, `logo_large.png`
- Group related assets: `icon_home.png`, `icon_settings.png`

### 3. Using Generated Dart Code
After running `generate-dart`, use assets in your Flutter app:

```dart
import 'package:your_app/generated/assets.dart';

Widget build(BuildContext context) {
  return Column(
    children: [
      Image.asset(Assets.logoImage),
      Lottie.asset(Assets.loadingAnimation),
      // Audio player
      AudioPlayer().play(AssetSource(Assets.notificationSound)),
    ],
  );
}
```

### 4. Asset Optimization
- Optimize images before adding them to your project
- Use appropriate formats (PNG for transparency, JPG for photos)
- Consider using SVGs for scalable graphics
- Compress audio/video files appropriately

## Common Workflows

### Setting up assets for a new project
```bash
# Add initial assets
fdawg asset add assets/logo.png
fdawg asset add assets/loading.json --type animations
fdawg asset add assets/notification.mp3 --type audio

# Generate Dart code
fdawg asset generate-dart
```

### Organizing an existing project
```bash
# Migrate existing assets
fdawg asset migrate

# Add any new assets
fdawg asset add new-icon.png

# Regenerate Dart code
fdawg asset generate-dart
```

### Adding platform-specific assets
```bash
# Add different resolution images
fdawg asset add assets/images/logo.png
fdawg asset add assets/images/2.0x/logo.png
fdawg asset add assets/images/3.0x/logo.png

# Add platform-specific sounds
fdawg asset add assets/audio/notification_ios.mp3
fdawg asset add assets/audio/notification_android.mp3
```

### Cleaning up assets
```bash
# Remove unused assets
fdawg asset remove old-logo.png
fdawg asset remove unused-sound.mp3

# Regenerate Dart code
fdawg asset generate-dart
```

## File Structure

After using asset commands, your project structure will look like:

```
your_flutter_project/
├── assets/
│   ├── images/
│   │   ├── logo.png
│   │   ├── background.jpg
│   │   └── icons/
│   │       └── home.png
│   ├── animations/
│   │   ├── loading.json
│   │   └── success.json
│   ├── audio/
│   │   └── notification.mp3
│   ├── fonts/
│   │   └── custom-font.ttf
│   └── misc/
│       └── data.txt
├── assets.backup/          # Created by migrate command
├── lib/
│   └── generated/
│       └── assets.dart     # Generated by generate-dart
└── pubspec.yaml            # Updated automatically
```

## Integration with pubspec.yaml

FDAWG automatically updates your `pubspec.yaml` file when adding or removing assets:

```yaml
flutter:
  assets:
    - assets/images/
    - assets/animations/
    - assets/audio/
    - assets/json/
    - assets/svgs/
    - assets/misc/
  
  fonts:
    - family: CustomFont
      fonts:
        - asset: assets/fonts/custom-font.ttf
```

---

**Next:** Learn about [Localization Commands](localization.html) for managing translations.

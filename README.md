# FDAWG - Flutter Development Assistant with Go

FDAWG is a CLI tool designed to help manage Flutter projects. It provides various commands to streamline Flutter development workflows including asset management, environment configuration, localization, app naming across platforms, and includes a modern web interface for project management.

## Installation

```bash
# Clone the repository
git clone https://github.com/Jerinji2016/fdawg.git

# Build the binary
cd fdawg
go build -o fdawg ./cmd/flutter-manager

# Optional: Move to a directory in your PATH
mv fdawg /usr/local/bin/
```

## Usage

```bash
fdawg [command] [options]
```

## Available Commands

| Command | Description |
|---------|-------------|
| `serve [directory]` | Start a web server for project management |
| `init [directory]`  | Check if the specified directory is a Flutter project |
| `env list` | List all environment files in the `.environment` directory |
| `env show <env-name>` | Show all variables in a specific environment file |
| `env create <env-name>` | Create a new environment file |
| `env add <key> <value>` | Add or update a variable in an environment file |
| `env delete <env-name>` | Delete an environment file |
| `env remove <key>` | Remove a variable from an environment file |
| `env generate-dart` | Generate a Dart environment file with all environment variables |
| `asset add <asset-path>` | Add an asset to the project |
| `asset remove <asset-name>` | Remove an asset from the project |
| `asset list` | List all assets in the project |
| `asset generate-dart` | Generate a Dart asset file with all assets |
| `asset migrate` | Migrate assets to organized folders by type and clean up empty directories |
| `lang init` | Initialize localization using easy_localization package |
| `lang add <language-code>` | Add support for a new language |
| `lang remove <language-code>` | Remove support for a language |
| `lang insert <key>` | Add a new translation key to all languages |
| `lang delete <key>` | Delete a translation key from all languages |
| `lang list` | List all supported languages in the project |
| `namer get` | Get current app names from all or specific platforms |
| `namer set` | Set app names universally or for specific platforms |
| `namer list` | List current app names for all available platforms |

### Environment Command Options

- `env show <env-name>`:
  - `<env-name>`: Name of the environment file

- `env create <env-name>`:
  - `<env-name>`: Name of the environment file
  - `--copy` or `-c`: Copy variables from an existing environment file

- `env add <key> <value>`:
  - `<key>`: Variable key (must start with a letter or underscore)
  - `<value>`: Variable value
  - `--env` or `-e`: Environment file to add the variable to (default: development)

- `env delete <env-name>`:
  - `<env-name>`: Name of the environment file to delete

- `env remove <key>`:
  - `<key>`: Variable key to remove
  - `--env` or `-e`: Environment file to remove the variable from (default: development)

### Environment Command Examples

List all environment files:

```bash
fdawg env list
```

Show variables in the development environment:

```bash
fdawg env show development
```

Create a new environment file:

```bash
fdawg env create production
```

Create a new environment file by copying from an existing one:

```bash
fdawg env create staging --copy development
```

Add a variable to the default (development) environment:

```bash
fdawg env add API_URL https://api.example.com
```

Add a variable to a specific environment:

```bash
fdawg env add DEBUG_MODE false --env production
```

Remove a variable from an environment:

```bash
fdawg env remove API_URL --env staging
```

Delete an environment file:

```bash
fdawg env delete staging
```

Generate the Dart environment file:

```bash
fdawg env generate-dart
```

### Asset Command Options

- `asset add <asset-path>`:
  - `--type, -t`: Type of asset (images, animations, audio, videos, json, svgs, misc)
  - If no type is specified, it will be determined from the file extension

- `asset remove <asset-name>`:
  - `--type, -t`: Type of asset (images, animations, audio, videos, json, svgs, misc)
  - If no type is specified, it will search for the asset in all directories

- `asset migrate`:
  - Migrates assets to organized folders by type (images, animations, audio, etc.)
  - Creates a backup of all assets in assets.backup directory
  - Automatically categorizes assets based on file type and content
  - Recursively processes all files in nested directories
  - Removes empty directories after migration

### Asset Command Examples

Add an image asset:

```bash
fdawg asset add path/to/image.png
```

Add an asset with a specific type:

```bash
fdawg asset add path/to/file.json --type animations
```

Remove an asset:

```bash
fdawg asset remove image.png
```

Remove an asset with a specific type:

```bash
fdawg asset remove font.ttf --type fonts
```

List all assets:

```bash
fdawg asset list
```

Generate the Dart asset file:

```bash
fdawg asset generate-dart
```

Migrate assets to organized folders:

```bash
fdawg asset migrate
```

### Localization Command Options

- `lang init`:
  - Initializes localization using the easy_localization package
  - Creates the translations directory and default English (US) translation file
  - Updates pubspec.yaml to add the easy_localization dependency
  - Updates main.dart to initialize easy_localization

- `lang add <language-code>`:
  - `<language-code>`: ISO language code (e.g., 'en', 'es', 'fr')
  - For country-specific variants, use format like 'en_US', 'pt_BR'

- `lang remove <language-code>`:
  - `<language-code>`: ISO language code of the language to remove
  - Cannot remove the default language (en_US)

- `lang insert <key>`:
  - `<key>`: Translation key in dot notation (e.g., 'app.welcome', 'auth.login.title')
  - Will prompt for translation values for each supported language

- `lang delete <key>`:
  - `<key>`: Translation key to delete from all language files

- `lang list`:
  - Lists all supported languages in the project

### Localization Command Examples

Initialize localization:

```bash
fdawg lang init
```

Add support for Spanish:

```bash
fdawg lang add es
```

Add support for French (France):

```bash
fdawg lang add fr_FR
```

Remove support for a language:

```bash
fdawg lang remove es
```

Add a new translation key:

```bash
fdawg lang insert app.welcome
```

Delete a translation key:

```bash
fdawg lang delete app.welcome
```

List all supported languages:

```bash
fdawg lang list
```

### App Namer Command Options

The `namer` command manages your Flutter app's display name across all platforms. It supports both universal naming (same name for all platforms) and platform-specific naming.

- `namer get`:
  - `--platforms, -p`: Specify platforms to get names from (android, ios, macos, linux, windows, web)
  - If no platforms specified, gets names from all available platforms

- `namer set`:
  - `--value, -v`: App name to set universally across all platforms
  - `--platforms, -p`: Limit universal setting to specific platforms
  - `--android`: App name for Android platform only
  - `--ios`: App name for iOS platform only
  - `--macos`: App name for macOS platform only
  - `--linux`: App name for Linux platform only
  - `--windows`: App name for Windows platform only
  - `--web`: App name for Web platform only

- `namer list`:
  - Lists current app names for all available platforms (alias for `namer get`)

### App Namer Command Examples

Get current app names for all platforms:

```bash
fdawg namer list
# or
fdawg namer get
```

Get app names for specific platforms:

```bash
fdawg namer get --platforms android,ios,web
```

Set the same name for all platforms:

```bash
fdawg namer set --value "My Awesome App"
```

Set the same name for specific platforms only:

```bash
fdawg namer set --value "My App" --platforms android,ios
```

Set different names for different platforms:

```bash
fdawg namer set --android "Android App" --ios "iOS App" --web "Web App"
```

Mix universal and platform-specific naming:

```bash
# Set "My App" for most platforms, but use custom names for Android and iOS
fdawg namer set --value "My App" --android "Android Version" --ios "iOS Version"
```

### Platform Configuration Details

The `namer` command updates the following files for each platform:

- **Android**: `android/app/src/main/AndroidManifest.xml`
  - Updates the `android:label` attribute in the `<application>` tag

- **iOS**: `ios/Runner/Info.plist`
  - Updates both `CFBundleDisplayName` and `CFBundleName` properties

- **macOS**: `macos/Runner/Configs/AppInfo.xcconfig`
  - Updates the `PRODUCT_NAME` configuration variable

- **Linux**: `linux/CMakeLists.txt`
  - Updates the `BINARY_NAME` variable

- **Windows**: `windows/CMakeLists.txt`
  - Updates both the `project()` name and `BINARY_NAME` variable

- **Web**: `web/manifest.json` and `web/index.html`
  - Updates `name` and `short_name` in manifest.json
  - Updates `<title>` and `apple-mobile-web-app-title` in index.html

### Safety Features

The namer command includes several safety features:

- **Automatic backups**: Creates backups of all modified files before making changes
- **Rollback capability**: Automatically restores from backups if any errors occur
- **Platform validation**: Only attempts to update platforms that exist in your project
- **Flutter project validation**: Ensures you're running the command in a valid Flutter project

## Server and Init Command Examples

Start the web server for the current directory:

```bash
fdawg serve
```

Start the web server for a specific Flutter project:

```bash
fdawg serve /path/to/my/flutter/project
```

Start the web server on a custom port:

```bash
fdawg serve --port 3000
# or
fdawg serve -p 3000
```

Start the web server for a specific project on a custom port:

```bash
fdawg serve --port 3000 /path/to/my/flutter/project
# or
fdawg serve -p 3000 /path/to/my/flutter/project

# NOTE: The following will not take the port details
fdawg serve /path/to/my/flutter/project --port 3000
# or
fdawg serve /path/to/my/flutter/project -p 3000
```

Check if the current directory is a Flutter project:

```bash
fdawg init
```

Check if a specific directory is a Flutter project:

```bash
fdawg init /path/to/my/flutter/project
```

### Web Interface Features

The web server provides a modern, responsive interface for managing your Flutter project with the following features:

- **Overview**: Project information and quick stats
- **Environment**: Manage environment variables with a visual interface
- **Assets**: Upload, organize, and manage project assets with drag-and-drop support
- **Localizations**: Manage translations with inline editing and Google Translate integration
- **App Namer**: Set app display names across all platforms with real-time preview
- **Fastlane**: (Coming soon) Fastlane configuration management
- **Run Configs**: (Coming soon) Run configuration management

All web interface features include:

- Real-time updates and validation
- Toast notifications for user feedback
- Confirmation dialogs for destructive actions
- Responsive design for mobile and desktop
- Auto-dismissing notifications

## Project Structure

```text
fdawg/
├── cmd/                  # Application entry points
│   └── flutter-manager/  # Main CLI application
├── internal/             # Private application code
│   ├── commands/         # CLI command implementations
│   └── server/           # Web server implementation
│       ├── web/          # Web assets
│       │   ├── static/   # Static assets
│       │   │   ├── css/  # Compiled CSS
│       │   │   ├── js/   # JavaScript files
│       │   │   └── scss/ # SASS source files
│       │   └── templates/# HTML templates
├── pkg/                  # Reusable packages
│   ├── asset/            # Asset management functionality
│   ├── config/           # Configuration management
│   ├── environment/      # Environment variables management
│   ├── flutter/          # Flutter project utilities
│   ├── localization/     # Localization management
│   ├── namer/            # App name management across platforms
│   ├── translate/        # Translation services
│   └── utils/            # Utility functions
```

## SASS Structure

The project uses SASS for styling with a modular structure:

```text
scss/
├── _variables.scss      # Variables for colors, fonts, sizes, etc.
├── _base.scss           # Base styles and typography
├── _layout.scss         # Layout components (header, sidebar, content)
├── _components.scss     # UI components (buttons, cards, etc.)
├── _responsive.scss     # Media queries for responsive design
└── main.scss            # Main file that imports all partials
```

## Development

### SASS Compilation

The project includes npm scripts for SASS compilation:

```bash
# Install dependencies
npm install

# Compile SASS to CSS
npm run sass

# Watch SASS files for changes
npm run sass:watch
```

### Makefile Commands

```bash
# Build the project
make build

# Compile SASS and build
make all

# Compile SASS
make sass

# Watch SASS files for changes
make sass-watch

# Run the application (compiles SASS first)
make run

# Run in development mode (watches SASS files)
make dev
```

## License

[MIT License](LICENSE)

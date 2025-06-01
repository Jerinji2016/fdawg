---
layout: default
title: Localization Commands
parent: Command Reference
nav_order: 3
description: "Manage translations and multi-language support"
permalink: /commands/localization/
---

# Localization Commands

The `lang` command group provides comprehensive localization management for Flutter projects using the `easy_localization` package. FDAWG simplifies the process of adding multi-language support to your Flutter applications.

## Overview

```bash
fdawg lang [subcommand] [options] [arguments]
```

FDAWG uses the popular `easy_localization` package for Flutter localization, automatically managing translation files, updating dependencies, and configuring your app for multi-language support.

## Available Subcommands

### `init` - Initialize Localization

Sets up localization for your Flutter project using the `easy_localization` package.

```bash
fdawg lang init
```

**What it does:**
- Adds `easy_localization` dependency to `pubspec.yaml`
- Creates `assets/translations/` directory
- Creates default `en_US.json` translation file
- Updates `main.dart` to initialize `easy_localization`
- Configures asset paths in `pubspec.yaml`

**Example Output:**
```
Initializing localization...
✓ Added easy_localization dependency
✓ Created translations directory
✓ Created default en_US.json file
✓ Updated main.dart with localization setup
✓ Updated pubspec.yaml with translation assets
Localization initialized successfully!
```

### `list` - List Supported Languages

Lists all languages currently supported in your project.

```bash
fdawg lang list
```

**Example Output:**
```
Supported Languages:
- en_US (English - United States) [Default]
- es (Spanish)
- fr_FR (French - France)
- de (German)
- ja (Japanese)

Total: 5 languages
```

### `add` - Add Language Support

Adds support for a new language by creating a translation file.

```bash
fdawg lang add <language-code>
```

**Parameters:**
- `<language-code>`: ISO language code (e.g., 'en', 'es', 'fr') or locale code (e.g., 'en_US', 'pt_BR')

**Supported Language Codes:**
- Simple codes: `en`, `es`, `fr`, `de`, `ja`, `ko`, `zh`, `ar`, `hi`, `pt`, `ru`, `it`
- Locale-specific: `en_US`, `en_GB`, `es_ES`, `es_MX`, `fr_FR`, `fr_CA`, `pt_BR`, `pt_PT`, `zh_CN`, `zh_TW`

**Examples:**
```bash
# Add Spanish support
fdawg lang add es

# Add French (France) support
fdawg lang add fr_FR

# Add Portuguese (Brazil) support
fdawg lang add pt_BR
```

### `remove` - Remove Language Support

Removes support for a language by deleting its translation file.

```bash
fdawg lang remove <language-code>
```

**Parameters:**
- `<language-code>`: Language code to remove

**Note:** Cannot remove the default language (`en_US`).

**Examples:**
```bash
# Remove Spanish support
fdawg lang remove es

# Remove French (France) support
fdawg lang remove fr_FR
```

### `insert` - Add Translation Key

Adds a new translation key to all supported language files.

```bash
fdawg lang insert <key>
```

**Parameters:**
- `<key>`: Translation key in dot notation (e.g., 'app.welcome', 'auth.login.title')

The command will prompt you to enter translations for each supported language.

**Example:**
```bash
fdawg lang insert app.welcome
```

**Interactive Prompt:**
```
Adding translation key: app.welcome

Enter translation for en_US (English - United States): Welcome to our app!
Enter translation for es (Spanish): ¡Bienvenido a nuestra aplicación!
Enter translation for fr_FR (French - France): Bienvenue dans notre application !

✓ Added translation key 'app.welcome' to 3 languages
```

### `delete` - Remove Translation Key

Removes a translation key from all language files.

```bash
fdawg lang delete <key>
```

**Parameters:**
- `<key>`: Translation key to delete

**Example:**
```bash
fdawg lang delete app.old_message
```

## Translation File Structure

Translation files are stored as JSON in the `assets/translations/` directory:

```
your_flutter_project/
├── assets/
│   └── translations/
│       ├── en_US.json    # Default language
│       ├── es.json       # Spanish
│       ├── fr_FR.json    # French (France)
│       └── de.json       # German
├── lib/
│   └── main.dart         # Updated with localization setup
└── pubspec.yaml          # Updated with dependencies and assets
```

**Example Translation File (`en_US.json`):**
```json
{
  "app": {
    "title": "My Flutter App",
    "welcome": "Welcome to our app!",
    "description": "This is an amazing Flutter application"
  },
  "auth": {
    "login": {
      "title": "Login",
      "email": "Email Address",
      "password": "Password",
      "button": "Sign In"
    },
    "register": {
      "title": "Create Account",
      "button": "Sign Up"
    }
  },
  "common": {
    "ok": "OK",
    "cancel": "Cancel",
    "save": "Save",
    "delete": "Delete"
  }
}
```

## Using Translations in Your App

After initializing localization, use translations in your Flutter app:

### 1. Import the package
```dart
import 'package:easy_localization/easy_localization.dart';
```

### 2. Use translations in widgets
```dart
class MyWidget extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    return Column(
      children: [
        Text('app.title'.tr()),           // Simple translation
        Text('app.welcome'.tr()),         // Nested translation
        Text('auth.login.title'.tr()),    // Deep nested translation
        
        // With parameters
        Text('welcome_user'.tr(namedArgs: {'name': 'John'})),
        
        // Plural forms
        Text('items_count'.plural(itemCount)),
      ],
    );
  }
}
```

### 3. Change language dynamically
```dart
// Change to Spanish
await context.setLocale(Locale('es'));

// Change to French (France)
await context.setLocale(Locale('fr', 'FR'));

// Get current locale
Locale currentLocale = context.locale;
```

## Best Practices

### 1. Translation Key Organization
Use hierarchical keys with dot notation:

```json
{
  "app": {
    "title": "My App",
    "subtitle": "Welcome"
  },
  "screens": {
    "home": {
      "title": "Home",
      "welcome": "Welcome back!"
    },
    "profile": {
      "title": "Profile",
      "edit": "Edit Profile"
    }
  },
  "common": {
    "buttons": {
      "save": "Save",
      "cancel": "Cancel",
      "delete": "Delete"
    }
  }
}
```

### 2. Naming Conventions
- Use descriptive, hierarchical keys
- Group related translations together
- Use snake_case or camelCase consistently
- Avoid overly deep nesting (max 3-4 levels)

### 3. Translation Management
```bash
# Add all common languages at once
fdawg lang add es
fdawg lang add fr
fdawg lang add de
fdawg lang add ja

# Add common UI elements
fdawg lang insert common.ok
fdawg lang insert common.cancel
fdawg lang insert common.save
fdawg lang insert common.delete

# Add screen-specific translations
fdawg lang insert screens.home.title
fdawg lang insert screens.home.welcome
```

### 4. Handling Plurals and Parameters

**Plurals:**
```json
{
  "items_count": {
    "zero": "No items",
    "one": "One item", 
    "other": "{} items"
  }
}
```

**Parameters:**
```json
{
  "welcome_user": "Welcome, {name}!",
  "items_found": "Found {count} items in {category}"
}
```

## Common Workflows

### Setting up localization for a new project
```bash
# Initialize localization
fdawg lang init

# Add target languages
fdawg lang add es
fdawg lang add fr
fdawg lang add de

# Add basic translations
fdawg lang insert app.title
fdawg lang insert app.welcome
fdawg lang insert common.ok
fdawg lang insert common.cancel
```

### Adding new features with translations
```bash
# Add translations for new screen
fdawg lang insert screens.settings.title
fdawg lang insert screens.settings.theme
fdawg lang insert screens.settings.language
fdawg lang insert screens.settings.notifications
```

### Managing existing translations
```bash
# List current languages
fdawg lang list

# Add new language
fdawg lang add pt_BR

# Remove unused language
fdawg lang remove de

# Update existing translations
fdawg lang delete old.unused.key
fdawg lang insert new.feature.title
```

## Integration with main.dart

After running `fdawg lang init`, your `main.dart` will be updated:

```dart
import 'package:easy_localization/easy_localization.dart';
import 'package:flutter/material.dart';

void main() async {
  WidgetsFlutterBinding.ensureInitialized();
  await EasyLocalization.ensureInitialized();
  
  runApp(
    EasyLocalization(
      supportedLocales: [
        Locale('en', 'US'),
        Locale('es'),
        Locale('fr', 'FR'),
      ],
      path: 'assets/translations',
      fallbackLocale: Locale('en', 'US'),
      child: MyApp(),
    ),
  );
}

class MyApp extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      localizationsDelegates: context.localizationDelegates,
      supportedLocales: context.supportedLocales,
      locale: context.locale,
      home: MyHomePage(),
    );
  }
}
```

## Web Interface Integration

The FDAWG web interface provides visual translation management:
- View all translations in a table format
- Edit translations inline
- Google Translate integration for quick translations
- Add/remove languages visually
- Export/import translation files

Access via: `fdawg serve` → Localizations tab

---

**Next:** Learn about [App Namer Commands](namer.html) for cross-platform app naming.

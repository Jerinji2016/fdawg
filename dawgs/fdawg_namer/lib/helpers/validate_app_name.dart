part of '../fdawg_namer.dart';

/// Validates name against provided [platforms].
///
/// If [platforms] is null, name is validated against all platforms
void _validateAppName(String name, [List<PlatformOptions>? platforms]) {
  platforms ??= PlatformOptions.values.toList();

  final appNamePlatforms = [
    PlatformOptions.windows,
    PlatformOptions.linux,
    PlatformOptions.macos,
  ];

  if (platforms.any((platform) => appNamePlatforms.contains(platform))) {
    validateAppName(name);
  }

  final appLabelPlatforms = [
    PlatformOptions.android,
    PlatformOptions.ios,
    PlatformOptions.web,
  ];

  if (platforms.any((platform) => appLabelPlatforms.contains(platform))) {
    validateAppLabel(name);
  }
}

// Validate app name (Windows, macOS, Linux)
void validateAppName(String name) {
  if (name.isEmpty) {
    throw errorNameEmpty;
  }

// Regex for valid characters: Letters, numbers, spaces, dashes, underscores
  final regex = RegExp(r'^[a-zA-Z0-9_\- ]+$');
  if (!regex.hasMatch(name)) {
    throw errorInvalidCharacters;
  }

  if (name.length > 255) {
    throw errorInvalidLength;
  }

  final startWithSpecialChars = RegExp(r'^[_\- ]').hasMatch(name);
  final endsWithSpecialChars = RegExp(r'[_\- ]$').hasMatch(name);
  if (name.trim() != name || startWithSpecialChars || endsWithSpecialChars) {
    throw errorStartsOrEndsWithSpecialChar;
  }

  if (reservedWindowsNames.contains(name.toUpperCase())) {
    throw errorReservedName;
  }
}

// Validate app label/title (Android, iOS, Web)
void validateAppLabel(String label) {
  if (label.isEmpty) {
    throw errorNameEmpty;
  }

  final regex = RegExp(r'^[a-zA-Z0-9 ]+$');
  if (!regex.hasMatch(label)) {
    throw errorInvalidCharacters;
  }

  if (label.length > 30) {
    throw errorInvalidLength;
  }
}

const reservedWindowsNames = {
  'CON',
  'PRN',
  'AUX',
  'NUL',
  'COM1',
  'COM2',
  'COM3',
  'COM4',
  'COM5',
  'COM6',
  'COM7',
  'COM8',
  'COM9',
  'LPT1',
  'LPT2',
  'LPT3',
  'LPT4',
  'LPT5',
  'LPT6',
  'LPT7',
  'LPT8',
  'LPT9'
};

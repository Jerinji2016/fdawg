import 'package:fdawg_core/fdawg_core.dart';
import 'package:fdawg_namer/common/errors.dart';
import 'package:fdawg_namer/fdawg_namer.dart';
import 'package:flutter_test/flutter_test.dart';

void main() {
  group('validateAppName', () {
    final platforms = [
      PlatformOptions.windows,
      PlatformOptions.linux,
      PlatformOptions.macos,
    ];

    test('Valid app names', () {
      expect(() => FdawgNamer.isValidAppName('MyApp', platforms), returnsNormally);
      expect(() => FdawgNamer.isValidAppName('Cool-App_2024', platforms), returnsNormally);
      expect(() => FdawgNamer.isValidAppName('My App', platforms), returnsNormally);
    });

    test('Throws error for empty app name', () {
      expect(
        () => FdawgNamer.isValidAppName('', platforms),
        throwsA(equals(errorNameEmpty)),
      );
    });

    test('Throws error for invalid characters in app name', () {
      expect(
        () => FdawgNamer.isValidAppName('My@pp', platforms),
        throwsA(equals(errorInvalidCharacters)),
      );
      expect(
        () => FdawgNamer.isValidAppName('My#App', platforms),
        throwsA(equals(errorInvalidCharacters)),
      );
    });

    test('Throws error for reserved Windows names', () {
      expect(
        () => FdawgNamer.isValidAppName('CON', platforms),
        throwsA(equals(errorReservedName)),
      );
      expect(
        () => FdawgNamer.isValidAppName('PRN', platforms),
        throwsA(equals(errorReservedName)),
      );
    });

    test('Throws error for app name starting or ending with special characters', () {
      expect(
        () => FdawgNamer.isValidAppName('_AppName', platforms),
        throwsA(equals(errorStartsOrEndsWithSpecialChar)),
      );
      expect(
        () => FdawgNamer.isValidAppName('AppName-', platforms),
        throwsA(equals(errorStartsOrEndsWithSpecialChar)),
      );
    });

    test('Throws error for app name exceeding max length', () {
      final longName = 'A' * 256; // 256 characters
      expect(
        () => FdawgNamer.isValidAppName(longName, platforms),
        throwsA(equals(errorInvalidLength)),
      );
    });
  });

  group('validateAppLabel', () {
    final platforms = [
      PlatformOptions.android,
      PlatformOptions.ios,
      PlatformOptions.web,
    ];

    test('Valid app labels/titles', () {
      expect(() => FdawgNamer.isValidAppName('My App', platforms), returnsNormally);
      expect(() => FdawgNamer.isValidAppName('CoolApp', platforms), returnsNormally);
    });

    test('Throws error for empty app label', () {
      expect(
        () => FdawgNamer.isValidAppName('', platforms),
        throwsA(equals(errorNameEmpty)),
      );
    });

    test('Throws error for invalid characters in app label', () {
      expect(
        () => FdawgNamer.isValidAppName('App@Label', platforms),
        throwsA(equals(errorInvalidCharacters)),
      );
    });

    test('Throws error for app label exceeding max length', () {
      final longLabel = 'A' * 31; // 31 characters
      expect(
        () => FdawgNamer.isValidAppName(longLabel, platforms),
        throwsA(equals(errorInvalidLength)),
      );
    });
  });
}

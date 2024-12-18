import 'package:fdawg_namer/common/errors.dart';
import 'package:fdawg_namer/fdawg_namer.dart';
import 'package:flutter_test/flutter_test.dart';

void main() {
  group('FdawgNamer.isValidProjectName', () {
    test('Valid names do not throw errors', () {
      expect(() => FdawgNamer.isValidProjectName('valid_name'), returnsNormally);
      expect(() => FdawgNamer.isValidProjectName('validname123'), returnsNormally);
      expect(() => FdawgNamer.isValidProjectName('flutter_app'), returnsNormally);
    });

    test('Throws error for empty name', () {
      expect(
            () => FdawgNamer.isValidProjectName(''),
        throwsA(equals(errorNameEmpty)),
      );
    });

    test('Throws error for name starting or ending with underscore', () {
      expect(
            () => FdawgNamer.isValidProjectName('_name'),
        throwsA(equals(errorStartEndsWithUnderscore)),
      );
      expect(
            () => FdawgNamer.isValidProjectName('name_'),
        throwsA(equals(errorStartEndsWithUnderscore)),
      );
    });

    test('Throws error for invalid format', () {
      expect(
            () => FdawgNamer.isValidProjectName('InvalidName'),
        throwsA(equals(errorInvalidFormat)),
      );
      expect(
            () => FdawgNamer.isValidProjectName('name-with-dash'),
        throwsA(equals(errorInvalidFormat)),
      );
      expect(
            () => FdawgNamer.isValidProjectName('123name'),
        throwsA(equals(errorInvalidFormat)),
      );
      expect(
            () => FdawgNamer.isValidProjectName('name\$special'),
        throwsA(equals(errorInvalidFormat)),
      );
    });

    test('Throws error for reserved Dart keywords', () {
      expect(
            () => FdawgNamer.isValidProjectName('class'),
        throwsA(equals(errorReservedKeyword)),
      );
      expect(
            () => FdawgNamer.isValidProjectName('import'),
        throwsA(equals(errorReservedKeyword)),
      );
      expect(
            () => FdawgNamer.isValidProjectName('void'),
        throwsA(equals(errorReservedKeyword)),
      );
    });
  });
}
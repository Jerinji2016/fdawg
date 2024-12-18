import 'package:fdawg_namer/common/errors.dart';
import 'package:fdawg_namer/fdawg_namer.dart';
import 'package:flutter_test/flutter_test.dart';

void main() {
  group('FdawgNamer.isValidName', () {
    test('Valid names do not throw errors', () {
      expect(() => FdawgNamer.isValidName('valid_name'), returnsNormally);
      expect(() => FdawgNamer.isValidName('validname123'), returnsNormally);
      expect(() => FdawgNamer.isValidName('flutter_app'), returnsNormally);
    });

    test('Throws error for empty name', () {
      expect(
            () => FdawgNamer.isValidName(''),
        throwsA(equals(errorNameEmpty)),
      );
    });

    test('Throws error for name starting or ending with underscore', () {
      expect(
            () => FdawgNamer.isValidName('_name'),
        throwsA(equals(errorStartEndsWithUnderscore)),
      );
      expect(
            () => FdawgNamer.isValidName('name_'),
        throwsA(equals(errorStartEndsWithUnderscore)),
      );
    });

    test('Throws error for invalid format', () {
      expect(
            () => FdawgNamer.isValidName('InvalidName'),
        throwsA(equals(errorInvalidFormat)),
      );
      expect(
            () => FdawgNamer.isValidName('name-with-dash'),
        throwsA(equals(errorInvalidFormat)),
      );
      expect(
            () => FdawgNamer.isValidName('123name'),
        throwsA(equals(errorInvalidFormat)),
      );
      expect(
            () => FdawgNamer.isValidName('name\$special'),
        throwsA(equals(errorInvalidFormat)),
      );
    });

    test('Throws error for reserved Dart keywords', () {
      expect(
            () => FdawgNamer.isValidName('class'),
        throwsA(equals(errorReservedKeyword)),
      );
      expect(
            () => FdawgNamer.isValidName('import'),
        throwsA(equals(errorReservedKeyword)),
      );
      expect(
            () => FdawgNamer.isValidName('void'),
        throwsA(equals(errorReservedKeyword)),
      );
    });
  });
}

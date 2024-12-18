part of '../fdawg_namer.dart';

void _validateName(String name) {
  if (name.isEmpty) {
    throw errorNameEmpty;
  }

  if (name.startsWith('_') || name.endsWith('_')) {
    throw errorStartEndsWithUnderscore;
  }

  final regex = RegExp(r'^[a-z][a-z0-9_]*$');
  if (!regex.hasMatch(name)) {
    throw errorInvalidFormat;
  }

  if (_kDartReservedKeywords.contains(name)) {
    throw errorReservedKeyword;
  }
}

const _kDartReservedKeywords = {
  'abstract',
  'as',
  'assert',
  'async',
  'await',
  'break',
  'case',
  'catch',
  'class',
  'const',
  'continue',
  'covariant',
  'default',
  'deferred',
  'do',
  'dynamic',
  'else',
  'enum',
  'export',
  'extends',
  'extension',
  'external',
  'factory',
  'false',
  'final',
  'finally',
  'for',
  'Function',
  'get',
  'hide',
  'if',
  'implements',
  'import',
  'in',
  'interface',
  'is',
  'library',
  'mixin',
  'new',
  'null',
  'on',
  'operator',
  'part',
  'return',
  'set',
  'rethrow',
  'show',
  'static',
  'super',
  'switch',
  'sync',
  'this',
  'throw',
  'true',
  'try',
  'typedef',
  'var',
  'void',
  'while',
  'with',
  'yield',
};

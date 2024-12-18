part of '../fdawg_core.dart';

enum PlatformOptions {
  android('Android', 'android'),
  ios('iOS', 'ios'),
  web('Web', 'web'),
  linux('Linux', 'linux'),
  macos('MacOS', 'macos'),
  windows('Windows', 'windows');

  const PlatformOptions(this.label, this.value);

  final String label;
  final String value;
}

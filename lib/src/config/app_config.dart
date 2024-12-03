import 'package:flutter/material.dart';

class AppConfig extends ChangeNotifier {
  factory AppConfig() => _mInstance;

  AppConfig._();

  static final AppConfig _mInstance = AppConfig._();

  ThemeMode _themeMode = ThemeMode.system;

  ThemeMode get themeMode => _themeMode;

  void toggleTheme() {
    final currentThemeIndex = ThemeMode.values.indexOf(_themeMode);
    final nextIndex = (currentThemeIndex + 1) % ThemeMode.values.length;
    _themeMode = ThemeMode.values.elementAt(nextIndex);
    notifyListeners();
  }
}

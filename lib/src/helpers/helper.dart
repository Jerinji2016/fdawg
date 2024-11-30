import 'package:flutter/material.dart';

bool isDarkTheme(BuildContext context) {
  return Theme.of(context).brightness == Brightness.dark;
}

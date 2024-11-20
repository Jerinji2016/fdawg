import 'package:flutter/material.dart';

import '../config/route_config.dart';

bool get isDarkTheme {
  final context = globalNavigatorKey.currentState?.context;
  if(context == null) return false;
  return Theme.of(context).brightness == Brightness.dark;
}
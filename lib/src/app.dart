import 'package:flutter/material.dart';

import 'config/route_config.dart';
import 'config/theme_config.dart';

class FdawgApp extends StatelessWidget {
  const FdawgApp({super.key});

  @override
  Widget build(BuildContext context) {
    return MaterialApp.router(
      routerConfig: router,
      theme: lightTheme,
      darkTheme: darkTheme,
      debugShowCheckedModeBanner: false,
    );
  }
}

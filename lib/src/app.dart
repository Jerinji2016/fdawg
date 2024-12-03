import 'package:flutter/material.dart';
import 'package:provider/provider.dart';

import 'config/app_config.dart';
import 'config/route_config.dart';
import 'config/theme_config.dart';
import 'helpers/helper.dart';

const lightBgGradient = BoxDecoration(
  gradient: RadialGradient(
    colors: [
      Color(0xFFe8e8ea),
      Color(0xffced5da),
      Color(0xFFb2c1cd),
    ],
  ),
);
const darkBgGradient = BoxDecoration(
  gradient: RadialGradient(
    colors: [
      Color(0xff2c2c2c),
      Color(0xff1c1e1e),
      Color(0xff131315),
    ],
  ),
);

class FdawgApp extends StatelessWidget {
  const FdawgApp({super.key});

  @override
  Widget build(BuildContext context) {
    return ChangeNotifierProvider(
      create: (context) => AppConfig(),
      builder: (context, _) {
        final appConfig = Provider.of<AppConfig>(context);

        return MaterialApp.router(
          routerConfig: router,
          theme: lightTheme,
          darkTheme: darkTheme,
          themeMode: appConfig.themeMode,
          debugShowCheckedModeBanner: false,
          builder: (context, child) {
            return Stack(
              children: [
                Container(
                  decoration: isDarkTheme(context) ? darkBgGradient : lightBgGradient,
                  child: child,
                ),
                Positioned(
                  right: 16,
                  top: 16,
                  child: _buildThemeToggleButton(context),
                ),
              ],
            );
          },
        );
      },
    );
  }

  Widget _buildThemeToggleButton(BuildContext context) {
    final appConfig = Provider.of<AppConfig>(context);
    IconData icon;
    String text;

    switch (appConfig.themeMode) {
      case ThemeMode.system:
        icon = Icons.auto_awesome;
        text = 'System';
      case ThemeMode.light:
        icon = Icons.light_mode_outlined;
        text = 'Light';
      case ThemeMode.dark:
        icon = Icons.dark_mode_outlined;
        text = 'Dark';
    }

    return IconButton(
      onPressed: appConfig.toggleTheme,
      icon: Row(
        mainAxisSize: MainAxisSize.min,
        children: [
          Icon(
            icon,size: 18,

          ),
          const SizedBox(width: 8),
          Text(
            text,
            style: const TextStyle(fontWeight: FontWeight.bold),
          ),
        ],
      ),
    );
  }
}

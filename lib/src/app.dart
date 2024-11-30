import 'package:flutter/material.dart';

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
    return MaterialApp.router(
      routerConfig: router,
      theme: lightTheme,
      darkTheme: darkTheme,
      themeMode: ThemeMode.light,
      debugShowCheckedModeBanner: false,
      builder: (context, child) {
        return Container(
          decoration: isDarkTheme(context) ? darkBgGradient : lightBgGradient,
          child: child,
        );
      },
    );
  }
}

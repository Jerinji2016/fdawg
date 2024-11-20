import 'package:flutter/cupertino.dart';
import 'package:go_router/go_router.dart';

import '../ui/intro/splash_page.dart';
import '../ui/intro/welcome_page.dart';

final globalNavigatorKey = GlobalKey<NavigatorState>();

final router = GoRouter(
  navigatorKey: globalNavigatorKey,
  routes: [
    GoRoute(
        name: 'intro',
        path: '/',
        builder: (context, state) => const SplashPage(),
    ),
    GoRoute(
        name: 'welcome',
        path: '/welcome',
        builder: (context, state) => const WelcomePage(),
    ),
  ],
);

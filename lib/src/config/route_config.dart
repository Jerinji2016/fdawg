import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';

import '../ui/create_project/create_project_page.dart';
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
    GoRoute(
      name: 'create-project',
      path: '/create',
      builder: (context, state) => const CreateProjectPage(),
    ),
  ],
);

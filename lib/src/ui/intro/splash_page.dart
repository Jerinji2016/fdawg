import 'package:flutter/material.dart';
import 'package:flutter/scheduler.dart';
import 'package:go_router/go_router.dart';

import '../../widgets/logo.dart';

class SplashPage extends StatelessWidget {
  const SplashPage({super.key});

  @override
  Widget build(BuildContext context) {
    SchedulerBinding.instance.addPostFrameCallback((duration) {
      Future<void>.delayed(const Duration(seconds: 2)).then((value) {
        if (context.mounted) {
          context.go('/welcome');
        }
      });
    });

    return const Scaffold(
      body: Center(
        child: Logo(),
      ),
    );
  }
}

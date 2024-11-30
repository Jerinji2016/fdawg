import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';

import '../../widgets/logo.dart';
import '../../widgets/primary_card.dart';

class WelcomePage extends StatelessWidget {
  const WelcomePage({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: Center(
        child: ConstrainedBox(
          constraints: BoxConstraints(
            maxWidth: MediaQuery.of(context).size.width * 0.4,
          ),
          child: Column(
            mainAxisSize: MainAxisSize.min,
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              const Logo(),
              const SizedBox(height: 32),
              WelcomeActionCard(
                icon: Icons.add,
                title: 'New Project',
                description: 'Setup a new Flutter Project',
                onTap: () => context.push('/create'),
              ),
              const SizedBox(height: 16),
              WelcomeActionCard(
                icon: Icons.get_app_outlined,
                title: 'Import Project',
                description: 'Import an existing Flutter Project',
                onTap: () {},
              ),
            ],
          ),
        ),
      ),
    );
  }
}

class WelcomeActionCard extends StatelessWidget {
  const WelcomeActionCard({
    required this.icon,
    required this.title,
    required this.description,
    required this.onTap,
    super.key,
  });

  final IconData icon;
  final String title;
  final String description;
  final VoidCallback onTap;

  @override
  Widget build(BuildContext context) {
    return PrimaryCard(
      onTap: onTap,
      child: Padding(
        padding: const EdgeInsets.all(16),
        child: Row(
          children: [
            SizedBox.square(
              dimension: 48,
              child: Icon(icon),
            ),
            const SizedBox(width: 16),
            Expanded(
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(
                    title,
                    style: const TextStyle(
                      fontSize: 16,
                      fontWeight: FontWeight.w500,
                    ),
                  ),
                  Text(
                    description,
                    style: TextStyle(
                      fontSize: 14,
                      color: Colors.grey.shade500,
                    ),
                  ),
                ],
              ),
            ),
          ],
        ),
      ),
    );
  }
}

import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';

import '../../widgets/logo.dart';
import '../../widgets/primary_button.dart';
import 'steps/name_your_project.dart';

class CreateProjectPage extends StatefulWidget {
  const CreateProjectPage({super.key});

  @override
  State<CreateProjectPage> createState() => _CreateProjectPageState();
}

class _CreateProjectPageState extends State<CreateProjectPage> {
  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: Center(
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            const Logo(),
            const SizedBox(height: 48),
            ConstrainedBox(
              constraints: BoxConstraints(
                maxWidth: MediaQuery.of(context).size.width * 0.3,
              ),
              child: const NameYourProject(),
            ),
            const SizedBox(height: 120),
            PrimaryButton(
              text: 'Home',
              color: Colors.transparent,
              prefix: const Icon(Icons.arrow_back),
              side: BorderSide(
                color: Theme.of(context).dividerColor,
                width: 1.5,
              ),
              textColor: Theme.of(context).textTheme.bodySmall?.color,
              onTap: () => context.pop(),
            )
          ],
        ),
      ),
    );
  }
}

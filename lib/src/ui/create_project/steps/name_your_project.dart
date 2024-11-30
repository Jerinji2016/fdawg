import 'package:flutter/material.dart';

import '../../../widgets/primary_button.dart';

class NameYourProject extends StatefulWidget {
  const NameYourProject({super.key});

  @override
  State<NameYourProject> createState() => _NameYourProjectState();
}

class _NameYourProjectState extends State<NameYourProject> {
  @override
  Widget build(BuildContext context) {
    return Column(
      mainAxisSize: MainAxisSize.min,
      children: [
        const Text(
          'Name your project',
          style: TextStyle(
            fontSize: 24,
            fontWeight: FontWeight.bold,
          ),
        ),
        const SizedBox(height: 24),
        TextField(
          controller: TextEditingController(),
          decoration: InputDecoration(
            label: const Text('Project Name'),
            floatingLabelStyle: TextStyle(
              color: Theme.of(context).colorScheme.primary,
            ),
          ),
        ),
        const SizedBox(height: 24),
        Align(
          alignment: Alignment.centerRight,
          child: PrimaryButton(
            text: 'Next',
            onTap: () {},
            suffix: const Icon(
              Icons.arrow_forward,
              color: Colors.white,
            ),
          ),
        ),
      ],
    );
  }
}

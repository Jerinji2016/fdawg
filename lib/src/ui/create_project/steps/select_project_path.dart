import 'package:flutter/material.dart';
import 'package:provider/provider.dart';

import '../../../widgets/primary_button.dart';
import '../create_project.vm.dart';

class SelectProjectPath extends StatelessWidget {
  const SelectProjectPath({super.key});

  @override
  Widget build(BuildContext context) {
    final viewModel = Provider.of<CreateProjectViewModel>(context);
    final directoryPath = viewModel.directoryPath;

    return Column(
      mainAxisSize: MainAxisSize.min,
      children: [
        const Text(
          'Select Project Path',
          style: TextStyle(
            fontSize: 24,
            fontWeight: FontWeight.bold,
          ),
        ),
        const SizedBox(height: 24),
        if (directoryPath != null)
          Container(
            decoration: BoxDecoration(
              color: Colors.white12,
              borderRadius: BorderRadius.circular(12),
            ),
            padding: const EdgeInsets.symmetric(vertical: 8, horizontal: 12),
            margin: const EdgeInsets.only(bottom: 16),
            child: Text(
              directoryPath,
              style: const TextStyle(
                fontWeight: FontWeight.w500,
              ),
            ),
          ),
        PrimaryButton(
          text: directoryPath == null ? 'Choose' : 'Change',
          onTap: viewModel.pickDirectory,
        ),
        if (directoryPath != null)
          Align(
            alignment: Alignment.centerRight,
            child: PrimaryButton(
              text: 'Next',
              onTap: viewModel.onNextTapped,
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

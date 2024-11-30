import 'package:flutter/material.dart';
import 'package:provider/provider.dart';

import '../../../widgets/primary_button.dart';
import '../create_project.vm.dart';

class SetAppIdentity extends StatelessWidget {
  const SetAppIdentity({super.key});

  @override
  Widget build(BuildContext context) {
    final viewModel = Provider.of<CreateProjectViewModel>(context);

    return Column(
      mainAxisSize: MainAxisSize.min,
      children: [
        const Text(
          'Set App Identity',
          style: TextStyle(
            fontSize: 24,
            fontWeight: FontWeight.bold,
          ),
        ),
        const SizedBox(height: 24),
        _buildImagePicker(),
        const SizedBox(height: 16),
        TextField(
          controller: viewModel.projectNameController,
          decoration: InputDecoration(
            label: const Text('App Name'),
            hintText: 'Cool App',
            floatingLabelStyle: TextStyle(
              color: Theme.of(context).colorScheme.primary,
            ),
          ),
        ),
        const SizedBox(height: 24),
        Row(
          children: [
            PrimaryButton(
              text: 'Back',
              onTap: viewModel.onPreviousTapped,
              prefix: const Icon(
                Icons.arrow_back,
                color: Colors.white,
              ),
            ),
            const Spacer(),
            PrimaryButton(
              text: 'Next',
              onTap: viewModel.onNextTapped,
              suffix: const Icon(
                Icons.arrow_forward,
                color: Colors.white,
              ),
            ),
          ],
        ),
      ],
    );
  }

  Widget _buildImagePicker() {
    return Container(
      height: 120,
      width: 120,
      decoration: BoxDecoration(
        borderRadius: BorderRadius.circular(24),
        border: Border.all(
          style: BorderStyle.solid
        ),
      ),
    );
  }
}

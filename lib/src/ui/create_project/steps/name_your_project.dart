import 'package:flutter/material.dart';
import 'package:provider/provider.dart';

import '../../../widgets/primary_button.dart';
import '../create_project.vm.dart';

class NameYourProject extends StatefulWidget {
  const NameYourProject({super.key});

  @override
  State<NameYourProject> createState() => _NameYourProjectState();
}

class _NameYourProjectState extends State<NameYourProject> {
  @override
  Widget build(BuildContext context) {
    final viewModel = Provider.of<CreateProjectViewModel>(context);

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
          controller: viewModel.projectNameController,
          decoration: InputDecoration(
            label: const Text('Project Name'),
            hintText: 'My Awesome Project',
            floatingLabelStyle: TextStyle(
              color: Theme.of(context).colorScheme.primary,
            ),
          ),
        ),
        const SizedBox(height: 16),
        TextField(
          controller: TextEditingController(),
          maxLines: 2,
          decoration: InputDecoration(
            label: const Text('Description'),
            hintText: 'My awesome project is going to ...',
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

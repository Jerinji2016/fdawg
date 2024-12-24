import 'package:fdawg_namer/fdawg_namer.dart';
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
  String? _projectNameError;

  void _onProjectNameChanged(String text) {
    if (_projectNameError != null) {
      setState(() => _projectNameError = null);
    }
  }

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
          onChanged: _onProjectNameChanged,
          decoration: InputDecoration(
            errorText: _projectNameError,
            label: const Text('Project Name'),
            hintText: 'My Awesome Project',
            floatingLabelStyle: TextStyle(
              color: Theme.of(context).colorScheme.primary,
            ),
          ),
        ),
        const SizedBox(height: 16),
        TextField(
          controller: viewModel.projectDescriptionController,
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
            onTap: _onNextTapped,
            suffix: const Icon(
              Icons.arrow_forward,
              color: Colors.white,
            ),
          ),
        ),
      ],
    );
  }

  void _onNextTapped() {
    final viewModel = Provider.of<CreateProjectViewModel>(context, listen: false);
    final name = viewModel.projectNameController.text;
    try {
      FdawgNamer.isValidProjectName(name);
    } catch (e) {
      debugPrint('_NameYourProjectState._onNextTapped: ❌ERROR: $e');
      return setState(
        () => _projectNameError = e.toString(),
      );
    }

    viewModel.onNextTapped();
  }
}

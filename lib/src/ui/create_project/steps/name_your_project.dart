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
    final screenHeight = MediaQuery.of(context).size.height;
    final isCompactHeight = screenHeight < 600;

    return SingleChildScrollView(
      child: Column(
        mainAxisSize: MainAxisSize.min,
        children: [
          Text(
            'Name your project',
            style: TextStyle(
              fontSize: isCompactHeight ? 18 : 24,
              fontWeight: FontWeight.bold,
            ),
          ),
          SizedBox(height: isCompactHeight ? 12 : 24),
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
          SizedBox(height: isCompactHeight ? 8 : 16),
          TextField(
            controller: viewModel.projectDescriptionController,
            maxLines: isCompactHeight ? 1 : 2,
            decoration: InputDecoration(
              label: const Text('Description'),
              hintText: 'My awesome project is going to ...',
              floatingLabelStyle: TextStyle(
                color: Theme.of(context).colorScheme.primary,
              ),
            ),
          ),
          SizedBox(height: isCompactHeight ? 12 : 24),
          Row(
            children: [
              PrimaryButton(
                text: isCompactHeight ? '' : 'Back',
                onTap: viewModel.onPreviousTapped,
                prefix: Icon(
                  Icons.arrow_back,
                  color: Colors.white,
                  size: isCompactHeight ? 18 : 24,
                ),
                padding: isCompactHeight
                    ? const EdgeInsets.all(10)
                    : const EdgeInsets.symmetric(horizontal: 16, vertical: 12),
              ),
              const Spacer(),
              PrimaryButton(
                text: isCompactHeight ? '' : 'Next',
                onTap: _onNextTapped,
                suffix: Icon(
                  Icons.arrow_forward,
                  color: Colors.white,
                  size: isCompactHeight ? 18 : 24,
                ),
                padding: isCompactHeight
                    ? const EdgeInsets.all(10)
                    : const EdgeInsets.symmetric(horizontal: 16, vertical: 12),
              ),
            ],
          ),
        ],
      ),
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

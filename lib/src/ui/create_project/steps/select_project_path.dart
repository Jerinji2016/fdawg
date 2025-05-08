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
    final screenHeight = MediaQuery.of(context).size.height;
    final isCompactHeight = screenHeight < 600;

    return SingleChildScrollView(
      child: Column(
        mainAxisSize: MainAxisSize.min,
        children: [
          Text(
            'Select Project Path',
            style: TextStyle(
              fontSize: isCompactHeight ? 18 : 24,
              fontWeight: FontWeight.bold,
            ),
          ),
          SizedBox(height: isCompactHeight ? 12 : 24),
          if (directoryPath != null)
            Container(
              decoration: BoxDecoration(
                color: Colors.white12,
                borderRadius: BorderRadius.circular(12),
              ),
              padding: EdgeInsets.symmetric(
                vertical: isCompactHeight ? 4 : 8, 
                horizontal: 12
              ),
              margin: EdgeInsets.only(bottom: isCompactHeight ? 8 : 12),
              child: Row(
                children: [
                  Expanded(
                    child: Text(
                      directoryPath,
                      style: const TextStyle(
                        fontWeight: FontWeight.w500,
                      ),
                      overflow: TextOverflow.ellipsis,
                      maxLines: isCompactHeight ? 1 : 2,
                    ),
                  ),
                  IconButton(
                    icon: const Icon(Icons.edit, size: 16),
                    onPressed: viewModel.pickDirectory,
                    tooltip: 'Change directory',
                    padding: EdgeInsets.zero,
                    constraints: const BoxConstraints(
                      minWidth: 32,
                      minHeight: 32,
                    ),
                  ),
                ],
              ),
            ),
          if (directoryPath == null)
            PrimaryButton(
              text: isCompactHeight ? '' : 'Choose Directory',
              onTap: viewModel.pickDirectory,
              prefix: const Icon(
                Icons.folder_open,
                color: Colors.white,
              ),
              padding: isCompactHeight
                  ? const EdgeInsets.all(10)
                  : const EdgeInsets.symmetric(horizontal: 16, vertical: 12),
            ),
          SizedBox(height: isCompactHeight ? 8 : 16),
          if (directoryPath != null)
            Row(
              mainAxisAlignment: MainAxisAlignment.end,
              children: [
                PrimaryButton(
                  text: isCompactHeight ? '' : 'Next',
                  onTap: viewModel.onNextTapped,
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
}

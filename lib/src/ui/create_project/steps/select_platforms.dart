import 'package:flutter/material.dart';
import 'package:provider/provider.dart';

import '../../../widgets/primary_button.dart';
import '../create_project.vm.dart';

enum PlatformOptions {
  android('Android'),
  ios('iOS'),
  linux('Linux'),
  macos('MacOS'),
  web('Web'),
  windows('Windows');

  const PlatformOptions(this.label);

  final String label;
}

class SelectPlatforms extends StatelessWidget {
  const SelectPlatforms({super.key});

  @override
  Widget build(BuildContext context) {
    final viewModel = Provider.of<CreateProjectViewModel>(context);

    return Column(
      mainAxisSize: MainAxisSize.min,
      children: [
        const Text(
          'Select Platforms',
          style: TextStyle(
            fontSize: 24,
            fontWeight: FontWeight.bold,
          ),
        ),
        const SizedBox(height: 24),
        ConstrainedBox(
          constraints: BoxConstraints(
            maxHeight: MediaQuery.of(context).size.height * 0.3,
            maxWidth: MediaQuery.of(context).size.width * 0.4,
          ),
          child: SingleChildScrollView(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              mainAxisSize: MainAxisSize.min,
              children: List.generate(
                PlatformOptions.values.length,
                (index) => _buildPlatformOptions(context, index),
              ),
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
              text: 'Finish',
              onTap: viewModel.onFinishTapped,
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

  Widget _buildPlatformOptions(BuildContext context, int index) {
    final viewModel = Provider.of<CreateProjectViewModel>(context);
    final option = PlatformOptions.values.elementAt(index);
    final isSelected = viewModel.platformOptions[option] ?? false;

    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 4),
      child: Row(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Checkbox(
            value: isSelected,
            onChanged: (value) => viewModel.onPlatformSelected(option, isSelected: value),
          ),
          const SizedBox(width: 12),
          Expanded(
            child: Padding(
              padding: const EdgeInsets.symmetric(vertical: 6),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                mainAxisSize: MainAxisSize.min,
                children: [
                  Text(
                    option.label,
                    style: TextStyle(
                      fontWeight: isSelected ? FontWeight.bold : FontWeight.normal,
                    ),
                  ),
                  AnimatedSize(
                    duration: const Duration(milliseconds: 100),
                    child: isSelected
                        ? const Padding(
                            padding: EdgeInsets.only(top: 8),
                            child: TextField(),
                          )
                        : const SizedBox.shrink(),
                  ),
                ],
              ),
            ),
          ),
        ],
      ),
    );
  }
}

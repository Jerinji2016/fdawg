import 'package:fdawg_core/fdawg_core.dart';
import 'package:flutter/material.dart';
import 'package:provider/provider.dart';

import '../../../widgets/primary_button.dart';
import '../create_project.vm.dart';

extension PlatformOptionsExtension on PlatformOptions {
  String get placeholder {
    switch (this) {
      case PlatformOptions.android:
        return 'Android Bundle ID';
      case PlatformOptions.ios:
        return 'iOS Bundle ID';
      case PlatformOptions.web:
        return '';
      case PlatformOptions.linux:
        return 'Linux Package Name';
      case PlatformOptions.macos:
        return 'MacOS Bundle ID';
      case PlatformOptions.windows:
        return 'Windows Package ID';
    }
  }

  bool get hasBundleId => this != PlatformOptions.web;
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
            // maxWidth: MediaQuery.of(context).size.width * 0.4,
          ),
          child: SingleChildScrollView(
            child: Wrap(
              // crossAxisAlignment: CrossAxisAlignment.start,
              // mainAxisSize: MainAxisSize.min,
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
      padding: const EdgeInsets.symmetric(vertical: 4, horizontal: 8),
      child: Row(
        mainAxisSize: MainAxisSize.min,
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Checkbox(
            value: isSelected,
            onChanged: (value) => viewModel.onPlatformSelected(option, isSelected: value),
          ),
          const SizedBox(width: 2),
          Container(
            constraints: BoxConstraints(
              // maxHeight: MediaQuery.of(context).size.height * 0.3,
              maxWidth: MediaQuery.of(context).size.width * 0.28,
            ),
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
                  duration: const Duration(milliseconds: 300),
                  curve: Curves.fastLinearToSlowEaseIn,
                  child: isSelected && option.hasBundleId
                      ? Padding(
                          padding: const EdgeInsets.only(top: 8),
                          child: TextField(
                            decoration: InputDecoration(
                              hintText: option.placeholder,
                            ),
                          ),
                        )
                      : const SizedBox.shrink(),
                ),
              ],
            ),
          ),
        ],
      ),
    );
  }
}

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

  IconData get icon {
    switch (this) {
      case PlatformOptions.android:
        return Icons.android;
      case PlatformOptions.ios:
        return Icons.phone_iphone;
      case PlatformOptions.web:
        return Icons.web;
      case PlatformOptions.linux:
        return Icons.computer;
      case PlatformOptions.macos:
        return Icons.laptop_mac;
      case PlatformOptions.windows:
        return Icons.window;
    }
  }
}

class SelectPlatforms extends StatelessWidget {
  const SelectPlatforms({super.key});

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
            'Select Platforms',
            style: TextStyle(
              fontSize: isCompactHeight ? 18 : 24,
              fontWeight: FontWeight.bold,
            ),
          ),
          SizedBox(height: isCompactHeight ? 12 : 24),
          _buildPlatformIconGrid(context, isCompactHeight),
          SizedBox(height: isCompactHeight ? 12 : 24),
          _buildBundleIdSection(context, isCompactHeight),
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
                text: isCompactHeight ? '' : 'Finish',
                onTap: viewModel.onFinishTapped,
                suffix: Icon(
                  Icons.check,
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

  Widget _buildPlatformIconGrid(BuildContext context, bool isCompact) {
    final viewModel = Provider.of<CreateProjectViewModel>(context);

    return Wrap(
      spacing: isCompact ? 12 : 16,
      runSpacing: isCompact ? 12 : 16,
      alignment: WrapAlignment.center,
      children: PlatformOptions.values.map((platform) {
        final isSelected = viewModel.platformOptions[platform] ?? false;

        return Tooltip(
          message: platform.label,
          waitDuration: const Duration(milliseconds: 500),
          child: InkWell(
            onTap: () => viewModel.onPlatformSelected(
              platform,
              isSelected: !isSelected,
            ),
            borderRadius: BorderRadius.circular(isCompact ? 20 : 24),
            child: Container(
              width: isCompact ? 40 : 48,
              height: isCompact ? 40 : 48,
              decoration: BoxDecoration(
                color: isSelected
                    ? Theme.of(context).colorScheme.primary
                    : Colors.grey.withValues(alpha: .1),
                borderRadius: BorderRadius.circular(isCompact ? 20 : 24),
                border: Border.all(
                  color: isSelected
                      ? Theme.of(context).colorScheme.primary
                      : Colors.grey.withValues(alpha: .3),
                ),
              ),
              child: Stack(
                children: [
                  Center(
                    child: Icon(
                      platform.icon,
                      size: isCompact ? 22 : 26,
                      color: isSelected ? Colors.white : Colors.grey,
                    ),
                  ),
                  if (isSelected)
                    Positioned(
                      right: 0,
                      bottom: 0,
                      child: Container(
                        width: isCompact ? 14 : 16,
                        height: isCompact ? 14 : 16,
                        decoration: const BoxDecoration(
                          color: Colors.white,
                          shape: BoxShape.circle,
                        ),
                        child: Icon(
                          Icons.check,
                          size: isCompact ? 10 : 12,
                          color: Theme.of(context).colorScheme.primary,
                        ),
                      ),
                    ),
                ],
              ),
            ),
          ),
        );
      }).toList(),
    );
  }

  Widget _buildBundleIdSection(BuildContext context, bool isCompact) {
    final viewModel = Provider.of<CreateProjectViewModel>(context);
    final selectedPlatforms = PlatformOptions.values
        .where(
          (platform) =>
              (viewModel.platformOptions[platform] ?? false) &&
              platform.hasBundleId,
        )
        .toList();

    if (selectedPlatforms.isEmpty) return const SizedBox.shrink();

    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Padding(
          padding: const EdgeInsets.only(bottom: 8),
          child: Text(
            'Bundle IDs',
            style: TextStyle(
              fontSize: isCompact ? 14 : 16,
              fontWeight: FontWeight.w500,
            ),
          ),
        ),
        ...selectedPlatforms.map((platform) {
          TextEditingController? controller;

          switch (platform) {
            case PlatformOptions.android:
              controller = viewModel.androidBundleIdController;
            case PlatformOptions.ios:
              controller = viewModel.iosBundleIdController;
            case PlatformOptions.linux:
              controller = viewModel.linuxBundleIdController;
            case PlatformOptions.macos:
              controller = viewModel.macBundleIdController;
            case PlatformOptions.windows:
              controller = viewModel.windowsBundleIdController;
            case PlatformOptions.web:
              controller = null;
          }

          return Padding(
            padding: const EdgeInsets.only(bottom: 8),
            child: TextField(
              controller: controller,
              style: TextStyle(fontSize: isCompact ? 13 : 14),
              decoration: InputDecoration(
                prefixIcon: Icon(platform.icon, size: isCompact ? 18 : 20),
                hintText: platform.placeholder,
                contentPadding: isCompact
                    ? const EdgeInsets.symmetric(horizontal: 8, vertical: 8)
                    : null,
                isDense: isCompact,
              ),
            ),
          );
        }),
      ],
    );
  }
}

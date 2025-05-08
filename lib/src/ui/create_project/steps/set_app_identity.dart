import 'package:flutter/material.dart';
import 'package:provider/provider.dart';

import '../../../widgets/primary_button.dart';
import '../create_project.vm.dart';
import '../widgets/app_icon_avatar.dart';

class SetAppIdentity extends StatelessWidget {
  const SetAppIdentity({super.key});

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
            'Set App Identity',
            style: TextStyle(
              fontSize: isCompactHeight ? 18 : 24,
              fontWeight: FontWeight.bold,
            ),
          ),
          SizedBox(height: isCompactHeight ? 16 : 32),
          _buildAppIdentityContainer(context, viewModel, isCompactHeight),
          SizedBox(height: isCompactHeight ? 16 : 24),
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

  Widget _buildAppIdentityContainer(
    BuildContext context,
    CreateProjectViewModel viewModel,
    bool isCompact,
  ) {
    return Container(
      decoration: BoxDecoration(
        color: Theme.of(context).colorScheme.surface.withValues(alpha: 0.7),
        borderRadius: BorderRadius.circular(16),
        border: Border.all(
          color: Theme.of(context).colorScheme.primary.withValues(alpha: 0.2),
        ),
        boxShadow: [
          BoxShadow(
            color: Theme.of(context).shadowColor.withValues(alpha: 0.1),
            blurRadius: 10,
          ),
        ],
        // Glassmorphic effect
        gradient: LinearGradient(
          begin: Alignment.topLeft,
          end: Alignment.bottomRight,
          colors: [
            Theme.of(context).colorScheme.surface.withValues(alpha: 0.8),
            Theme.of(context).colorScheme.surface.withValues(alpha: 0.6),
          ],
        ),
      ),
      padding: EdgeInsets.all(isCompact ? 16 : 24),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Expanded(
                flex: 2,
                child: _buildAppIconSection(context, viewModel, isCompact),
              ),
              SizedBox(width: isCompact ? 12 : 20),
              Expanded(
                flex: 3,
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(
                      'App Name',
                      style: TextStyle(
                        fontSize: isCompact ? 14 : 16,
                        fontWeight: FontWeight.w500,
                        color: Theme.of(context)
                            .colorScheme
                            .onSurface
                            .withValues(alpha: 0.8),
                      ),
                    ),
                    SizedBox(height: isCompact ? 8 : 12),
                    TextField(
                      controller: viewModel.appNameController,
                      decoration: InputDecoration(
                        hintText: 'Cool App',
                        prefixIcon: const Icon(Icons.app_shortcut),
                        border: OutlineInputBorder(
                          borderRadius: BorderRadius.circular(8),
                        ),
                        contentPadding: isCompact
                            ? const EdgeInsets.symmetric(
                                horizontal: 12,
                                vertical: 8,
                              )
                            : null,
                      ),
                    ),
                    SizedBox(height: isCompact ? 6 : 8),
                    Text(
                      'This name will be displayed on the device home screen',
                      style: TextStyle(
                        fontSize: isCompact ? 11 : 12,
                        color: Theme.of(context).hintColor,
                      ),
                    ),
                    // Add spacer to match height with icon section
                    SizedBox(height: isCompact ? 20 : 30),
                  ],
                ),
              ),
            ],
          ),
          Padding(
            padding: EdgeInsets.only(top: isCompact ? 12 : 16),
            child: Text(
              "Your app identity defines how users will recognize your app on their devices. Choose a distinctive icon and a memorable name that reflects your app's purpose.",
              style: TextStyle(
                fontSize: isCompact ? 12 : 13,
                color: Theme.of(context)
                    .colorScheme
                    .onSurface
                    .withValues(alpha: .6),
                fontStyle: FontStyle.italic,
              ),
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildAppIconSection(
    BuildContext context,
    CreateProjectViewModel viewModel,
    bool isCompact,
  ) {
    final iconSize = isCompact ? 80.0 : 100.0;

    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(
          'App Icon',
          style: TextStyle(
            fontSize: isCompact ? 14 : 16,
            fontWeight: FontWeight.w500,
            color:
                Theme.of(context).colorScheme.onSurface.withValues(alpha: 0.8),
          ),
        ),
        SizedBox(height: isCompact ? 8 : 12),
        SizedBox(
          height: iconSize,
          width: iconSize,
          child: AppIconAvatar(
            initialIconAsBytes: viewModel.appIconAsBytes,
            onRemove: viewModel.clearSelectedAppIcon,
            onImagePicked: viewModel.setAppIcon,
          ),
        ),
        if (viewModel.appIconAsBytes != null)
          Padding(
            padding: const EdgeInsets.only(top: 6),
            child: Text(
              'Icon selected',
              style: TextStyle(
                fontSize: isCompact ? 11 : 12,
                color: Theme.of(context).colorScheme.primary,
                fontWeight: FontWeight.w500,
              ),
            ),
          )
        else
          Padding(
            padding: const EdgeInsets.only(top: 6),
            child: Text(
              'Tap to select',
              style: TextStyle(
                fontSize: isCompact ? 11 : 12,
                color: Theme.of(context).hintColor,
                fontStyle: FontStyle.italic,
              ),
            ),
          ),
      ],
    );
  }
}

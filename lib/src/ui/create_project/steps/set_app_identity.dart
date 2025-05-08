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
          SizedBox(height: isCompactHeight ? 12 : 24),
          SizedBox(
            height: isCompactHeight ? 80 : 100,
            child: AppIconAvatar(
              initialIconAsBytes: viewModel.appIconAsBytes,
              onRemove: viewModel.clearSelectedAppIcon,
              onImagePicked: viewModel.setAppIcon,
            ),
          ),
          SizedBox(height: isCompactHeight ? 8 : 16),
          TextField(
            controller: viewModel.appNameController,
            decoration: InputDecoration(
              label: const Text('App Name'),
              hintText: 'Cool App',
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

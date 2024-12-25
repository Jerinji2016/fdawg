import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:provider/provider.dart';

import '../../widgets/fade_transition_switcher.dart';
import '../../widgets/logo.dart';
import '../../widgets/page_index_indicator.dart';
import '../../widgets/primary_button.dart';
import 'create_project.vm.dart';

class CreateProjectPage extends StatefulWidget {
  const CreateProjectPage({super.key});

  @override
  State<CreateProjectPage> createState() => _CreateProjectPageState();
}

class _CreateProjectPageState extends State<CreateProjectPage> {
  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: Center(
        child: ChangeNotifierProvider(
          create: (context) => CreateProjectViewModel(),
          builder: (context, child) {
            final viewModel = Provider.of<CreateProjectViewModel>(context, listen: false);

            return Column(
              mainAxisSize: MainAxisSize.min,
              children: [
                const Logo(),
                const SizedBox(height: 48),
                ConstrainedBox(
                  constraints: BoxConstraints(
                    maxWidth: MediaQuery.of(context).size.width * 0.35,
                  ),
                  child: _buildPages(context),
                ),
                const SizedBox(height: 24),
                PageIndexIndicator(
                  controller: viewModel.pageIndexController,
                ),
                const SizedBox(height: 80),
                PrimaryButton(
                  text: 'Home',
                  color: Colors.transparent,
                  prefix: const Icon(Icons.arrow_back),
                  side: BorderSide(
                    color: Theme.of(context).dividerColor,
                    width: 1.5,
                  ),
                  textColor: Theme.of(context).textTheme.bodySmall?.color,
                  onTap: () => context.pop(),
                ),
              ],
            );
          },
        ),
      ),
    );
  }

  Widget _buildPages(BuildContext context) {
    final viewModel = Provider.of<CreateProjectViewModel>(context);

    return AnimatedSize(
      duration: const Duration(milliseconds: 100),
      curve: Curves.fastLinearToSlowEaseIn,
      child: FadeTransitionSwitcher(
        item: viewModel.pages.elementAt(viewModel.currentPage),
      ),
    );
  }
}

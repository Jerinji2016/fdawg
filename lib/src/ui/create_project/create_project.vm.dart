import 'package:flutter/cupertino.dart';

import '../../widgets/page_index_indicator.dart';
import 'steps/select_platforms.dart';

class CreateProjectViewModel extends ChangeNotifier {
  final pageIndexController = PageIndexController(total: 3);

  final projectNameController = TextEditingController();
  final projectDescriptionController = TextEditingController();

  final appNameController = TextEditingController();

  final platformOptions = Map<PlatformOptions, bool>.fromEntries(
    PlatformOptions.values.map(
      (e) => MapEntry(e, false),
    ),
  );

  int _currentPageIndex = 0;

  int get currentPage => _currentPageIndex;

  void onNextTapped() {
    _currentPageIndex++;
    pageIndexController.nextIndex();
    notifyListeners();
  }

  void onPreviousTapped() {
    _currentPageIndex--;
    pageIndexController.previousIndex();
    notifyListeners();
  }

  void onPlatformSelected(PlatformOptions option, {bool? isSelected}) {
    if (isSelected == null) return;
    platformOptions[option] = isSelected;
    notifyListeners();
  }

  void onFinishTapped() {
    debugPrint('CreateProjectViewModel.onFinishTapped: 🐞');
  }

  @override
  void dispose() {
    pageIndexController.dispose();

    projectNameController.dispose();
    projectDescriptionController.dispose();
    appNameController.dispose();

    super.dispose();
  }
}

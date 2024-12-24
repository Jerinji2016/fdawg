import 'package:fdawg_core/fdawg_core.dart';
import 'package:file_picker/file_picker.dart';
import 'package:flutter/cupertino.dart';

import '../../widgets/page_index_indicator.dart';

class CreateProjectViewModel extends ChangeNotifier {
  final pageIndexController = PageIndexController(total: 3);

  final projectNameController = TextEditingController();
  final projectDescriptionController = TextEditingController();

  final appNameController = TextEditingController();

  String? _directoryPath;

  String? get directoryPath => _directoryPath;

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

  Future<void> pickDirectory() async {
    _directoryPath = await FilePicker.platform.getDirectoryPath();
    if (_directoryPath == null) return;

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

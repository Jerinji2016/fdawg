import 'package:fdawg_core/fdawg_core.dart';
import 'package:file_picker/file_picker.dart';
import 'package:flutter/cupertino.dart';
import 'package:flutter/foundation.dart';

import '../../widgets/page_index_indicator.dart';

class CreateProjectViewModel extends ChangeNotifier {
  CreateProjectViewModel() {
    if (kDebugMode) {
      projectNameController.text = 'test_project';
      projectDescriptionController.text = 'This is a project description';
      appNameController.text = 'Test Project 1';
      _directoryPath = '/Volumes/Macintosh HD/Users/jerin/Documents/Projects/dwag_tests';
    }
  }

  final pageIndexController = PageIndexController(total: 4);

  final projectNameController = TextEditingController();
  final projectDescriptionController = TextEditingController();

  final appNameController = TextEditingController();

  final androidBundleIdController = TextEditingController();
  final iosBundleIdController = TextEditingController();
  final macBundleIdController = TextEditingController();
  final windowsBundleIdController = TextEditingController();
  final linuxBundleIdController = TextEditingController();

  String? _directoryPath;

  String? get directoryPath => _directoryPath;

  Uint8List? _appIconAsBytes;

  Uint8List? get appIconAsBytes => _appIconAsBytes;

  final platformOptions = Map<PlatformOptions, bool>.fromEntries(
    PlatformOptions.values.map(
      (e) => MapEntry(e, true),
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

  void setAppIcon(Uint8List iconAsBytes) {
    _appIconAsBytes = iconAsBytes;
    notifyListeners();
  }

  void clearSelectedAppIcon() {
    _appIconAsBytes = null;
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

    androidBundleIdController.dispose();
    iosBundleIdController.dispose();
    windowsBundleIdController.dispose();
    linuxBundleIdController.dispose();
    macBundleIdController.dispose();

    super.dispose();
  }
}

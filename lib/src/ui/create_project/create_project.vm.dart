import 'package:flutter/cupertino.dart';

import '../../widgets/page_index_indicator.dart';

class CreateProjectViewModel extends ChangeNotifier {
  final pageIndexController = PageIndexController(total: 3);

  final projectNameController = TextEditingController();
  final projectDescriptionController = TextEditingController();

  final appNameController = TextEditingController();

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

  @override
  void dispose() {
    pageIndexController.dispose();

    projectNameController.dispose();
    projectDescriptionController.dispose();
    appNameController.dispose();

    super.dispose();
  }
}

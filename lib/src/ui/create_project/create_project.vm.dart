import 'package:flutter/cupertino.dart';

class CreateProjectViewModel extends ChangeNotifier {
  final projectNameController = TextEditingController();

  @override
  void dispose() {
    super.dispose();
    projectNameController.dispose();
  }
}

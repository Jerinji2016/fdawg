library fdawg_namer;

part 'helpers/validate_project_name.dart';

class FdawgNamer {
  FdawgNamer._();

  static void isValidName(String name) => _validateName(name);
}

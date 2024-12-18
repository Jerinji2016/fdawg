library fdawg_namer;

import 'package:fdawg_core/fdawg_core.dart';

import 'common/errors.dart';

part 'helpers/validate_app_name.dart';
part 'helpers/validate_project_name.dart';

class FdawgNamer {
  FdawgNamer._();

  static void isValidProjectName(String name) => _validateName(name);

  static void isValidAppName(String name, [List<PlatformOptions>? platforms]) => _validateAppName(name, platforms);
}

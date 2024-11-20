import 'package:flutter/material.dart';

import '../constants/assets.dart';
import '../helpers/helper.dart';

class Logo extends StatelessWidget {
  const Logo({
    this.size,
    this.noHero = false,
    super.key,
  });

  final bool noHero;
  final double? size;

  @override
  Widget build(BuildContext context) {
    return Hero(
      tag: noHero ? DateTime.now() : 'logo',
      child: Image.asset(
        Assets.logo,
        color: isDarkTheme ? Colors.white : Colors.black,
        height: size ?? 200,
      ),
    );
  }
}

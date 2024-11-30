import 'package:flutter/material.dart';

import '../config/theme_config.dart';

class PrimaryCard extends StatelessWidget {
  const PrimaryCard({
    required this.child,
    this.onTap,
    this.borderRadius = 16,
    super.key,
  });

  final Widget child;
  final double borderRadius;
  final VoidCallback? onTap;

  @override
  Widget build(BuildContext context) {
    return Material(
      color: Theme.of(context).cardColor,
      shape: RoundedRectangleBorder(
        borderRadius: BorderRadius.circular(borderRadius),
        side: const BorderSide(
          color: Colors.white24,
        ),
      ),
      child: InkWell(
        borderRadius: BorderRadius.circular(borderRadius),
        hoverDuration: hoverDuration,
        onTap: onTap,
        child: child,
      ),
    );
  }
}

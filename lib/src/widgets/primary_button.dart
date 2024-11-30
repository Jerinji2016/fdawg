import 'package:flutter/material.dart';

import '../config/theme_config.dart';

class PrimaryButton extends StatelessWidget {
  const PrimaryButton({
    required this.text,
    required this.onTap,
    this.padding = const EdgeInsets.symmetric(horizontal: 24, vertical: 12),
    this.side = BorderSide.none,
    this.color,
    Color? textColor,
    this.prefix,
    this.suffix,
    this.spacing = 4,
    super.key,
  }) : textColor = textColor ?? Colors.white;

  final String text;
  final VoidCallback onTap;
  final double spacing;
  final EdgeInsetsGeometry padding;
  final Color? color;
  final Color textColor;
  final BorderSide side;
  final Widget? prefix;
  final Widget? suffix;

  @override
  Widget build(BuildContext context) {
    return Material(
      color: color ?? Theme.of(context).colorScheme.primary,
      shape: RoundedRectangleBorder(
        side: side,
        borderRadius: BorderRadius.circular(16),
      ),
      child: InkWell(
        onTap: onTap,
        hoverDuration: hoverDuration,
        borderRadius: BorderRadius.circular(16),
        child: Padding(
          padding: padding,
          child: Row(
            mainAxisSize: MainAxisSize.min,
            children: [
              if (prefix != null) ...[
                prefix!,
                SizedBox(width: spacing),
              ],
              Text(
                text,
                style: TextStyle(color: textColor),
              ),
              if (suffix != null) ...[
                SizedBox(width: spacing),
                suffix!,
              ],
            ],
          ),
        ),
      ),
    );
  }
}

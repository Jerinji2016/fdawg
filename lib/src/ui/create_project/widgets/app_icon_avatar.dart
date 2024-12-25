import 'dart:io';
import 'dart:typed_data';

import 'package:file_picker/file_picker.dart';
import 'package:flutter/material.dart';

import '../../../config/theme_config.dart';

class AppIconAvatar extends StatefulWidget {
  const AppIconAvatar({
    required this.initialIconAsBytes,
    required this.onRemove,
    this.onImagePicked,
    super.key,
  });

  final Uint8List? initialIconAsBytes;
  final VoidCallback onRemove;
  final ValueChanged<Uint8List>? onImagePicked;

  @override
  State<AppIconAvatar> createState() => _AppIconAvatarState();
}

class _AppIconAvatarState extends State<AppIconAvatar> {
  bool _isHovered = false;

  @override
  Widget build(BuildContext context) {
    return InkWell(
      onTap: _onPickTapped,
      onHover: (value) => setState(() => _isHovered = value),
      hoverDuration: hoverDuration,
      borderRadius: BorderRadius.circular(24),
      child: Container(
        height: 120,
        width: 120,
        decoration: BoxDecoration(
          borderRadius: BorderRadius.circular(24),
          border: Border.all(
            color: Theme.of(context).dividerColor,
          ),
          boxShadow: [
            if (widget.initialIconAsBytes != null)
              BoxShadow(
                color: Colors.purple.shade300.withOpacity(0.8),
                spreadRadius: 3,
                blurRadius: 5,
              ),
          ],
        ),
        child: widget.initialIconAsBytes != null
            ? Stack(
                children: [
                  Positioned.fill(
                    child: ClipRRect(
                      borderRadius: BorderRadius.circular(24),
                      child: Image.memory(
                        widget.initialIconAsBytes!,
                        fit: BoxFit.cover,
                      ),
                    ),
                  ),
                  if (_isHovered)
                    Positioned(
                      top: 4,
                      right: 4,
                      child: IconButton.filled(
                        padding: EdgeInsets.zero,
                        onPressed: widget.onRemove,
                        iconSize: 14,
                        icon: const Icon(Icons.close),
                      ),
                    ),
                ],
              )
            : Center(
                child: Column(
                  mainAxisSize: MainAxisSize.min,
                  children: [
                    Icon(
                      Icons.add_a_photo_outlined,
                      color: Theme.of(context).disabledColor,
                    ),
                    const SizedBox(height: 4),
                    Text(
                      'Add App Icon',
                      style: TextStyle(
                        color: Theme.of(context).disabledColor,
                      ),
                    ),
                  ],
                ),
              ),
      ),
    );
  }

  Future<void> _onPickTapped() async {
    final result = await FilePicker.platform.pickFiles(
      type: FileType.image,
    );

    final appIconPath = result?.files.firstOrNull?.path;
    if (appIconPath == null) return;
    final file = File(appIconPath);
    final iconAsBytes = await file.readAsBytes();
    widget.onImagePicked?.call(iconAsBytes);
  }
}

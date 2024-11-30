import 'package:flutter/material.dart';
import 'package:provider/provider.dart';

class PageIndexController extends ChangeNotifier {
  PageIndexController({
    required this.total,
    int initialIndex = 0,
  })  : assert(initialIndex <= total, 'initialIndex should be <= total'),
        _currentIndex = initialIndex;

  final int total;

  int _currentIndex = 0;

  int get currentIndex => _currentIndex;

  void nextIndex() {
    _currentIndex++;
    notifyListeners();
  }

  void previousIndex() {
    _currentIndex--;
    notifyListeners();
  }

  void changeIndex(int newIndex) {
    _currentIndex = newIndex;
    notifyListeners();
  }
}

class PageIndexIndicator extends StatefulWidget {
  const PageIndexIndicator({
    required this.controller,
    super.key,
  });

  final PageIndexController controller;

  @override
  State<PageIndexIndicator> createState() => _PageIndexIndicatorState();
}

class _PageIndexIndicatorState extends State<PageIndexIndicator> {
  @override
  Widget build(BuildContext context) {
    return ChangeNotifierProvider.value(
      value: widget.controller,
      builder: (context, child) {
        final controller = Provider.of<PageIndexController>(context);

        return Row(
          mainAxisSize: MainAxisSize.min,
          children: List.generate(
            controller.total,
            (index) => _buildIndicatorDot(context, index),
          ),
        );
      },
    );
  }

  Widget _buildIndicatorDot(BuildContext context, int index) {
    final controller = Provider.of<PageIndexController>(context);

    final isSelected = index <= controller.currentIndex;
    final size = (isSelected ? 10 : 6).toDouble();
    final color = Theme.of(context).colorScheme.primary;

    return AnimatedContainer(
      duration: const Duration(milliseconds: 300),
      margin: const EdgeInsets.symmetric(horizontal: 4),
      decoration: BoxDecoration(
        color: isSelected ? color : Colors.transparent,
        borderRadius: BorderRadius.circular(size),
        border: Border.all(
          color: color,
        ),
      ),
      width: size,
      height: size,
    );
  }
}

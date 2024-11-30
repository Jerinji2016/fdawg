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
            (index) {
              final isSelected = index <= controller.currentIndex;
              return _buildIndicatorDot(isSelected);
            },
          ),
        );
      },
    );
  }

  Widget _buildIndicatorDot(bool isSelected) {
    final size = (isSelected ? 10 : 6).toDouble();
    return AnimatedContainer(
      duration: const Duration(milliseconds: 300),
      margin: const EdgeInsets.symmetric(horizontal: 4),
      decoration: BoxDecoration(
        color: isSelected ? Theme.of(context).primaryColor : Colors.transparent,
        borderRadius: BorderRadius.circular(size),
        border: Border.all(
          color: Theme.of(context).primaryColor,
        ),
      ),
      width: size,
      height: size,
    );
  }
}

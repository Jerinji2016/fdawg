import 'package:flutter/material.dart';
import 'package:flutter/scheduler.dart';

class FadeTransitionChild {
  FadeTransitionChild({
    required this.child,
    required this.index,
  });

  final Widget child;
  final int index;

  String get key => '$child-$index';
}

class FadeTransitionSwitcher extends StatefulWidget {
  const FadeTransitionSwitcher({
    required this.item,
    super.key = const ValueKey('FadeTransitionSwitcher'),
  });

  final FadeTransitionChild item;

  @override
  State<FadeTransitionSwitcher> createState() => _FadeTransitionSwitcherState();
}

class _FadeTransitionSwitcherState extends State<FadeTransitionSwitcher> with TickerProviderStateMixin {
  late final _entryAnimationController = AnimationController(
    vsync: this,
    duration: const Duration(milliseconds: 100),
    animationBehavior: AnimationBehavior.preserve,
  );

  late final _exitAnimationController = AnimationController(
    vsync: this,
    duration: const Duration(milliseconds: 300),
  );

  late FadeTransitionChild item;
  bool isReverse = false;

  @override
  void initState() {
    super.initState();
    item = widget.item;

    SchedulerBinding.instance.addPostFrameCallback((_) {
      _entryAnimationController.forward(from: 0);
    });
  }

  T getAnimation<T>(AnimationController controller, {required T begin, required T end}) {
    return controller
        .drive(
          Tween<T>(
            begin: begin,
            end: end,
          ).chain(CurveTween(curve: Curves.fastLinearToSlowEaseIn)),
        )
        .value;
  }

  @override
  Widget build(BuildContext context) {
    return AnimatedBuilder(
      animation: _exitAnimationController,
      builder: (context, child) {
        return Transform.translate(
          offset: getAnimation(
            _exitAnimationController,
            begin: Offset.zero,
            end: Offset(100 * (isReverse ? 1 : -1), 0),
          ),
          child: Opacity(
            opacity: getAnimation<double>(
              _exitAnimationController,
              begin: 1,
              end: 0,
            ),
            child: child,
          ),
        );
      },
      child: AnimatedBuilder(
        animation: _entryAnimationController,
        builder: (context, child) {
          return Transform.translate(
            offset: getAnimation(
              _entryAnimationController,
              begin: Offset(100 * (isReverse ? -1 : 1), 0),
              end: Offset.zero,
            ),
            child: Opacity(
              opacity: getAnimation<double>(
                _entryAnimationController,
                begin: 0,
                end: 1,
              ),
              child: item.child,
            ),
          );
        },
      ),
    );
  }

  @override
  void didUpdateWidget(covariant FadeTransitionSwitcher oldWidget) {
    if (oldWidget.item.key != widget.item.key) {
      transitionToNewChild(oldWidget.item.index > widget.item.index);
    }
    super.didUpdateWidget(oldWidget);
  }

  Future<void> transitionToNewChild(bool isReverse) async {
    debugPrint('_FadeTransitionSwitcherState.transitionToNewChild: 🐞$isReverse');
    this.isReverse = isReverse;
    await _exitAnimationController.forward();
    item = widget.item;
    _entryAnimationController.reset();
    _exitAnimationController.reset();
    await _entryAnimationController.forward();
  }

  @override
  void dispose() {
    if (_entryAnimationController.isAnimating) {
      _entryAnimationController.stop();
    }
    _entryAnimationController.dispose();

    if (_exitAnimationController.isAnimating) {
      _exitAnimationController.stop();
    }
    _exitAnimationController.dispose();

    super.dispose();
  }
}

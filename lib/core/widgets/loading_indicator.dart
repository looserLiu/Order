import 'package:flutter/material.dart';

enum LoadingIndicatorSize { small, medium, large }

class LoadingIndicator extends StatelessWidget {
  final LoadingIndicatorSize size;
  final Color? color;
  final String? message;

  const LoadingIndicator({
    super.key,
    this.size = LoadingIndicatorSize.medium,
    this.color,
    this.message,
  });

  @override
  Widget build(BuildContext context) {
    final colorScheme = Theme.of(context).colorScheme;

    return Center(
      child: Column(
        mainAxisSize: MainAxisSize.min,
        children: [
          SizedBox(
            width: _diameter,
            height: _diameter,
            child: CircularProgressIndicator(
              color: color ?? colorScheme.primary,
              strokeWidth: _strokeWidth,
            ),
          ),
          if (message != null) ...[
            const SizedBox(height: 16),
            Text(
              message!,
              style: Theme.of(context).textTheme.bodyMedium?.copyWith(
                    color: colorScheme.onSurfaceVariant,
                  ),
            ),
          ],
        ],
      ),
    );
  }

  double get _diameter {
    switch (size) {
      case LoadingIndicatorSize.small:
        return 24;
      case LoadingIndicatorSize.medium:
        return 40;
      case LoadingIndicatorSize.large:
        return 64;
    }
  }

  double get _strokeWidth {
    switch (size) {
      case LoadingIndicatorSize.small:
        return 2;
      case LoadingIndicatorSize.medium:
        return 3;
      case LoadingIndicatorSize.large:
        return 4;
    }
  }
}

import 'package:flutter/material.dart';

enum AppButtonVariant { primary, secondary, outline, text }

enum AppButtonSize { small, medium, large }

class AppButton extends StatelessWidget {
  final String label;
  final VoidCallback? onPressed;
  final AppButtonVariant variant;
  final AppButtonSize size;
  final IconData? icon;
  final bool isLoading;
  final bool isExpanded;

  const AppButton({
    super.key,
    required this.label,
    this.onPressed,
    this.variant = AppButtonVariant.primary,
    this.size = AppButtonSize.medium,
    this.icon,
    this.isLoading = false,
    this.isExpanded = false,
  });

  @override
  Widget build(BuildContext context) {
    final colorScheme = Theme.of(context).colorScheme;

    ButtonStyle style;
    switch (variant) {
      case AppButtonVariant.primary:
        style = ElevatedButton.styleFrom(
          backgroundColor: colorScheme.primary,
          foregroundColor: colorScheme.onPrimary,
        );
      case AppButtonVariant.secondary:
        style = ElevatedButton.styleFrom(
          backgroundColor: colorScheme.secondary,
          foregroundColor: colorScheme.onSecondary,
        );
      case AppButtonVariant.outline:
        style = OutlinedButton.styleFrom(
          foregroundColor: colorScheme.primary,
        );
      case AppButtonVariant.text:
        style = TextButton.styleFrom(
          foregroundColor: colorScheme.primary,
        );
    }

    EdgeInsetsGeometry padding;
    switch (size) {
      case AppButtonSize.small:
        padding = const EdgeInsets.symmetric(horizontal: 12, vertical: 6);
      case AppButtonSize.medium:
        padding = const EdgeInsets.symmetric(horizontal: 24, vertical: 12);
      case AppButtonSize.large:
        padding = const EdgeInsets.symmetric(horizontal: 32, vertical: 16);
    }

    Widget buttonChild = isLoading
        ? SizedBox(
            height: 20,
            width: 20,
            child: CircularProgressIndicator(
              strokeWidth: 2,
              color: variant == AppButtonVariant.primary
                  ? colorScheme.onPrimary
                  : colorScheme.primary,
            ),
          )
        : Row(
            mainAxisSize: MainAxisSize.min,
            children: [
              if (icon != null) ...[
                Icon(icon, size: _iconSize),
                const SizedBox(width: 8),
              ],
              Text(label),
            ],
          );

    Widget button;
    switch (variant) {
      case AppButtonVariant.primary:
      case AppButtonVariant.secondary:
        button = ElevatedButton(
          onPressed: isLoading ? null : onPressed,
          style: style.copyWith(
            padding: WidgetStateProperty.all(padding),
          ),
          child: buttonChild,
        );
      case AppButtonVariant.outline:
        button = OutlinedButton(
          onPressed: isLoading ? null : onPressed,
          style: style.copyWith(
            padding: WidgetStateProperty.all(padding),
          ),
          child: buttonChild,
        );
      case AppButtonVariant.text:
        button = TextButton(
          onPressed: isLoading ? null : onPressed,
          style: style.copyWith(
            padding: WidgetStateProperty.all(padding),
          ),
          child: buttonChild,
        );
    }

    return isExpanded
        ? SizedBox(width: double.infinity, child: button)
        : button;
  }

  double get _iconSize {
    switch (size) {
      case AppButtonSize.small:
        return 16;
      case AppButtonSize.medium:
        return 20;
      case AppButtonSize.large:
        return 24;
    }
  }
}

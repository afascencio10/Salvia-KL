import 'package:flutter/material.dart';
import '../../global/my_colors.dart';

class LinkText extends StatefulWidget {
  final String text;
  final VoidCallback onLinkTap;
  final double fontSize;
  final Color color;
  final double letterSpacing;
  final bool asOutlinedButton;
  final Color outlinedBorderColor;
  final Color outlinedBackgroundColor;
  final double borderRadius;
  final double outlinedBorderWidth;
  final EdgeInsetsGeometry padding;

  const LinkText({
    Key? key,
    required this.text,
    required this.onLinkTap,
    this.fontSize = 16.0,
    this.color = AppColors.salviaPurpura3,
    this.letterSpacing = 0.0,
    this.asOutlinedButton = false,
    this.outlinedBorderColor = AppColors.salviaLinkOutline,
    this.outlinedBackgroundColor = Colors.transparent,
    this.borderRadius = 18,
    this.outlinedBorderWidth = 1.2,
    this.padding = const EdgeInsets.symmetric(vertical: 14, horizontal: 18),
  }) : super(key: key);

  @override
  State<LinkText> createState() => _LinkTextState();
}

class _LinkTextState extends State<LinkText> {
  static const Duration _pressDuration = Duration(milliseconds: 90);
  static const Duration _releaseDuration = Duration(milliseconds: 110);
  bool _isPressed = false;
  bool _isTriggering = false;

  void _setPressed(bool value) {
    if (_isPressed == value) return;
    setState(() {
      _isPressed = value;
    });
  }

  Future<void> _handleTap() async {
    if (_isTriggering) return;
    _isTriggering = true;
    _setPressed(false);
    await Future<void>.delayed(_releaseDuration);
    if (!mounted) return;
    widget.onLinkTap();
    _isTriggering = false;
  }

  @override
  Widget build(BuildContext context) {
    final textWidget = Text(
      widget.text,
      style: TextStyle(
        decoration: widget.asOutlinedButton
            ? TextDecoration.none
            : TextDecoration.underline,
        color: widget.color,
        fontSize: widget.fontSize,
        letterSpacing: widget.letterSpacing,
        fontWeight: FontWeight.w400,
      ),
      textAlign: TextAlign.center,
    );

    final buttonChild = widget.asOutlinedButton
        ? Ink(
            padding: widget.padding,
            decoration: BoxDecoration(
              color: widget.outlinedBackgroundColor,
              borderRadius: BorderRadius.circular(widget.borderRadius),
              border: Border.all(
                color: widget.outlinedBorderColor,
                width: widget.outlinedBorderWidth,
              ),
            ),
            child: textWidget,
          )
        : textWidget;

    return Listener(
      onPointerDown: (_) => _setPressed(true),
      onPointerUp: (_) => _setPressed(false),
      onPointerCancel: (_) => _setPressed(false),
      child: AnimatedScale(
        scale: _isPressed ? 0.982 : 1.0,
        duration: _isPressed ? _pressDuration : _releaseDuration,
        curve: Curves.easeOutCubic,
        child: AnimatedContainer(
          duration: _isPressed ? _pressDuration : _releaseDuration,
          transform: Matrix4.translationValues(0, _isPressed ? 2 : 0, 0),
          child: Material(
            color: Colors.transparent,
            child: InkWell(
              onTap: _handleTap,
              borderRadius: BorderRadius.circular(widget.borderRadius),
              child: buttonChild,
            ),
          ),
        ),
      ),
    );
  }
}

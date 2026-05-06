import 'package:flutter/material.dart';

class NavigationButton extends StatefulWidget {
  final String text;
  final Color buttonColor;
  final Color textColor;

  final VoidCallback? onPressed;
  final IconData icon;
  final Color iconColor;
  final double iconSize;
  final double fontSize;
  final double buttonHeight;
  final double elevation;
  final bool iconOnLeft;
  final bool contentLeftAligned;
  final double leftContentPadding;

  const NavigationButton({
    super.key,
    required this.text,
    this.buttonColor = Colors.blue,
    this.textColor = Colors.white,
    this.onPressed,
    this.icon = Icons.arrow_forward,
    this.iconColor = Colors.white,
    this.iconSize = 24.0,
    this.fontSize = 18.0,
    this.buttonHeight = 70.0,
    this.elevation = 5.0,
    this.iconOnLeft = false,
    this.contentLeftAligned = false,
    this.leftContentPadding = 18.0,
  });

  @override
  State<NavigationButton> createState() => _NavigationButtonState();
}

class _NavigationButtonState extends State<NavigationButton> {
  static const Duration _pressDuration = Duration(milliseconds: 90);
  static const Duration _releaseDuration = Duration(milliseconds: 110);
  bool _isPressed = false;
  bool _isTriggering = false;

  bool get _isEnabled => widget.onPressed != null;

  void _setPressed(bool value) {
    if (!_isEnabled) return;
    if (_isPressed == value) return;
    setState(() {
      _isPressed = value;
    });
  }

  Future<void> _handleTap() async {
    if (!_isEnabled || _isTriggering) return;
    _isTriggering = true;
    _setPressed(false);
    await Future<void>.delayed(_releaseDuration);
    if (!mounted) return;
    widget.onPressed?.call();
    _isTriggering = false;
  }

  @override
  Widget build(BuildContext context) {
    final textWidget = Text(
      widget.text,
      style: TextStyle(
        color: widget.textColor,
        fontSize: widget.fontSize,
      ),
      textAlign: widget.contentLeftAligned ? TextAlign.left : TextAlign.center,
    );

    final iconWidget = Icon(
      widget.icon,
      color: widget.iconColor,
      size: widget.iconSize,
    );

    final orderedChildren = widget.iconOnLeft
        ? <Widget>[
            iconWidget,
            const SizedBox(width: 10),
            textWidget,
          ]
        : <Widget>[
            textWidget,
            const SizedBox(width: 8),
            iconWidget,
          ];

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
          child: ElevatedButton(
            style: ElevatedButton.styleFrom(
              backgroundColor: widget.buttonColor,
              disabledBackgroundColor: Colors.grey.shade300,
              disabledForegroundColor: Colors.black38,
              shape: RoundedRectangleBorder(
                borderRadius: BorderRadius.circular(30),
              ),
              minimumSize: Size(double.infinity, widget.buttonHeight),
              elevation: _isPressed ? 1 : widget.elevation,
            ),
            onPressed: _isEnabled ? _handleTap : null,
            child: Row(
              mainAxisSize: widget.contentLeftAligned
                  ? MainAxisSize.max
                  : MainAxisSize.min,
              mainAxisAlignment: widget.contentLeftAligned
                  ? MainAxisAlignment.start
                  : MainAxisAlignment.center,
              children: [
                if (widget.contentLeftAligned)
                  SizedBox(width: widget.leftContentPadding),
                ...orderedChildren,
              ],
            ),
          ),
        ),
      ),
    );
  }
}

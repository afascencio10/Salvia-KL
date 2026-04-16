import 'package:flutter/material.dart';
import '../../global/my_colors.dart';

class InfoButton extends StatelessWidget {
  final String infoMessage;
  final IconData buttonIcon;
  final Color buttonBackgroundColor;
  final Color iconColor;
  final double iconSize;
  final bool boxedStyle;
  final double boxSize;
  final double borderRadius;
  final Color borderColor;
  final double borderWidth;
  final Color snackBackgroundColor;
  final Color snackTextColor;
  final IconData snackIcon;
  final Color snackIconColor;
  final Duration snackDuration;

  const InfoButton({
    super.key,
    required this.infoMessage,
    this.buttonIcon = Icons.help_outline_rounded,
    this.buttonBackgroundColor = AppColors.salviaPurpura5,
    this.iconColor = AppColors.salviaPurpura3, // Color predeterminado del ícono
    this.iconSize = 30, // Tamaño predeterminado del ícono
    this.boxedStyle = false,
    this.boxSize = 82,
    this.borderRadius = 20,
    this.borderColor = const Color(0xFFD0D0D0),
    this.borderWidth = 1.2,
    this.snackBackgroundColor = AppColors.salviaPurpura1,
    this.snackTextColor = Colors.white,
    this.snackIcon = Icons.info_rounded,
    this.snackIconColor = Colors.white,
    this.snackDuration = const Duration(seconds: 3),
  });

  @override
  Widget build(BuildContext context) {
    final widgetButton = Container(
      width: boxedStyle ? boxSize : null,
      height: boxedStyle ? boxSize : null,
      decoration: BoxDecoration(
        color: buttonBackgroundColor,
        shape: boxedStyle ? BoxShape.rectangle : BoxShape.circle,
        borderRadius: boxedStyle ? BorderRadius.circular(borderRadius) : null,
        border: boxedStyle
            ? Border.all(color: borderColor, width: borderWidth)
            : null,
      ),
      child: IconButton(
        icon: Icon(buttonIcon, color: iconColor, size: iconSize),
        onPressed: () {
          ScaffoldMessenger.of(context).hideCurrentSnackBar();
          ScaffoldMessenger.of(context).showSnackBar(
            SnackBar(
              behavior: SnackBarBehavior.floating,
              backgroundColor: snackBackgroundColor,
              shape: RoundedRectangleBorder(
                borderRadius: BorderRadius.circular(14),
              ),
              duration: snackDuration,
              content: Row(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Icon(
                    snackIcon,
                    color: snackIconColor,
                    size: 22,
                  ),
                  const SizedBox(width: 10),
                  Expanded(
                    child: Text(
                      infoMessage,
                      style: TextStyle(
                        color: snackTextColor,
                        fontSize: 14,
                        fontWeight: FontWeight.w400,
                        height: 1.3,
                      ),
                    ),
                  ),
                ],
              ),
            ),
          );
        },
        padding: EdgeInsets.zero,
        alignment: Alignment.center,
      ),
    );

    return widgetButton;
  }
}

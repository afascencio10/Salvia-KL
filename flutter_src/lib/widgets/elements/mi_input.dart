import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import '../../global/my_colors.dart';

class IconoYTextField extends StatefulWidget {
  final TextEditingController? controller;
  final double fontSize;
  final Color fontColor;
  final String hintText;
  final bool isPassword;
  final bool isNumeric;
  final TextInputType keyboardType;
  final int? maxLength;
  final int? minLines;
  final Function(String)? onChanged;

  const IconoYTextField({
    Key? key,
    this.controller,
    this.fontSize = 16.0,
    this.fontColor = const Color(0xff45008b),
    this.hintText = '',
    this.isPassword = false,
    this.isNumeric = false,
    this.keyboardType = TextInputType.text,
    this.maxLength,
    this.minLines = 1,
    this.onChanged,
  }) : super(key: key);

  @override
  IconoYTextFieldState createState() => IconoYTextFieldState();
}

class IconoYTextFieldState extends State<IconoYTextField> {
  @override
  Widget build(BuildContext context) {
    return Row(
      mainAxisAlignment: MainAxisAlignment.start,
      children: <Widget>[
        const SizedBox(width: 20.0),
        Expanded(
          child: Container(
            padding: const EdgeInsets.symmetric(vertical: 8, horizontal: 16),
            decoration: BoxDecoration(
              border: Border.all(color: const Color(0xffe7d9f6), width: 2),
              borderRadius: BorderRadius.circular(20),
            ),
            child: TextField(
              controller: widget.controller, // ✅ NUEVO
              style: TextStyle(
                fontSize: widget.fontSize,
                color: widget.fontColor,
              ),
              decoration: InputDecoration(
                border: InputBorder.none,
                hintText: widget.hintText,
                hintStyle: const TextStyle(
                  color: AppColors.salviaPlaceholderColor,
                  fontSize: AppColors.salviaPlaceholderFontSize,
                  fontWeight: FontWeight.w400,
                ),
              ),
              obscureText: widget.isPassword,
              keyboardType:
                  widget.isNumeric ? TextInputType.number : widget.keyboardType,
              inputFormatters: widget.isNumeric
                  ? [FilteringTextInputFormatter.digitsOnly]
                  : [],
              maxLength: widget.maxLength,
              maxLines: widget.isPassword ? 1 : null,
              minLines: widget.isPassword ? 1 : widget.minLines,
              onChanged: widget.onChanged,
            ),
          ),
        ),
      ],
    );
  }
}

import 'package:flutter/material.dart';
import '../../global/my_colors.dart';

Widget lineaDivisoria() {
  return const Padding(
    padding: EdgeInsets.symmetric(
        vertical: 20), // Agrega 20 píxeles arriba y abajo
    child: Divider(
      color: AppColors.salviaPurpura4,
      height: 16,
      thickness: 0,
      indent: 0,
      endIndent: 0,
    ),
  );
}

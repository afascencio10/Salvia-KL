import 'package:flutter/material.dart';
import '../../global/my_colors.dart';

Widget subtitulo(String texto) {
  return Text(
    texto,
    style: const TextStyle(
      color: AppColors.salviaSubtitulos, // Color del texto
      letterSpacing: 0.1, // Espaciado entre letras
      fontWeight: FontWeight.normal, // Peso de la fuente
      fontSize: 16.0, // Tamaño de la fuente
      // Puedes añadir más propiedades de estilo si lo necesitas
    ),
    textAlign: TextAlign.center, // Alineación del texto
  );
}

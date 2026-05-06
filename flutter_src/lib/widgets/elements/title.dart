import 'package:flutter/material.dart';
import '../../global/my_colors.dart';

Widget miTitulo(String texto) {
  return Text(
    texto,
    style: const TextStyle(
      color: AppColors.salviaTitulos, // Color del texto
      letterSpacing: 0.5, // Espaciado entre letras
      fontWeight: FontWeight.bold, // Peso de la fuente
      fontSize: 26.0, // Tamaño de la fuente
      // Puedes añadir más propiedades de estilo si lo necesitas
    ),
    textAlign: TextAlign.center, // Alineación del texto
  );
}

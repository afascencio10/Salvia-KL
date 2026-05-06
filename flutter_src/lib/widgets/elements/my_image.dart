import 'package:flutter/material.dart';

class ImagenConAspecto extends StatelessWidget {
  final String nombreImagen;
  final double aspectRatioWidth;
  final double aspectRatioHeight;

  const ImagenConAspecto({
    Key? key,
    required this.nombreImagen,
    required this.aspectRatioWidth,
    required this.aspectRatioHeight,
  }) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return AspectRatio(
      aspectRatio: aspectRatioWidth / aspectRatioHeight,
      child: Container(
        width:
            MediaQuery.of(context).size.width, // Ancho igual al de la pantalla
        decoration: BoxDecoration(
          image: DecorationImage(
            image: AssetImage(nombreImagen), // Ruta de la imagen
            fit: BoxFit
                .fitWidth, // La imagen se ajusta al ancho, manteniendo proporciones
          ),
          borderRadius: BorderRadius.circular(10), // Esquinas redondeadas
        ),
      ),
    );
  }
}

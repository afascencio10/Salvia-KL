import 'package:flutter/material.dart';
import '../global/my_screens.dart';
import '../global/my_colors.dart';
import '../widgets/elements/info_button_snack.dart';
import '../widgets/elements/go_button.dart';
import '../widgets/elements/espacio.dart';
import '../widgets/elements/title.dart';
import '../widgets/elements/subtitle.dart';
import "../widgets/elements/linea_horizontal.dart";
import '../data/data.dart';
import '../widgets/elements/badge.dart';

class CasosScreen extends StatefulWidget {
  const CasosScreen({super.key});

  @override
  State<CasosScreen> createState() => _CasosScreenState();
}

class _CasosScreenState extends State<CasosScreen> {
  @override
  void initState() {
    super.initState();
  }

  //llamada desde el botón de navegación a otras páginas
  void stopAudioIfNeeded() {
    // Audio temporalmente desactivado.
  }

  // void _toggleAudio(String fileName, int buttonNumber) {
  //   String path = "lib/assets/sounds/$fileName";
  //   if (isPlaying && playingButton == buttonNumber) {
  //     audioPlayer.stop();
  //   } else {
  //     audioPlayer
  //         .open(Audio(path), autoStart: true, showNotification: true)
  //         .then((_) {
  //       setState(() {
  //         playingButton = buttonNumber;
  //         isPlaying = true;
  //       });
  //     }).catchError((error) {
  //       //print("Error al reproducir audio: $error");
  //       setState(() {
  //         playingButton = null;
  //         isPlaying = false;
  //       });
  //     });
  //   }
  // }

  @override
  void dispose() {
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text(
          "Consultar | Casos",
          style: TextStyle(
            color: Color(0xffffffff),
            fontSize: 20,
          ),
        ),
        backgroundColor: AppColors.salviaTurquesa1,
        elevation: 0,
      ),
      body: cuerpoReportar(context),
    );
  }

  Widget cuerpoReportar(BuildContext context) {
    return Stack(
      children: [
        Container(
          decoration: const BoxDecoration(
            color: AppColors.salviaTurquesa1,
          ),
        ),
        // Agregar aquí el Container con la imagen
        Container(
          width: MediaQuery.of(context).size.width,
          height: MediaQuery.of(context).size.height *
              0.2, // Ajusta esta altura según tus necesidades
          decoration: const BoxDecoration(
            image: DecorationImage(
              image: AssetImage('lib/assets/images/brand.png'),
              fit: BoxFit.cover,
              alignment: Alignment.topCenter,
            ),
          ),
        ),
        // Positioned widget para el ícono con opacidad
        const Positioned(
          top: 16.0, // Padding desde arriba
          left: 20.0, // Padding desde la izquierda
          child: Opacity(
            opacity: 0.3, // Ajusta el nivel de opacidad aquí
            child: Icon(
              Icons.person, // Reemplaza con el ícono que deseas usar
              color: Colors.white,
              size: 30.0, // Ajusta el tamaño según tus necesidades
            ),
          ),
        ),
        SingleChildScrollView(
          // Añadir SingleChildScrollView aquí
          child: ConstrainedBox(
            constraints: BoxConstraints(
              minHeight: MediaQuery.of(context).size.height,
            ),
            child: Container(
              margin: const EdgeInsets.fromLTRB(0, 50, 0, 0),
              padding: const EdgeInsets.all(0),
              decoration: BoxDecoration(
                color: const Color(0xffffffff),
                shape: BoxShape.rectangle,
                borderRadius: const BorderRadius.only(
                  topLeft: Radius.circular(30.0),
                  topRight: Radius.circular(30.0),
                ),
                border: Border.all(color: const Color(0x4d9e9e9e), width: 1),
              ),
              child: Padding(
                padding: const EdgeInsets.symmetric(horizontal: 20),
                child: Column(
                  mainAxisAlignment: MainAxisAlignment.start,
                  crossAxisAlignment: CrossAxisAlignment.stretch,
                  //--------------------------------------------------------------
                  //INTERFAZ COMPONENTES------------------------------------------
                  //--------------------------------------------------------------
                  children: [
                    espacio(50),
                    miTitulo("Bienvenid@"),
                    espacio(20),
                    subtitulo(
                        "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Quisque pulvinar porta nibh at tempor."),
                    espacio(20),

                    _buildComponentRowSoundSnackButton(
                      1,
                      "sonidoGeneral.mp3",
                      "Información sobre Sonido 1", //Snack
                      "Actualizar", //Button
                      () {
                        stopAudioIfNeeded();
                        Navigator.of(context).push(
                          MaterialPageRoute(
                              builder: (context) => const InicioScreen()),
                        );
                      },
                      AppColors.salviaTurquesa1, //Button
                      Colors.white, //Button
                      Icons.refresh_rounded, // Ícono del botón
                      Colors.white, //Icon
                      24.0, //Icon
                    ),

                    lineaDivisoria(),

                    espacio(20),
                    miTitulo("Lista de casos"),
                    espacio(20),

                    ListView.builder(
                      physics: const NeverScrollableScrollPhysics(),
                      shrinkWrap: true,
                      itemCount: options.length,
                      itemBuilder: (BuildContext context, int index) {
                        //print("Índice actual: $index");
                        return Column(
                          children: [
                            Padding(
                              padding:
                                  const EdgeInsets.symmetric(vertical: 8.0),
                              child: Row(
                                children: [
                                  Container(
                                    width: 5.0, // Ancho de la línea vertical
                                    height:
                                        60.0, // Altura aproximada de la línea (ajústala según tus necesidades)
                                    color: AppColors
                                        .salviaPurpura3, // Color de la línea
                                  ),
                                  const SizedBox(
                                      width:
                                          10), // Espacio entre la línea y el texto
                                  Expanded(
                                    child: Column(
                                      crossAxisAlignment:
                                          CrossAxisAlignment.start,
                                      children: [
                                        const SizedBox(
                                          height: 15,
                                        ),
                                        const Text(
                                          "Caso",
                                          style: TextStyle(
                                            fontSize: 18,
                                            color: AppColors.salviaTurquesa1,
                                          ),
                                        ),
                                        const SizedBox(
                                          height: 15,
                                        ),
                                        Text(
                                          options[index].title,
                                          style: const TextStyle(
                                            fontSize: 16,
                                            color: AppColors.salviaTitulos,
                                            fontWeight: FontWeight.normal,
                                          ),
                                        ),
                                        const SizedBox(
                                          height: 10,
                                        ),
                                      ],
                                    ),
                                  ),
                                  options[index].pendingMessages > 0
                                      ? MyBadge(
                                          messages:
                                              options[index].pendingMessages)
                                      : const SizedBox
                                          .shrink(), // Si no hay mensajes, no muestra nada
                                ],
                              ),
                            ),
                            _buildComponentRowSoundSnackButton(
                              2,
                              "sonidoGeneral.mp3",
                              "Información sobre Sonido 1", //Snack
                              "Leer Caso", //Button
                              () {
                                stopAudioIfNeeded();
                                Navigator.of(context).push(
                                  MaterialPageRoute(
                                      builder: (context) =>
                                          const InicioScreen()),
                                );
                              },
                              AppColors.salviaPurpura3, //Button
                              Colors.white, //Button
                              Icons.sort, // Ícono del botón
                              Colors.white, //Icon
                              24.0, //Icon
                            ),
                            const SizedBox(
                              height: 30,
                            ),
                            const Divider(height: 0),
                          ],
                        );
                      },
                    ),

                    //------------------------------------------

                    espacio(30),
                    // ... Repetir IconoYDropdown y otros widgets según sea necesario
                  ],
                ),
              ),
            ),
          ),
        ),
      ],
    );
  }


  Widget _buildComponentRowSoundSnackButton(
      int buttonNumber,
      String audioName,
      String infoMessage,
      String
          navButtonText, // Nuevo parámetro para el texto del botón de navegación
      VoidCallback
          onPressed, // Callback para la acción del botón // Nuevo parámetro para la ruta de navegación
      Color navButtonColor, // Nuevo parámetro para el color del botón
      Color
          navButtonTextColor, // Nuevo parámetro para el color del texto del botón
      IconData navButtonIcon, // Nuevo parámetro para el ícono
      Color navButtonIconColor, // Nuevo parámetro para el color del ícono
      double navButtonIconSize // Nuevo parámetro para el tamaño del ícono
      ) {
    return Row(
      mainAxisAlignment: MainAxisAlignment.start,
      crossAxisAlignment: CrossAxisAlignment.center,
      mainAxisSize: MainAxisSize.max,
      children: [
        Padding(
          padding: const EdgeInsets.fromLTRB(0, 0, 18, 0),
          child: InfoButton(
            infoMessage: infoMessage,
          ),
        ),
        Expanded(
          flex: 1,
          child: NavigationButton(
            text: navButtonText,
            onPressed: onPressed, // Usamos el callback aquí
            buttonColor: navButtonColor,
            textColor: navButtonTextColor,
            icon: navButtonIcon,
            iconColor: navButtonIconColor,
            iconSize: navButtonIconSize,
          ),
        ),
      ],
    );
  }



}

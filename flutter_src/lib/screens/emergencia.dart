import 'package:flutter/material.dart';
import '../global/my_screens.dart';
import '../global/my_colors.dart';
import '../widgets/elements/info_button_snack.dart';
import '../widgets/elements/go_button.dart';
import '../widgets/elements/espacio.dart';
import '../widgets/elements/title.dart';
import '../widgets/elements/subtitle.dart';

class EmergenciaScreen extends StatefulWidget {
  const EmergenciaScreen({super.key});
  @override
  State<EmergenciaScreen> createState() => _EmergenciaScreenState();
}

class _EmergenciaScreenState extends State<EmergenciaScreen> {
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
          "Emergencia",
          style: TextStyle(
            color: Color(0xffffffff),
            fontSize: 20,
          ),
        ),
        backgroundColor: AppColors.salviaAmarillo1,
        elevation: 0,
      ),
      body: cuerpoEmergencia(context),
    );
  }

  Widget cuerpoEmergencia(BuildContext context) {
    return Stack(
      children: [
        Container(
          decoration: const BoxDecoration(
            color: AppColors.salviaAmarillo1, // Color azul de fondo
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
              Icons.call_rounded, // Reemplaza con el ícono que deseas usar
              color: Colors.white,
              size: 30.0, // Ajusta el tamaño según tus necesidades
            ),
          ),
        ),
        SizedBox(
          height: MediaQuery.of(context).size.height,
          width: MediaQuery.of(context).size.width,
          child: Stack(
            alignment: Alignment.topCenter,
            children: [
              Container(
                margin: const EdgeInsets.fromLTRB(0, 50, 0, 0),
                padding: const EdgeInsets.all(0),
                width: MediaQuery.of(context).size.width,
                height: MediaQuery.of(context).size.height,
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
                  padding: const EdgeInsets.symmetric(
                      vertical: 0, horizontal: 20), //padding contenedor
                  child: Column(
                    mainAxisAlignment: MainAxisAlignment.start,
                    crossAxisAlignment: CrossAxisAlignment.stretch,
                    mainAxisSize: MainAxisSize.max,
                    children: [
                      espacio(50), // Espacio después de label_A
                      miTitulo("Comunícate"),
                      espacio(20), // Espacio después de label_A
                      subtitulo(
                          "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Quisque pulvinar porta nibh at tempor."),
                      espacio(20),

                      const SizedBox(height: 20),
                      _buildComponentRowSoundSnackButton(
                        1,
                        "sonidoGeneral.mp3",
                        "Información sobre Sonido 3", //Snack
                        "llamar", //Button
                        () {
                          // onPressed (VoidCallback)
                          Navigator.of(context).push(
                            MaterialPageRoute(
                                builder: (context) => const InicioScreen()),
                          );
                        },
                        AppColors.salviaAmarillo1, //Button
                        Colors.white, //Button
                        Icons.call, // Ícono del botón
                        Colors.white, //Icon
                        24.0, //Icon
                      ),
                    ],
                  ),
                ),
              ),
            ],
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
            //onButtonPressed: stopAudioIfNeeded,
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

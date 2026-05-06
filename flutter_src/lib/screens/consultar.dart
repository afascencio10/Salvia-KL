import 'package:flutter/material.dart';
import '../global/my_screens.dart';
import '../global/my_colors.dart';
import '../widgets/elements/info_button_snack.dart';
import '../widgets/elements/go_button.dart';
import '../widgets/elements/mi_input.dart';
import '../widgets/elements/espacio.dart';
import '../widgets/elements/title.dart';
import '../widgets/elements/subtitle.dart';

class ConsultarScreen extends StatefulWidget {
  const ConsultarScreen({super.key});
  @override
  State<ConsultarScreen> createState() => _ConsultarScreenState();
}

class _ConsultarScreenState extends State<ConsultarScreen> {
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
      extendBodyBehindAppBar:
          true, // Esto permite que el cuerpo se extienda detrás del AppBar
      appBar: AppBar(
        title: const Text(
          "Consultar",
          style: TextStyle(
            color: Color(0xffffffff),
            fontSize: 20,
          ),
        ),
        backgroundColor: AppColors.salviaTransparente, // tu color con opacidad
        elevation: 0,
        automaticallyImplyLeading: true,
        iconTheme: const IconThemeData(
          color: Colors.white, // Color blanco para el botón back
        ),
      ),
      body: cuerpoInicio(context),
    );
  }

  Widget cuerpoInicio(BuildContext context) {
    return Stack(
      children: [
        Container(
          decoration: const BoxDecoration(
            color: AppColors.salviaTurquesa1, // Color azul de fondo
          ),
        ),
        // Agregar aquí el Container con la imagen
        Container(
          width: MediaQuery.of(context).size.width,
          height: MediaQuery.of(context).size.height *
              0.2, // Ajusta esta altura según tus necesidades
          decoration: const BoxDecoration(
            image: DecorationImage(
              image: AssetImage('lib/assets/images/brand-consultar.png'),
              fit: BoxFit.cover,
              alignment: Alignment.topCenter,
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
                margin: const EdgeInsets.fromLTRB(0, 100, 0, 0),
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
                    //--------------------------------------------------------------
                    //INTERFAZ COMPONENTES------------------------------------------
                    //--------------------------------------------------------------
                    children: [
                      espacio(50), // Espacio después de label_A
                      miTitulo("Iniciar Sesión"),
                      espacio(20), // Espacio después de label_A
                      subtitulo(
                          "Ingresa el usuario y contraseña brindado por el personal de Salvia, para hacer seguimiento a tu caso."),
                      espacio(20),

                      const SizedBox(height: 20),
                      _componentPlayInput(
                          1, // buttonNumber
                          "sonidoGeneral.mp3", // audioName
                          "Usuario", // hintText para el TextField
                          16.0, // Tamaño de fuente para el TextField
                          Colors.black, // Color de fuente para el TextField
                          false, // isPasswordField
                          false, // isNumericField
                          TextInputType.text, // keyboardType
                          null, // maxLength
                          1 // minLines
                          ),
                      const SizedBox(height: 20),
                      _componentPlayInput(
                          2, // buttonNumber
                          "sonidoGeneral.mp3", // audioName
                          "Contraseña", // hintText para el TextField
                          16.0, // Tamaño de fuente para el TextField
                          Colors.black, // Color de fuente para el TextField
                          true, // isPasswordField
                          false, // isNumericField
                          TextInputType.text, // keyboardType
                          null, // maxLength
                          1 // minLines
                          ),
                      const SizedBox(height: 20),
                      _buildComponentRowSoundSnackButton(
                        3,
                        "sonidoGeneral.mp3",
                        "Información sobre Sonido 1", //Snack
                        "Ingresar", //Button
                        () {
                          // onPressed (VoidCallback)
                          stopAudioIfNeeded();
                          Navigator.pushAndRemoveUntil(
                            context,
                            MaterialPageRoute(
                                builder: (context) => const CasosScreen()),
                            (Route<dynamic> route) =>
                                false, // Esto asegura que todas las pantallas anteriores sean eliminadas
                          );
                        },
                        AppColors.salviaTurquesa1, //Button
                        Colors.white, //Button
                        Icons.key, // Ícono del botón
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

  Widget _componentPlayInput(
      int buttonNumber,
      String audioName,
      String
          hintText, // Nuevo parámetro para el texto del placeholder del campo de texto
      double
          textFieldFontSize, // Nuevo parámetro para el tamaño de fuente del campo de texto
      Color
          textFieldFontColor, // Nuevo parámetro para el color de la fuente del campo de texto
      bool
          isPasswordField, // Nuevo parámetro para determinar si es un campo de contraseña
      bool
          isNumericField, // Nuevo parámetro para determinar si el campo es numérico
      TextInputType
          keyboardType, // Nuevo parámetro para el tipo de teclado del campo de texto
      int?
          maxLength, // Nuevo parámetro para la longitud máxima del campo de texto
      int?
          minLines // Nuevo parámetro para el número mínimo de líneas en el campo de texto
      ) {
    return Row(
      mainAxisAlignment: MainAxisAlignment.start,
      crossAxisAlignment: CrossAxisAlignment.center,
      mainAxisSize: MainAxisSize.max,
      children: [
        Expanded(
          flex: 1,
          child: IconoYTextField(
            fontSize: textFieldFontSize,
            fontColor: textFieldFontColor,
            hintText: hintText,
            isPassword: isPasswordField,
            isNumeric: isNumericField,
            keyboardType: keyboardType,
            maxLength: maxLength,
            minLines: minLines,
          ),
        ),
      ],
    );
  }
}

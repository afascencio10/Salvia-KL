import 'package:flutter/material.dart';
import '../global/my_colors.dart';
import '../widgets/elements/espacio.dart';

class InfoScreen extends StatefulWidget {
  const InfoScreen({super.key});

  @override
  State<InfoScreen> createState() => _InfoScreenState();
}

class _InfoScreenState extends State<InfoScreen> {
  static const double _sectionHeaderFontSize = 20.0;
  static const int _headerTopFlex = 42;
  static const int _headerCenefaFlex = 6;
  static const int _headerBottomFlex = 46;

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
      extendBodyBehindAppBar: true,
      appBar: AppBar(
        title: const Text(
          "SALVIA",
          style: TextStyle(
            color: Color(0xffffffff),
            fontSize: 20,
          ),
        ),
        backgroundColor: AppColors.salviaTransparente,
        elevation: 0,
        automaticallyImplyLeading: true,
        iconTheme: const IconThemeData(
          color: Colors.white, // Color blanco para el botón back
        ),
      ),
      body: cuerpoInfo(context),
    );
  }

  Widget cuerpoInfo(BuildContext context) {
    final screenHeight = MediaQuery.of(context).size.height;
    final screenWidth = MediaQuery.of(context).size.width;
    final headerHeight = screenHeight * 0.3;
    const totalFlex = _headerTopFlex + _headerCenefaFlex + _headerBottomFlex;
    final panelTop =
        headerHeight * (_headerTopFlex + _headerCenefaFlex) / totalFlex;

    return Stack(
      children: [
        Container(
          decoration: const BoxDecoration(
            color: AppColors.salviaHeaderBottom,
          ),
        ),
        SizedBox(
          width: screenWidth,
          height: headerHeight,
          child: Column(
            children: [
              Expanded(
                flex: _headerTopFlex,
                child: Container(
                  width: double.infinity,
                  color: AppColors.salviaHeaderTopDark,
                ),
              ),
              Expanded(
                flex: _headerCenefaFlex,
                child: Container(
                  width: double.infinity,
                  color: AppColors.salviaHeaderBandSecondary,
                ),
              ),
              Expanded(
                flex: _headerBottomFlex,
                child: Container(
                  width: double.infinity,
                  color: AppColors.salviaHeaderBottom,
                ),
              ),
            ],
          ),
        ),
        SingleChildScrollView(
          child: ConstrainedBox(
            constraints: BoxConstraints(
              minHeight: screenHeight,
            ),
            child: Container(
              margin: EdgeInsets.fromLTRB(0, panelTop, 0, 0),
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
                  children: [
                    espacio(50),
                    _infoTitle("¿Qué es?"),
                    espacio(20),
                    _infoParagraph(
                        "Este sistema es una estrategia del Gobierno Nacional para la prevención del feminicidio y la eliminación de las Violencias Basadas en Género."),
                    espacio(30),
                    _infoTitle("¿A quién está dirigido?"),
                    espacio(20),
                    _infoBullet(
                        "A mujeres y personas LGBTIQ+, que vivan en Colombia en calidad de nacionales, residentes, migrantes regulares o en tránsito, que sean víctimas de violencias basadas en género."),
                    espacio(10),
                    _infoBullet("A familiares de víctimas de feminicidio."),
                    espacio(10),
                    _infoBullet(
                        "A hombres y personas responsables de violencias basadas en género."),
                    espacio(10),
                    _infoBullet(
                        "A servidores y servidoras públicas pertenecientes a las entidades responsables de atención a violencias basadas en género."),
                    espacio(30),
                    _infoTitle("Beneficios"),
                    espacio(20),
                    _infoParagraph(
                        "Marca a la Línea 155 Salvia para acceder a:"),
                    espacio(14),
                    _infoBullet(
                        "Orientación y activación de la ruta de atención integral"),
                    espacio(10),
                    _infoBullet(
                        "Acompañamiento permanente y diferencial a las mujeres y personas LGBTIQ+ víctimas de VBG para la eliminación de barreras de acceso a salud, justicia y medidas de atención, protección y estabilización."),
                    espacio(10),
                    _infoBullet(
                        "Atención psicosocial a víctimas de violencia basada en género"),
                    espacio(10),
                    _infoBullet(
                        "Medidas de emergencia para mujeres en riesgo de feminicidio y familiares de víctimas de feminicidio según valoración del caso."),
                    espacio(10),
                    _infoBullet(
                        "Acompañamiento a hombres y personas responsables de Violencias Basadas en Género."),
                    espacio(10),
                    _infoBullet(
                        "Información y datos cuantitativos organizados sobre las violencias basadas en género ocurridas a nivel nacional, regional y local."),
                    espacio(10),
                    _infoBullet(
                        "Articulación de organizaciones de la sociedad civil a través de la remisión de casos, generación de alertas y como red de apoyo para las víctimas de violencias."),
                    espacio(30),
                    _infoTitle("¿Cómo acceder?"),
                    espacio(20),
                    _infoBullet("Línea gratuita nacional 155 Salvia"),
                    espacio(10),
                    _infoBullet(
                        "Página web: salvia.minigualdadyequidad.gov.co"),
                    espacio(30),
                  ],
                ),
              ),
            ),
          ),
        ),
      ],
    );
  }

  Widget _infoTitle(String text) {
    return Text(
      text,
      textAlign: TextAlign.left,
      style: const TextStyle(
        fontSize: _sectionHeaderFontSize,
        fontWeight: FontWeight.w700,
        color: AppColors.salviaTituloPrincipal,
        letterSpacing: -0.1,
      ),
    );
  }

  Widget _infoParagraph(String text) {
    return Text(
      text,
      textAlign: TextAlign.left,
      style: const TextStyle(
        color: AppColors.salviaSubtitulos,
        letterSpacing: 0.6,
        fontWeight: FontWeight.normal,
        fontSize: 16.0,
      ),
    );
  }

  Widget _infoBullet(String text) {
    return Text(
      "• $text",
      textAlign: TextAlign.left,
      style: const TextStyle(
        color: AppColors.salviaSubtitulos,
        letterSpacing: 0.3,
        fontWeight: FontWeight.normal,
        fontSize: 16.0,
        height: 1.35,
      ),
    );
  }
}

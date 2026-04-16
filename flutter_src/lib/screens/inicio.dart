import 'package:flutter/material.dart';
import 'package:flutter/foundation.dart' show kIsWeb;
import 'package:flutter/services.dart';
import 'package:url_launcher/url_launcher.dart';
import '../global/my_screens.dart';
import '../global/my_colors.dart';
import '../widgets/elements/info_button_snack.dart';
import '../widgets/elements/go_button.dart';
import '../widgets/elements/espacio.dart';
import '../widgets/elements/title.dart';
import '../widgets/elements/subtitle.dart';
import '../widgets/elements/mi_link.dart';

class InicioScreen extends StatefulWidget {
  const InicioScreen({super.key});
  @override
  State<InicioScreen> createState() => _InicioScreenState();
}

class _InicioScreenState extends State<InicioScreen>
    with SingleTickerProviderStateMixin {
  static const double _actionInfoGap = 12;
  static const double _actionButtonHeight = 82;
  static const Color _reportarButtonColor = Color(0xFF380D44);
  static const Color _linea155ButtonColor = Color(0xFFC24877);
  static const Color _emergencia123ButtonColor = Color(0xFFEF875B);
  static const int _headerTopFlex = 42;
  static const int _headerCenefaFlex = 6;
  static const int _headerBottomFlex = 46;

  late final AnimationController _entryController;
  late final Animation<double> _contentFade;
  late final Animation<Offset> _reportarSlide;
  late final Animation<Offset> _emergenciaSlide;

  bool estadoVisibleGrupoA = false;
  bool estadoVisibleGrupoB = false;

  List<Widget> grupoA() {
    return [
      miTitulo("Show_1"),
    ];
  }

  List<Widget> grupoB() {
    return [
      miTitulo("Show_2"),
    ];
  }

  @override
  void initState() {
    super.initState();
    _entryController = AnimationController(
      vsync: this,
      duration: const Duration(seconds: 1),
    );
    final curve = CurvedAnimation(
      parent: _entryController,
      curve: Curves.fastOutSlowIn,
    );
    _contentFade = CurvedAnimation(
      parent: _entryController,
      curve: const Interval(0.2, 1.0, curve: Curves.easeIn),
    );
    _reportarSlide = Tween<Offset>(
      begin: const Offset(0, 0.14),
      end: Offset.zero,
    ).animate(curve);
    _emergenciaSlide = Tween<Offset>(
      begin: const Offset(0, -0.14),
      end: Offset.zero,
    ).animate(curve);
    _entryController.forward();

    estadoVisibleGrupoA = false;
    estadoVisibleGrupoB = false;
  }

  //llamada desde el botón de navegación a otras páginas
  void stopAudioIfNeeded() {
    // Audio temporalmente desactivado.
  }

  Future<void> _callEmergencyLine(String lineNumber) async {
    final uri = Uri(scheme: 'tel', path: lineNumber);
    try {
      final launched = await launchUrl(uri, mode: LaunchMode.platformDefault);
      if (!launched && mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
            content: Text(
              kIsWeb
                  ? 'En navegador de escritorio no siempre se puede abrir el marcador.'
                  : 'No se pudo iniciar la llamada al $lineNumber.',
            ),
          ),
        );
      }
    } on MissingPluginException {
      if (!mounted) return;
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(
          content: Text(
            'Plugin de llamada no disponible en esta sesión. Reinicia la app completa.',
          ),
        ),
      );
    } catch (_) {
      if (!mounted) return;
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(
          content:
              Text('No se pudo abrir el marcador para la línea $lineNumber.'),
        ),
      );
    }
  }

  void _restartEntryAnimation() {
    if (!mounted) return;
    _entryController
      ..stop()
      ..forward(from: 0);
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
    _entryController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: cuerpoInicio(context),
    );
  }

  Widget cuerpoInicio(BuildContext context) {
    final screenHeight = MediaQuery.of(context).size.height;
    final screenWidth = MediaQuery.of(context).size.width;
    final headerHeight = screenHeight * 0.3;
    // Mantiene el panel blanco pegado al bloque de cabecera para evitar
    // que aparezca una franja morada en pantallas bajas.
    final minPanelTop = screenHeight < 120.0 ? screenHeight : 120.0;
    final panelTop = (headerHeight - 28).clamp(minPanelTop, screenHeight);

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
                  color: AppColors.salviaHeaderTopInicio,
                  child: Center(
                    child: FractionallySizedBox(
                      heightFactor: 0.56,
                      widthFactor: 0.72,
                      child: Image.asset(
                        'lib/assets/images/igualdadwhite.png',
                        fit: BoxFit.contain,
                        errorBuilder: (_, __, ___) => const SizedBox.shrink(),
                      ),
                    ),
                  ),
                ),
              ),
              Expanded(
                flex: _headerCenefaFlex,
                child: Container(
                  width: double.infinity,
                  color: AppColors.salviaHeaderBandInicio,
                ),
              ),
              Expanded(
                flex: _headerBottomFlex,
                child: Container(
                  width: double.infinity,
                  color: AppColors.salviaHeaderBottom,
                  child: Center(
                    child: FractionallySizedBox(
                      heightFactor: 0.56,
                      widthFactor: 0.74,
                      child: Image.asset(
                        'lib/assets/images/salviaheader.png',
                        fit: BoxFit.contain,
                        errorBuilder: (_, __, ___) => const SizedBox.shrink(),
                      ),
                    ),
                  ),
                ),
              ),
            ],
          ),
        ),
        const Positioned(
          top: 16.0,
          left: 20.0,
          child: Opacity(
            opacity: 0.0,
            child: Icon(
              Icons.home_rounded,
              color: Colors.white,
              size: 30.0,
            ),
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
              width: screenWidth,
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
                    //-----------------------------------
                    //if (estadoVisibleGrupoA) ...grupoA(),
                    //if (estadoVisibleGrupoB) ...grupoB(),
                    //-----------------------------------
                    espacio(50), // Espacio después de label_A
                    FadeTransition(
                      opacity: _contentFade,
                      child: const Text(
                        "Bienvenida",
                        textAlign: TextAlign.center,
                        style: TextStyle(
                          fontSize: 34,
                          fontWeight: FontWeight.w400,
                          color: AppColors.salviaTituloPrincipal,
                          letterSpacing: -0.7,
                        ),
                      ),
                    ),
                    espacio(20), // Espacio después de label_A
                    FadeTransition(
                      opacity: _contentFade,
                      child: subtitulo("Elige una opción para poder ayudarte."),
                    ),
                    espacio(30),

                    SlideTransition(
                      position: _reportarSlide,
                      child: _buildActionRowWithInfoRight(
                        infoMessage:
                            "Elige esta opción si estás viviendo un hecho de violencia basada en género o si tienes conocimiento de un hecho de violencia experimentado por otra persona.",
                        infoSnackColor: _reportarButtonColor,
                        infoSnackIcon: Icons.chat_bubble_outline_rounded,
                        navButtonText: "Reportar",
                        onPressed: () async {
                          stopAudioIfNeeded();
                          await Navigator.of(context).push(
                            MaterialPageRoute(
                                builder: (context) => const ReportarScreen()),
                          );
                          _restartEntryAnimation();
                        },
                        navButtonColor: _reportarButtonColor,
                        navButtonTextColor: Colors.white,
                        navButtonIcon: Icons.chat_rounded,
                        navButtonIconColor: Colors.white,
                        navButtonIconSize: 24.0,
                        navButtonHeight: _actionButtonHeight,
                      ),
                    ),
                    const SizedBox(height: 14),
                    SlideTransition(
                      position: _emergenciaSlide,
                      child: Column(
                        children: [
                          _buildActionRowWithInfoRight(
                            infoMessage:
                                'Línea 155: Línea gratuita nacional del Sistema Nacional de Registro, Atención, Seguimiento y Monitoreo de Violencias Basadas en Género - Salvia.',
                            infoSnackColor: _linea155ButtonColor,
                            infoSnackIcon: Icons.woman_rounded,
                            navButtonText: 'Línea 155',
                            onPressed: () {
                              stopAudioIfNeeded();
                              _callEmergencyLine('155');
                            },
                            navButtonColor: _linea155ButtonColor,
                            navButtonTextColor: Colors.white,
                            navButtonIcon: Icons.call_rounded,
                            navButtonIconColor: Colors.white,
                            navButtonIconSize: 24.0,
                            navButtonHeight: _actionButtonHeight,
                          ),
                          const SizedBox(height: 14),
                          _buildActionRowWithInfoRight(
                            infoMessage:
                                'Línea 123: Línea de Emergencias de la Policía Nacional. Se debe marcar ante casos urgentes de riesgo inminente o violencia.',
                            infoSnackColor: _emergencia123ButtonColor,
                            infoSnackIcon: Icons.notifications_active_rounded,
                            navButtonText: 'Emergencia 123',
                            onPressed: () {
                              stopAudioIfNeeded();
                              _callEmergencyLine('123');
                            },
                            navButtonColor: _emergencia123ButtonColor,
                            navButtonTextColor: Colors.white,
                            navButtonIcon: Icons.call_rounded,
                            navButtonIconColor: Colors.white,
                            navButtonIconSize: 24.0,
                            navButtonHeight: _actionButtonHeight,
                          ),
                        ],
                      ),
                    ),
                    const SizedBox(height: 26),
                    FadeTransition(
                      opacity: _contentFade,
                      child: Padding(
                        padding: const EdgeInsets.symmetric(
                          horizontal: 24,
                        ),
                        child: LinkText(
                          text: "Conoce más sobre SALVIA",
                          color: AppColors.salviaLinkOutline,
                          fontSize: 13,
                          letterSpacing: -0.1,
                          asOutlinedButton: true,
                          outlinedBorderColor: AppColors.salviaLinkOutline,
                          outlinedBorderWidth: 1.2,
                          outlinedBackgroundColor: Colors.transparent,
                          borderRadius: 14,
                          padding: const EdgeInsets.symmetric(
                            vertical: 6,
                            horizontal: 14,
                          ),
                          onLinkTap: () {
                            Navigator.of(context).push(
                              MaterialPageRoute(
                                builder: (context) => const InfoScreen(),
                              ),
                            );
                          },
                        ),
                      ),
                    ),
                    espacio(18),
                  ],
                ),
              ),
            ),
          ),
        ),
      ],
    );
  }

  Widget _buildActionRowWithInfoRight({
    required String infoMessage,
    required Color infoSnackColor,
    required IconData infoSnackIcon,
    required String navButtonText,
    required VoidCallback onPressed,
    required Color navButtonColor,
    required Color navButtonTextColor,
    required IconData navButtonIcon,
    required Color navButtonIconColor,
    required double navButtonIconSize,
    required double navButtonHeight,
  }) {
    return Row(
      crossAxisAlignment: CrossAxisAlignment.center,
      children: [
        Expanded(
          child: NavigationButton(
            text: navButtonText,
            onPressed: onPressed,
            buttonColor: navButtonColor,
            textColor: navButtonTextColor,
            icon: navButtonIcon,
            iconColor: navButtonIconColor,
            iconSize: navButtonIconSize,
            buttonHeight: navButtonHeight,
            fontSize: 16,
            iconOnLeft: true,
            contentLeftAligned: true,
            leftContentPadding: 18,
          ),
        ),
        const SizedBox(width: _actionInfoGap),
        InfoButton(
          infoMessage: infoMessage,
          snackBackgroundColor: infoSnackColor,
          snackIcon: infoSnackIcon,
          boxedStyle: true,
          boxSize: navButtonHeight,
          borderRadius: 20,
          buttonBackgroundColor: const Color(0xFFF5F5F5),
          iconColor: AppColors.salviaPurpura3,
        ),
      ],
    );
  }
}

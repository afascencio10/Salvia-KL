import 'package:flutter/material.dart';

import '../global/my_colors.dart';
import '../global/my_screens.dart';
import '../widgets/elements/espacio.dart';
import '../widgets/elements/go_button.dart';
import '../widgets/elements/mi_link.dart';

class GraciasScreen extends StatefulWidget {
  const GraciasScreen({super.key});

  @override
  State<GraciasScreen> createState() => _GraciasScreenState();
}

class _GraciasScreenState extends State<GraciasScreen>
    with SingleTickerProviderStateMixin {
  static const double _actionSideSlot = 48;
  static const double _actionGap = 18;
  static const double _actionMirrorSpace = _actionSideSlot + _actionGap;
  static const int _headerTopFlex = 42;
  static const int _headerCenefaFlex = 6;
  static const int _headerBottomFlex = 46;

  late final AnimationController _entryController;
  late final Animation<double> _contentFade;
  late final Animation<Offset> _ctaSlide;

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
    _ctaSlide = Tween<Offset>(
      begin: const Offset(0, 0.12),
      end: Offset.zero,
    ).animate(curve);
    _entryController.forward();
  }

  @override
  void dispose() {
    _entryController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final screenHeight = MediaQuery.of(context).size.height;
    final screenWidth = MediaQuery.of(context).size.width;
    final headerHeight = screenHeight * 0.3;
    final minPanelTop = screenHeight < 120.0 ? screenHeight : 120.0;
    final panelTop = (headerHeight - 28).clamp(minPanelTop, screenHeight);

    return Scaffold(
      body: Stack(
        children: [
          const ColoredBox(
            color: AppColors.salviaHeaderBottom,
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
          SingleChildScrollView(
            child: ConstrainedBox(
              constraints: BoxConstraints(minHeight: screenHeight),
              child: Container(
                margin: EdgeInsets.fromLTRB(0, panelTop, 0, 0),
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
                  padding: const EdgeInsets.symmetric(horizontal: 20),
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.stretch,
                    children: [
                      espacio(60),
                      FadeTransition(
                        opacity: _contentFade,
                        child: const Text(
                          'Reporte recibido',
                          textAlign: TextAlign.center,
                          style: TextStyle(
                            fontSize: 34,
                            fontWeight: FontWeight.w400,
                            color: AppColors.salviaTituloPrincipal,
                            letterSpacing: -0.7,
                          ),
                        ),
                      ),
                      espacio(20),
                      FadeTransition(
                        opacity: _contentFade,
                        child: const Text(
                          'Hemos recibido tu reporte y ya inició su proceso de seguimiento. Si es necesario, nos comunicaremos contigo.',
                          textAlign: TextAlign.center,
                          style: TextStyle(
                            color: AppColors.salviaSubtitulos,
                            fontSize: 16,
                            fontWeight: FontWeight.w400,
                            height: 1.4,
                          ),
                        ),
                      ),
                      espacio(36),
                      SlideTransition(
                        position: _ctaSlide,
                        child: FadeTransition(
                          opacity: _contentFade,
                          child: Padding(
                            padding: const EdgeInsets.symmetric(
                              horizontal: _actionMirrorSpace,
                            ),
                            child: NavigationButton(
                              text: 'Ir al inicio',
                              onPressed: () {
                                Navigator.pushAndRemoveUntil(
                                  context,
                                  MaterialPageRoute(
                                    builder: (_) => const InicioScreen(),
                                  ),
                                  (Route<dynamic> route) => false,
                                );
                              },
                              buttonColor: AppColors.salviaPurpura1,
                              textColor: Colors.white,
                              icon: Icons.home_rounded,
                              iconColor: Colors.white,
                              iconSize: 24,
                              buttonHeight: 82,
                            ),
                          ),
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
                            text: 'Conoce más sobre SALVIA',
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
                      espacio(24),
                    ],
                  ),
                ),
              ),
            ),
          ),
        ],
      ),
    );
  }
}

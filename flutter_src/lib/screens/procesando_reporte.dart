import 'dart:async';
import 'dart:math' as math;

import 'package:flutter/material.dart';

import '../global/my_colors.dart';
import '../global/my_screens.dart';

class ProcesandoReporteScreen extends StatefulWidget {
  const ProcesandoReporteScreen({super.key});

  @override
  State<ProcesandoReporteScreen> createState() => _ProcesandoReporteScreenState();
}

class _ProcesandoReporteScreenState extends State<ProcesandoReporteScreen> {
  @override
  void initState() {
    super.initState();
    Timer(const Duration(seconds: 2), () {
      if (!mounted) return;
      Navigator.of(context).pushReplacement(
        MaterialPageRoute(builder: (_) => const GraciasScreen()),
      );
    });
  }

  @override
  Widget build(BuildContext context) {
    final screenHeight = MediaQuery.of(context).size.height;
    final screenWidth = MediaQuery.of(context).size.width;
    final headerHeight = screenHeight * 0.3;
    final panelTop = (headerHeight - 28).clamp(120.0, screenHeight);

    return Scaffold(
      body: Stack(
        children: [
          Container(
            decoration: const BoxDecoration(
              color: AppColors.salviaPurpura3,
            ),
          ),
          Container(
            width: screenWidth,
            height: headerHeight,
            decoration: const BoxDecoration(
              image: DecorationImage(
                image: AssetImage('lib/assets/images/header_app.jpg'),
                fit: BoxFit.cover,
                alignment: Alignment.topCenter,
              ),
            ),
          ),
          Container(
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
            child: const Align(
              alignment: Alignment.topCenter,
              child: Padding(
                padding: EdgeInsets.fromLTRB(24, 50, 24, 0),
                child: Column(
                  mainAxisSize: MainAxisSize.min,
                  children: [
                    Text(
                      'Un momento',
                      textAlign: TextAlign.center,
                      style: TextStyle(
                        fontSize: 34,
                        fontWeight: FontWeight.w400,
                        color: AppColors.salviaTituloPrincipal,
                        letterSpacing: -0.7,
                      ),
                    ),
                    SizedBox(height: 42),
                    _DotSpinner(
                      size: 110,
                      dotCount: 8,
                      dotSize: 12,
                      color: AppColors.salviaTituloPrincipal,
                    ),
                    SizedBox(height: 36),
                    Text(
                      'Analizando tu información',
                      textAlign: TextAlign.center,
                      style: TextStyle(
                        color: AppColors.salviaSubtitulos,
                        fontSize: 16,
                        fontWeight: FontWeight.w400,
                        height: 1.4,
                      ),
                    ),
                  ],
                ),
              ),
            ),
          ),
        ],
      ),
    );
  }
}

class _DotSpinner extends StatefulWidget {
  final double size;
  final int dotCount;
  final double dotSize;
  final Color color;

  const _DotSpinner({
    required this.size,
    required this.dotCount,
    required this.dotSize,
    required this.color,
  });

  @override
  State<_DotSpinner> createState() => _DotSpinnerState();
}

class _DotSpinnerState extends State<_DotSpinner>
    with SingleTickerProviderStateMixin {
  late final AnimationController _controller;

  @override
  void initState() {
    super.initState();
    _controller = AnimationController(
      vsync: this,
      duration: const Duration(milliseconds: 1150),
    )..repeat();
  }

  @override
  void dispose() {
    _controller.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final radius = (widget.size - widget.dotSize) / 2;
    return SizedBox(
      width: widget.size,
      height: widget.size,
      child: AnimatedBuilder(
        animation: _controller,
        builder: (context, _) {
          final turn = _controller.value * 2 * math.pi;
          return Stack(
            children: List.generate(widget.dotCount, (index) {
              final angle = turn + (index * 2 * math.pi / widget.dotCount);
              final x = radius + radius * math.cos(angle);
              final y = radius + radius * math.sin(angle);
              final phase =
                  ((index + _controller.value * widget.dotCount) % widget.dotCount) /
                      widget.dotCount;
              final opacity = 0.35 + (0.65 * phase);

              return Positioned(
                left: x,
                top: y,
                child: Opacity(
                  opacity: opacity.clamp(0.2, 1.0).toDouble(),
                  child: Container(
                    width: widget.dotSize,
                    height: widget.dotSize,
                    decoration: BoxDecoration(
                      color: widget.color,
                      shape: BoxShape.circle,
                    ),
                  ),
                ),
              );
            }),
          );
        },
      ),
    );
  }
}

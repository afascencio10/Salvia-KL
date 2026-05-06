import 'package:flutter/material.dart';
import 'package:demo_flut/global/my_screens.dart';

void main() {
  runApp(const SalviaApp());
}

class SalviaApp extends StatelessWidget {
  const SalviaApp({super.key});

  @override
  Widget build(BuildContext context) {
    final baseTextTheme = ThemeData.light().textTheme;

    return MaterialApp(
      debugShowCheckedModeBanner: false,
      title: 'Salvia',
      theme: ThemeData(
        primarySwatch: Colors.blue,
        fontFamily: 'NunitoSans',
        textTheme: baseTextTheme.copyWith(
          displayLarge:
              baseTextTheme.displayLarge?.copyWith(fontWeight: FontWeight.w700),
          displayMedium: baseTextTheme.displayMedium
              ?.copyWith(fontWeight: FontWeight.w700),
          displaySmall:
              baseTextTheme.displaySmall?.copyWith(fontWeight: FontWeight.w700),
          headlineLarge: baseTextTheme.headlineLarge
              ?.copyWith(fontWeight: FontWeight.w700),
          headlineMedium: baseTextTheme.headlineMedium
              ?.copyWith(fontWeight: FontWeight.w700),
          headlineSmall: baseTextTheme.headlineSmall
              ?.copyWith(fontWeight: FontWeight.w700),
          titleLarge:
              baseTextTheme.titleLarge?.copyWith(fontWeight: FontWeight.w700),
          titleMedium:
              baseTextTheme.titleMedium?.copyWith(fontWeight: FontWeight.w700),
          titleSmall:
              baseTextTheme.titleSmall?.copyWith(fontWeight: FontWeight.w700),
          bodyLarge:
              baseTextTheme.bodyLarge?.copyWith(fontWeight: FontWeight.w400),
          bodyMedium:
              baseTextTheme.bodyMedium?.copyWith(fontWeight: FontWeight.w400),
          bodySmall:
              baseTextTheme.bodySmall?.copyWith(fontWeight: FontWeight.w400),
          labelLarge:
              baseTextTheme.labelLarge?.copyWith(fontWeight: FontWeight.w400),
          labelMedium:
              baseTextTheme.labelMedium?.copyWith(fontWeight: FontWeight.w400),
          labelSmall:
              baseTextTheme.labelSmall?.copyWith(fontWeight: FontWeight.w400),
        ),
        appBarTheme: const AppBarTheme(
          titleTextStyle: TextStyle(
            fontFamily: 'NunitoSans',
            fontSize: 20,
            fontWeight: FontWeight.w700,
            color: Colors.white,
          ),
        ),
      ),
      home: const InicioScreen(),
    );
  }
}

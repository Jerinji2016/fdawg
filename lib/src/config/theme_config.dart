import 'package:flutter/material.dart';

const seedColor = Colors.purple;
const primaryBlack = Color(0xFF0D0D0C);

const hoverDuration = Duration(milliseconds: 200);

ThemeData get lightTheme => ThemeData.light(
      useMaterial3: true,
    ).copyWith(
      splashFactory: InkSparkle.splashFactory,
      scaffoldBackgroundColor: Colors.transparent,
      highlightColor: Colors.transparent,
      dividerColor: Colors.grey,
      cardColor: const Color(0x73FFFFFF),
      hoverColor: Colors.white54,
      textTheme: ThemeData.light(useMaterial3: true).textTheme.apply(
            fontFamily: 'Ubuntu Sans',
            displayColor: primaryBlack,
          ),
      inputDecorationTheme: InputDecorationTheme(
        border: OutlineInputBorder(
          borderRadius: BorderRadius.circular(16),
          borderSide: const BorderSide(
            color: Colors.grey,
            width: 2,
          ),
        ),
        hintStyle: const TextStyle(
          fontSize: 14,
          fontWeight: FontWeight.normal,
        ),
      ),
      colorScheme: ColorScheme.fromSeed(
        seedColor: seedColor,
      ),
    );

final darkTheme = ThemeData.dark(
  useMaterial3: true,
).copyWith(
  splashFactory: InkSparkle.splashFactory,
  scaffoldBackgroundColor: Colors.transparent,
  textTheme: ThemeData.dark(useMaterial3: true).textTheme.apply(
        fontFamily: 'Ubuntu Sans',
      ),
  inputDecorationTheme: InputDecorationTheme(
    border: OutlineInputBorder(
      borderRadius: BorderRadius.circular(16),
    ),
  ),
  cardColor: Colors.black26,
  colorScheme: const ColorScheme.dark(
    primary: Colors.purple,
  ),
);

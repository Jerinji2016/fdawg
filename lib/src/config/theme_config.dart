import 'package:flutter/material.dart';

final lightTheme = ThemeData.light(
  useMaterial3: true,
).copyWith(
  textTheme: ThemeData.light(useMaterial3: true).textTheme.apply(
        fontFamily: 'Ubuntu Sans',
      ),
  cardColor: Colors.black12,
  colorScheme: const ColorScheme.light(),
);

final darkTheme = ThemeData.dark(
  useMaterial3: true,
).copyWith(
  textTheme: ThemeData.dark(useMaterial3: true).textTheme.apply(
        fontFamily: 'Ubuntu Sans',
      ),
  cardColor: Colors.white10,
  colorScheme: const ColorScheme.dark(),
);

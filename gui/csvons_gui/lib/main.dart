import 'package:flutter/material.dart';

import 'screens/home_page.dart';

void main() {
  runApp(const CsvonsGuiApp());
}

class CsvonsGuiApp extends StatelessWidget {
  const CsvonsGuiApp({super.key});

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'csvons GUI',
      theme: ThemeData(useMaterial3: true, colorSchemeSeed: Colors.indigo),
      home: const HomePage(),
    );
  }
}

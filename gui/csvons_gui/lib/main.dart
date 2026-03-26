import 'package:flutter/material.dart';

import 'screens/home_page.dart';
import 'screens/workspace_page.dart';

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
      home: const _RootShell(),
    );
  }
}

class _RootShell extends StatefulWidget {
  const _RootShell();

  @override
  State<_RootShell> createState() => _RootShellState();
}

class _RootShellState extends State<_RootShell> {
  int _selectedIndex = 0;

  static const _pages = <Widget>[
    HomePage(),
    WorkspacePage(),
  ];

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: _pages[_selectedIndex],
      bottomNavigationBar: NavigationBar(
        selectedIndex: _selectedIndex,
        destinations: const [
          NavigationDestination(
            icon: Icon(Icons.playlist_play),
            label: 'Validate',
          ),
          NavigationDestination(
            icon: Icon(Icons.folder_open),
            label: 'Workspace',
          ),
        ],
        onDestinationSelected: (index) {
          setState(() => _selectedIndex = index);
        },
      ),
    );
  }
}

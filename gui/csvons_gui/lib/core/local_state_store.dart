import 'dart:convert';
import 'dart:io';

class LocalState {
  final List<String> recentBinaryPaths;
  final List<String> recentRulerPaths;
  final List<String> recentWorkspacePaths;

  const LocalState({
    required this.recentBinaryPaths,
    required this.recentRulerPaths,
    required this.recentWorkspacePaths,
  });

  factory LocalState.empty() => const LocalState(
        recentBinaryPaths: <String>[],
        recentRulerPaths: <String>[],
        recentWorkspacePaths: <String>[],
      );

  Map<String, dynamic> toJson() => {
        'recent_binary_paths': recentBinaryPaths,
        'recent_ruler_paths': recentRulerPaths,
        'recent_workspace_paths': recentWorkspacePaths,
      };

  factory LocalState.fromJson(Map<String, dynamic> json) {
    return LocalState(
      recentBinaryPaths:
          (json['recent_binary_paths'] as List<dynamic>? ?? const <dynamic>[])
              .map((e) => e.toString())
              .toList(growable: false),
      recentRulerPaths:
          (json['recent_ruler_paths'] as List<dynamic>? ?? const <dynamic>[])
              .map((e) => e.toString())
              .toList(growable: false),
      recentWorkspacePaths:
          (json['recent_workspace_paths'] as List<dynamic>? ?? const <dynamic>[])
              .map((e) => e.toString())
              .toList(growable: false),
    );
  }
}

class LocalStateStore {
  LocalStateStore({String? fileName})
      : _fileName = fileName ?? '.csvons_gui_state.json';

  final String _fileName;

  Future<LocalState> load() async {
    final file = File(_fileName);
    if (!await file.exists()) {
      return LocalState.empty();
    }

    final raw = await file.readAsString();
    if (raw.trim().isEmpty) {
      return LocalState.empty();
    }

    final jsonObj = jsonDecode(raw) as Map<String, dynamic>;
    return LocalState.fromJson(jsonObj);
  }

  Future<void> saveRecentPaths({
    required String binaryPath,
    required String rulerPath,
    String? workspacePath,
    int maxItems = 8,
  }) async {
    final existing = await load();

    List<String> pushTop(List<String> input, String value) {
      final trimmed = value.trim();
      if (trimmed.isEmpty) return input;

      final next = <String>[trimmed, ...input.where((v) => v != trimmed)];
      if (next.length > maxItems) {
        return next.sublist(0, maxItems);
      }
      return next;
    }

    final next = LocalState(
      recentBinaryPaths: pushTop(existing.recentBinaryPaths, binaryPath),
      recentRulerPaths: pushTop(existing.recentRulerPaths, rulerPath),
      recentWorkspacePaths: pushTop(
        existing.recentWorkspacePaths,
        workspacePath ?? '',
      ),
    );

    await File(_fileName).writeAsString(
      const JsonEncoder.withIndent('  ').convert(next.toJson()),
    );
  }

  Future<void> saveRecentWorkspace({
    required String workspacePath,
    int maxItems = 8,
  }) async {
    final existing = await load();
    final trimmed = workspacePath.trim();
    if (trimmed.isEmpty) return;

    final next = <String>[
      trimmed,
      ...existing.recentWorkspacePaths.where((v) => v != trimmed),
    ];

    final state = LocalState(
      recentBinaryPaths: existing.recentBinaryPaths,
      recentRulerPaths: existing.recentRulerPaths,
      recentWorkspacePaths: next.take(maxItems).toList(growable: false),
    );

    await File(_fileName).writeAsString(
      const JsonEncoder.withIndent('  ').convert(state.toJson()),
    );
  }
}

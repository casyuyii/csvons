import 'dart:convert';
import 'dart:io';

class LocalState {
  final List<String> recentBinaryPaths;
  final List<String> recentRulerPaths;

  const LocalState({
    required this.recentBinaryPaths,
    required this.recentRulerPaths,
  });

  factory LocalState.empty() => const LocalState(
        recentBinaryPaths: <String>[],
        recentRulerPaths: <String>[],
      );

  Map<String, dynamic> toJson() => {
        'recent_binary_paths': recentBinaryPaths,
        'recent_ruler_paths': recentRulerPaths,
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
    );

    await File(_fileName).writeAsString(
      const JsonEncoder.withIndent('  ').convert(next.toJson()),
    );
  }
}

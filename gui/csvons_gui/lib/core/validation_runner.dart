import 'dart:convert';
import 'dart:io';

import '../models/validation_report.dart';

typedef ProcessStarter = Future<Process> Function(
  String executable,
  List<String> arguments,
);

class ValidationRunner {
  ValidationRunner({
    required this.binaryPath,
    ProcessStarter? processStarter,
  }) : _processStarter = processStarter ?? Process.start;

  final String binaryPath;
  final ProcessStarter _processStarter;

  Future<ValidationResult> run({
    required String rulerPath,
    bool jsonFormat = true,
  }) async {
    final executable = resolveBinaryPath(binaryPath);
    final args = <String>[
      if (jsonFormat) ...['--format', 'json'],
      rulerPath,
    ];

    final proc = await _processStarter(executable, args);
    final outFuture = proc.stdout.transform(utf8.decoder).join();
    final errFuture = proc.stderr.transform(utf8.decoder).join();

    final exitCode = await proc.exitCode;
    final out = await outFuture;
    final err = await errFuture;

    ValidationReport? report;
    if (jsonFormat && out.trim().isNotEmpty) {
      try {
        report = ValidationReport.fromJson(
          jsonDecode(out) as Map<String, dynamic>,
        );
      } catch (_) {
        // Keep report null when stdout is not valid JSON yet.
      }
    }

    return ValidationResult(
      exitCode: exitCode,
      report: report,
      stdoutText: out,
      stderrText: err,
    );
  }

  static String defaultBinaryPath() {
    for (final candidate in _existingCandidatePaths()) {
      return candidate;
    }
    return devBinaryPath();
  }

  static String resolveBinaryPath(String explicitPath) {
    final trimmed = explicitPath.trim();
    if (trimmed.isNotEmpty) return trimmed;

    for (final candidate in _existingCandidatePaths()) {
      return candidate;
    }
    return devBinaryPath();
  }

  static String bundledBinaryPath() {
    final parts = <String>['bin', _platformDirectoryName(), _binaryName()];
    return parts.join(Platform.pathSeparator);
  }

  static String devBinaryPath() {
    if (Platform.isWindows) return 'bin/csvons.exe';
    if (Platform.isMacOS) return 'bin/csvons_macos';
    if (Platform.isLinux) return 'bin/csvons_linux';
    throw UnsupportedError('Unsupported platform for csvons binary');
  }

  static List<String> _existingCandidatePaths() {
    final candidates = <String>[
      _bundledBinaryPathNextToExecutable(),
      bundledBinaryPath(),
      devBinaryPath(),
    ];

    final seen = <String>{};
    final existing = <String>[];
    for (final candidate in candidates) {
      final normalized = candidate.trim();
      if (normalized.isEmpty || !seen.add(normalized)) continue;
      if (File(normalized).existsSync()) {
        existing.add(normalized);
      }
    }
    return existing;
  }

  static String _bundledBinaryPathNextToExecutable() {
    final executableDir = File(Platform.resolvedExecutable).parent.path;
    final parts = <String>[
      executableDir,
      'bin',
      _platformDirectoryName(),
      _binaryName(),
    ];
    return parts.join(Platform.pathSeparator);
  }

  static String _platformDirectoryName() {
    if (Platform.isWindows) return 'windows';
    if (Platform.isMacOS) return 'macos';
    if (Platform.isLinux) return 'linux';
    throw UnsupportedError('Unsupported platform for csvons binary');
  }

  static String _binaryName() {
    if (Platform.isWindows) return 'csvons.exe';
    return 'csvons';
  }
}

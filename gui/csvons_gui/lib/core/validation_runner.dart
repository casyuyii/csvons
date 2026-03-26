import 'dart:convert';
import 'dart:io';

import '../models/validation_report.dart';

class ValidationRunner {
  ValidationRunner({required this.binaryPath});

  final String binaryPath;

  Future<ValidationResult> run({
    required String rulerPath,
    bool jsonFormat = true,
  }) async {
    final args = <String>[
      if (jsonFormat) ...['--format', 'json', '--quiet'],
      rulerPath,
    ];

    final proc = await Process.start(binaryPath, args);
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
    if (Platform.isWindows) return 'bin/csvons.exe';
    if (Platform.isMacOS) return 'bin/csvons_macos';
    if (Platform.isLinux) return 'bin/csvons_linux';
    throw UnsupportedError('Unsupported platform for csvons binary');
  }
}

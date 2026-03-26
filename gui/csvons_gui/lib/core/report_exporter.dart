import 'dart:convert';
import 'dart:io';

import '../models/validation_report.dart';

class ReportExporter {
  static String toPrettyJson(ValidationReport report) {
    final map = <String, dynamic>{
      'schema_version': report.schemaVersion,
      'summary': {
        'files_checked': report.summary.filesChecked,
        'passed': report.summary.passed,
        'failed': report.summary.failed,
        'duration_ms': report.summary.durationMs,
      },
      'issues': report.issues
          .map(
            (issue) => {
              'file': issue.file,
              'rule': issue.rule,
              'field': issue.field,
              'row': issue.row,
              'value': issue.value,
              'message': issue.message,
              'severity': issue.severity,
            },
          )
          .toList(growable: false),
    };

    return const JsonEncoder.withIndent('  ').convert(map);
  }

  static String toMarkdown(ValidationReport report) {
    final buffer = StringBuffer()
      ..writeln('# csvons Validation Report')
      ..writeln()
      ..writeln('- Schema: `${report.schemaVersion}`')
      ..writeln('- Files checked: ${report.summary.filesChecked}')
      ..writeln('- Passed: ${report.summary.passed}')
      ..writeln('- Failed: ${report.summary.failed}')
      ..writeln('- Duration: ${report.summary.durationMs} ms')
      ..writeln();

    if (report.issues.isEmpty) {
      buffer.writeln('## Issues');
      buffer.writeln();
      buffer.writeln('No issues found.');
      return buffer.toString();
    }

    buffer
      ..writeln('## Issues')
      ..writeln()
      ..writeln('| Severity | File | Rule | Field | Row | Value | Message |')
      ..writeln('|---|---|---|---|---:|---|---|');

    for (final issue in report.issues) {
      buffer.writeln(
        '| ${_md(issue.severity)} | ${_md(issue.file)} | ${_md(issue.rule)} | '
        '${_md(issue.field)} | ${issue.row ?? ''} | ${_md(issue.value ?? '')} | '
        '${_md(issue.message)} |',
      );
    }

    return buffer.toString();
  }

  static Future<String> exportJson({
    required ValidationReport report,
    required String path,
  }) async {
    final normalized = _normalizePath(path, '.json');
    await _ensureParentDirectory(normalized);
    await File(normalized).writeAsString(toPrettyJson(report));
    return normalized;
  }

  static Future<String> exportMarkdown({
    required ValidationReport report,
    required String path,
  }) async {
    final normalized = _normalizePath(path, '.md');
    await _ensureParentDirectory(normalized);
    await File(normalized).writeAsString(toMarkdown(report));
    return normalized;
  }

  static Future<void> _ensureParentDirectory(String path) async {
    final parent = File(path).parent;
    if (await parent.exists()) return;
    await parent.create(recursive: true);
  }

  static String _normalizePath(String path, String extension) {
    final trimmed = path.trim();
    if (trimmed.isEmpty) {
      throw const FileSystemException('Export path is required');
    }
    return trimmed.toLowerCase().endsWith(extension) ? trimmed : '$trimmed$extension';
  }

  static String _md(String input) {
    return input.replaceAll('|', '\\|').replaceAll('\n', ' ');
  }
}

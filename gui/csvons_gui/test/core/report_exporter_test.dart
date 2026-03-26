import 'dart:io';

import 'package:csvons_gui/core/report_exporter.dart';
import 'package:csvons_gui/models/validation_report.dart';
import 'package:test/test.dart';

void main() {
  ValidationReport sampleReport() {
    return ValidationReport(
      schemaVersion: '1',
      summary: ValidationSummary(
        filesChecked: 2,
        passed: 1,
        failed: 1,
        durationMs: 12,
      ),
      issues: [
        ValidationIssue(
          file: 'orders.csv',
          rule: 'exists',
          field: 'user_id',
          row: 3,
          value: 'u404',
          message: 'missing | ref',
          severity: 'error',
        ),
      ],
    );
  }

  test('toMarkdown renders escaped table', () {
    final md = ReportExporter.toMarkdown(sampleReport());

    expect(md, contains('# csvons Validation Report'));
    expect(md, contains('| Severity | File | Rule | Field | Row | Value | Message |'));
    expect(md, contains('missing \\| ref'));
  });

  test('exportJson appends extension and creates missing parent folders', () async {
    final dir = await Directory.systemTemp.createTemp('report_exporter');
    final nested = '${dir.path}${Platform.pathSeparator}out${Platform.pathSeparator}daily${Platform.pathSeparator}report';

    final written = await ReportExporter.exportJson(
      report: sampleReport(),
      path: nested,
    );

    expect(written, endsWith('.json'));
    expect(await File(written).exists(), isTrue);

    await dir.delete(recursive: true);
  });
}

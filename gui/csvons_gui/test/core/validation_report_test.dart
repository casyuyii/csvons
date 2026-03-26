import 'package:csvons_gui/models/validation_report.dart';
import 'package:test/test.dart';

void main() {
  test('parses report with defaults for missing fields', () {
    final report = ValidationReport.fromJson({
      'schema_version': '1',
      'summary': {
        'files_checked': 2,
      },
      'issues': [
        {
          'file': 'orders.csv',
          'rule': 'exists',
          'field': 'user_id',
          'row': 11,
          'message': 'missing user',
        },
      ],
    });

    expect(report.schemaVersion, '1');
    expect(report.summary.filesChecked, 2);
    expect(report.summary.passed, 0);
    expect(report.summary.failed, 0);
    expect(report.summary.durationMs, 0);
    expect(report.issues, hasLength(1));
    expect(report.issues.first.severity, 'error');
    expect(report.issues.first.row, 11);
  });
}

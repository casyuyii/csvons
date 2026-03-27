import 'package:csvons_gui/models/validation_report.dart';
import 'package:csvons_gui/widgets/issues_table.dart';
import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';

void main() {
  ValidationReport buildReport() {
    return ValidationReport(
      schemaVersion: '1',
      summary: ValidationSummary(
        filesChecked: 1,
        passed: 0,
        failed: 1,
        durationMs: 12,
      ),
      issues: <ValidationIssue>[
        ValidationIssue(
          file: 'orders.csv',
          rule: 'exists',
          field: 'user_id',
          row: 10,
          value: 'u404',
          message: 'missing reference',
          severity: 'error',
        ),
        ValidationIssue(
          file: 'orders.csv',
          rule: 'vtype',
          field: 'amount',
          row: 7,
          value: 'abc',
          message: 'invalid float',
          severity: 'warning',
        ),
      ],
    );
  }

  Future<void> pumpTable(WidgetTester tester) async {
    await tester.pumpWidget(
      MaterialApp(
        home: Scaffold(
          body: IssuesTable(report: buildReport()),
        ),
      ),
    );
    await tester.pumpAndSettle();
  }

  Future<void> pumpCustomTable({
    required WidgetTester tester,
    required ValidationReport report,
  }) async {
    await tester.pumpWidget(
      MaterialApp(
        home: Scaffold(
          body: IssuesTable(report: report),
        ),
      ),
    );
    await tester.pumpAndSettle();
  }

  String firstValueCell(WidgetTester tester) {
    final table = tester.widget<DataTable>(find.byType(DataTable));
    final firstRow = table.rows.first;
    final valueCell = firstRow.cells[5].child as Text;
    return valueCell.data ?? '';
  }

  testWidgets('shows severity chips and filters with chips', (tester) async {
    await pumpTable(tester);

    expect(find.text('Showing 2 of 2 issue(s)'), findsOneWidget);
    expect(find.text('All (2)'), findsOneWidget);
    expect(find.text('Error (1)'), findsOneWidget);
    expect(find.text('Warning (1)'), findsOneWidget);
    expect(find.text('Value'), findsOneWidget);

    await tester.tap(find.text('Error (1)'));
    await tester.pumpAndSettle();

    expect(find.text('Showing 1 of 2 issue(s)'), findsOneWidget);
    expect(find.text('missing reference'), findsOneWidget);
    expect(find.text('invalid float'), findsNothing);

    final resetButton = tester.widget<TextButton>(find.widgetWithText(TextButton, 'Reset filters'));
    expect(resetButton.onPressed, isNotNull);
  });

  testWidgets('shows empty-state message when search has no matches', (tester) async {
    await pumpTable(tester);

    await tester.enterText(
      find.widgetWithText(TextField, 'Search message/file/rule/field/value/row/severity'),
      'no-match-value',
    );
    await tester.pumpAndSettle();

    expect(find.text('Showing 0 of 2 issue(s)'), findsOneWidget);
    expect(find.text('No issues match the current filters.'), findsOneWidget);
    expect(find.text('Reset all filters'), findsOneWidget);

    await tester.tap(find.text('Reset all filters'));
    await tester.pumpAndSettle();

    expect(find.text('Showing 2 of 2 issue(s)'), findsOneWidget);
    expect(find.text('No issues match the current filters.'), findsNothing);
  });

  testWidgets('filters by file/rule and resets filters', (tester) async {
    await pumpTable(tester);
    final initialResetButton = tester.widget<TextButton>(
      find.widgetWithText(TextButton, 'Reset filters'),
    );
    expect(initialResetButton.onPressed, isNull);

    await tester.tap(find.text('All files'));
    await tester.pumpAndSettle();
    await tester.tap(find.text('orders.csv').last);
    await tester.pumpAndSettle();

    await tester.tap(find.text('All rules'));
    await tester.pumpAndSettle();
    await tester.tap(find.text('exists').last);
    await tester.pumpAndSettle();

    expect(find.text('Showing 1 of 2 issue(s)'), findsOneWidget);
    expect(find.text('All (1)'), findsOneWidget);
    expect(find.text('Error (1)'), findsOneWidget);
    expect(find.text('Warning (0)'), findsOneWidget);
    expect(find.text('(filtered)'), findsOneWidget);
    expect(
      find.text('Active filters: visible=1/2, file=orders.csv, rule=exists'),
      findsOneWidget,
    );
    expect(find.byTooltip('Copy active filter summary'), findsOneWidget);
    await tester.tap(find.byTooltip('Copy active filter summary'));
    await tester.pump();
    expect(find.text('Active filter summary copied'), findsOneWidget);

    await tester.tap(find.text('Reset filters'));
    await tester.pumpAndSettle();

    expect(find.text('Showing 2 of 2 issue(s)'), findsOneWidget);
    expect(find.text('(filtered)'), findsNothing);
    expect(find.byTooltip('Copy active filter summary'), findsNothing);
  });

  testWidgets('formats custom severity labels in chips', (tester) async {
    final report = ValidationReport(
      schemaVersion: '1',
      summary: ValidationSummary(
        filesChecked: 1,
        passed: 0,
        failed: 1,
        durationMs: 1,
      ),
      issues: <ValidationIssue>[
        ValidationIssue(
          file: 'audit.csv',
          rule: 'vtype',
          field: 'status',
          row: 1,
          value: 'unknown',
          message: 'needs review',
          severity: 'needs_review',
        ),
      ],
    );

    await pumpCustomTable(tester: tester, report: report);

    expect(find.text('Needs Review (1)'), findsOneWidget);
  });

  testWidgets('shows explicit empty-report message when no issues exist', (tester) async {
    final report = ValidationReport(
      schemaVersion: '1',
      summary: ValidationSummary(
        filesChecked: 1,
        passed: 1,
        failed: 0,
        durationMs: 1,
      ),
      issues: const <ValidationIssue>[],
    );

    await pumpCustomTable(tester: tester, report: report);

    expect(find.text('No issues reported for this run.'), findsOneWidget);
    expect(find.textContaining('Showing'), findsNothing);
  });

  testWidgets('shows active filter summary for search and severity', (tester) async {
    await pumpTable(tester);

    await tester.enterText(
      find.widgetWithText(TextField, 'Search message/file/rule/field/value/row/severity'),
      'invalid',
    );
    await tester.pumpAndSettle();
    await tester.tap(find.text('Warning (1)'));
    await tester.pumpAndSettle();

    expect(
      find.text('Active filters: visible=1/2, severity=Warning, search=\"invalid\"'),
      findsOneWidget,
    );
  });

  testWidgets('orders severity chips by semantic priority', (tester) async {
    final report = ValidationReport(
      schemaVersion: '1',
      summary: ValidationSummary(
        filesChecked: 1,
        passed: 0,
        failed: 1,
        durationMs: 1,
      ),
      issues: <ValidationIssue>[
        ValidationIssue(
          file: 'a.csv',
          rule: 'vtype',
          field: 'x',
          row: 1,
          value: 'bad',
          message: 'm1',
          severity: 'info',
        ),
        ValidationIssue(
          file: 'a.csv',
          rule: 'vtype',
          field: 'x',
          row: 2,
          value: 'bad',
          message: 'm2',
          severity: 'warning',
        ),
        ValidationIssue(
          file: 'a.csv',
          rule: 'vtype',
          field: 'x',
          row: 3,
          value: 'bad',
          message: 'm3',
          severity: 'error',
        ),
        ValidationIssue(
          file: 'a.csv',
          rule: 'vtype',
          field: 'x',
          row: 4,
          value: 'bad',
          message: 'm4',
          severity: 'critical',
        ),
      ],
    );

    await pumpCustomTable(tester: tester, report: report);

    final chips = tester.widgetList<FilterChip>(find.byType(FilterChip)).toList();
    final labels = chips
        .map(
          (chip) => ((chip.label as Text).data ?? '').trim(),
        )
        .toList();

    expect(labels, ['All (4)', 'Critical (1)', 'Error (1)', 'Warning (1)', 'Info (1)']);
  });

  testWidgets('sorts by value column when header is tapped', (tester) async {
    await pumpTable(tester);

    expect(firstValueCell(tester), 'u404');

    await tester.tap(find.text('Value'));
    await tester.pumpAndSettle();

    expect(firstValueCell(tester), 'abc');
  });
}

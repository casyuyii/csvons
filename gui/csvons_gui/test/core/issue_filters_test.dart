import 'package:csvons_gui/core/issue_filters.dart';
import 'package:csvons_gui/models/validation_report.dart';
import 'package:test/test.dart';

void main() {
  List<ValidationIssue> seedIssues() {
    return [
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
      ValidationIssue(
        file: 'users.csv',
        rule: 'unique',
        field: 'email',
        row: 2,
        value: 'dupe@example.com',
        message: 'duplicate',
        severity: 'error',
      ),
    ];
  }

  test('filters by severity and search query', () {
    final result = filterAndSortIssues(
      issues: seedIssues(),
      query: 'orders',
      severityFilter: 'error',
      fileFilter: 'all',
      ruleFilter: 'all',
      sortField: IssueSortField.file,
      ascending: true,
    );

    expect(result, hasLength(1));
    expect(result.first.rule, 'exists');
  });

  test('search query matches issue value field', () {
    final result = filterAndSortIssues(
      issues: seedIssues(),
      query: 'u404',
      sortField: IssueSortField.file,
      ascending: true,
    );

    expect(result, hasLength(1));
    expect(result.first.value, 'u404');
  });

  test('search query matches issue row and severity text', () {
    final byRow = filterAndSortIssues(
      issues: seedIssues(),
      query: '10',
      sortField: IssueSortField.file,
      ascending: true,
    );
    expect(byRow, hasLength(1));
    expect(byRow.first.row, 10);

    final bySeverity = filterAndSortIssues(
      issues: seedIssues(),
      query: 'warning',
      sortField: IssueSortField.file,
      ascending: true,
    );
    expect(bySeverity, hasLength(1));
    expect(bySeverity.first.severity, 'warning');
  });

  test('sorts by row descending', () {
    final result = filterAndSortIssues(
      issues: seedIssues(),
      query: '',
      severityFilter: 'all',
      fileFilter: 'all',
      ruleFilter: 'all',
      sortField: IssueSortField.row,
      ascending: false,
    );

    expect(result.map((i) => i.row), [10, 7, 2]);
  });

  test('sorts null rows last when sorting row ascending', () {
    final issues = <ValidationIssue>[
      ValidationIssue(
        file: 'orders.csv',
        rule: 'exists',
        field: 'user_id',
        row: null,
        value: 'u404',
        message: 'missing reference',
        severity: 'error',
      ),
      ValidationIssue(
        file: 'orders.csv',
        rule: 'exists',
        field: 'user_id',
        row: 3,
        value: 'u500',
        message: 'bad row',
        severity: 'error',
      ),
    ];

    final result = filterAndSortIssues(
      issues: issues,
      sortField: IssueSortField.row,
      ascending: true,
    );

    expect(result.map((i) => i.row).toList(), [3, null]);
  });

  test('filters by file and rule selectors', () {
    final result = filterAndSortIssues(
      issues: seedIssues(),
      query: '',
      severityFilter: 'all',
      fileFilter: 'orders.csv',
      ruleFilter: 'vtype',
      sortField: IssueSortField.row,
      ascending: true,
    );

    expect(result, hasLength(1));
    expect(result.first.message, 'invalid float');
  });

  test('severity filtering is case-insensitive', () {
    final result = filterAndSortIssues(
      issues: seedIssues(),
      query: '',
      severityFilter: 'ERROR',
      fileFilter: 'all',
      ruleFilter: 'all',
      sortField: IssueSortField.file,
      ascending: true,
    );

    expect(result, hasLength(2));
    expect(result.every((issue) => issue.severity == 'error'), isTrue);
  });

  test('supports filtering custom severity labels', () {
    final issues = <ValidationIssue>[
      ...seedIssues(),
      ValidationIssue(
        file: 'users.csv',
        rule: 'vtype',
        field: 'status',
        row: 3,
        value: 'unknown',
        message: 'needs review',
        severity: 'info',
      ),
    ];
    final result = filterAndSortIssues(
      issues: issues,
      query: '',
      severityFilter: 'INFO',
      fileFilter: 'all',
      ruleFilter: 'all',
      sortField: IssueSortField.row,
      ascending: true,
    );

    expect(result, hasLength(1));
    expect(result.first.severity, 'info');
  });

  test('uses sensible defaults when optional filters are omitted', () {
    final result = filterAndSortIssues(issues: seedIssues());

    expect(result, hasLength(3));
    expect(result.first.severity, 'error');
  });

  test('sorts severity with semantic rank order', () {
    final issues = <ValidationIssue>[
      ...seedIssues(),
      ValidationIssue(
        file: 'meta.csv',
        rule: 'vtype',
        field: 'status',
        row: 0,
        value: 'panic',
        message: 'critical issue',
        severity: 'critical',
      ),
      ValidationIssue(
        file: 'meta.csv',
        rule: 'vtype',
        field: 'status',
        row: 1,
        value: 'ready',
        message: 'informational note',
        severity: 'info',
      ),
    ];

    final asc = filterAndSortIssues(
      issues: issues,
      sortField: IssueSortField.severity,
      ascending: true,
    );
    expect(asc.map((i) => i.severity).toList(), [
      'critical',
      'error',
      'error',
      'warning',
      'info',
    ]);

    final desc = filterAndSortIssues(
      issues: issues,
      sortField: IssueSortField.severity,
      ascending: false,
    );
    expect(desc.map((i) => i.severity).toList(), [
      'info',
      'warning',
      'error',
      'error',
      'critical',
    ]);
  });

  test('applies deterministic tie-breakers for equal sort keys', () {
    final issues = <ValidationIssue>[
      ValidationIssue(
        file: 'orders.csv',
        rule: 'exists',
        field: 'user_id',
        row: 12,
        value: 'u404',
        message: 'z-last',
        severity: 'error',
      ),
      ValidationIssue(
        file: 'orders.csv',
        rule: 'exists',
        field: 'user_id',
        row: 4,
        value: 'u500',
        message: 'a-first',
        severity: 'error',
      ),
    ];

    final sorted = filterAndSortIssues(
      issues: issues,
      sortField: IssueSortField.severity,
      ascending: true,
    );

    expect(sorted.map((issue) => issue.row).toList(), [4, 12]);
    expect(sorted.map((issue) => issue.message).toList(), ['a-first', 'z-last']);
  });

  test('sorts by value field', () {
    final issues = <ValidationIssue>[
      ValidationIssue(
        file: 'orders.csv',
        rule: 'exists',
        field: 'user_id',
        row: 2,
        value: 'z-last',
        message: 'm1',
        severity: 'error',
      ),
      ValidationIssue(
        file: 'orders.csv',
        rule: 'exists',
        field: 'user_id',
        row: 1,
        value: 'a-first',
        message: 'm2',
        severity: 'error',
      ),
    ];

    final sorted = filterAndSortIssues(
      issues: issues,
      sortField: IssueSortField.value,
      ascending: true,
    );

    expect(sorted.map((issue) => issue.value).toList(), ['a-first', 'z-last']);
  });

  test('sorts null/empty values last when sorting by value ascending', () {
    final issues = <ValidationIssue>[
      ValidationIssue(
        file: 'orders.csv',
        rule: 'exists',
        field: 'user_id',
        row: 2,
        value: null,
        message: 'missing value',
        severity: 'error',
      ),
      ValidationIssue(
        file: 'orders.csv',
        rule: 'exists',
        field: 'user_id',
        row: 1,
        value: '',
        message: 'empty value',
        severity: 'error',
      ),
      ValidationIssue(
        file: 'orders.csv',
        rule: 'exists',
        field: 'user_id',
        row: 3,
        value: 'abc',
        message: 'present value',
        severity: 'error',
      ),
    ];

    final sorted = filterAndSortIssues(
      issues: issues,
      sortField: IssueSortField.value,
      ascending: true,
    );

    expect(sorted.map((issue) => issue.value).toList(), ['abc', '', null]);
  });
}

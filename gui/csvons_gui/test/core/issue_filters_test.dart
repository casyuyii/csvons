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
      sortField: IssueSortField.file,
      ascending: true,
    );

    expect(result, hasLength(1));
    expect(result.first.rule, 'exists');
  });

  test('sorts by row descending', () {
    final result = filterAndSortIssues(
      issues: seedIssues(),
      query: '',
      severityFilter: 'all',
      sortField: IssueSortField.row,
      ascending: false,
    );

    expect(result.map((i) => i.row), [10, 7, 2]);
  });
}

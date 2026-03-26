import '../models/validation_report.dart';

enum IssueSortField {
  severity,
  file,
  rule,
  field,
  row,
  message,
}

List<ValidationIssue> filterAndSortIssues({
  required List<ValidationIssue> issues,
  required String query,
  required String severityFilter,
  required IssueSortField sortField,
  required bool ascending,
}) {
  final normalizedQuery = query.trim().toLowerCase();
  final filtered = issues.where((issue) {
    final severityOk =
        severityFilter == 'all' || issue.severity == severityFilter;
    if (!severityOk) return false;
    if (normalizedQuery.isEmpty) return true;

    final blob =
        '${issue.message} ${issue.file} ${issue.rule} ${issue.field}'.toLowerCase();
    return blob.contains(normalizedQuery);
  }).toList(growable: false);

  int compare(ValidationIssue a, ValidationIssue b) {
    switch (sortField) {
      case IssueSortField.severity:
        return a.severity.compareTo(b.severity);
      case IssueSortField.file:
        return a.file.compareTo(b.file);
      case IssueSortField.rule:
        return a.rule.compareTo(b.rule);
      case IssueSortField.field:
        return a.field.compareTo(b.field);
      case IssueSortField.row:
        return (a.row ?? -1).compareTo(b.row ?? -1);
      case IssueSortField.message:
        return a.message.compareTo(b.message);
    }
  }

  filtered.sort((a, b) {
    final result = compare(a, b);
    return ascending ? result : -result;
  });
  return filtered;
}

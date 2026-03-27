import '../models/validation_report.dart';

enum IssueSortField {
  severity,
  file,
  rule,
  field,
  row,
  value,
  message,
}

List<ValidationIssue> filterAndSortIssues({
  required List<ValidationIssue> issues,
  String query = '',
  String severityFilter = 'all',
  String fileFilter = 'all',
  String ruleFilter = 'all',
  IssueSortField sortField = IssueSortField.severity,
  bool ascending = true,
}) {
  final normalizedQuery = query.trim().toLowerCase();
  final normalizedSeverityFilter = severityFilter.toLowerCase();
  final filtered = issues.where((issue) {
    final severityOk = normalizedSeverityFilter == 'all' ||
        issue.severity.toLowerCase() == normalizedSeverityFilter;
    if (!severityOk) return false;
    final fileOk = fileFilter == 'all' || issue.file == fileFilter;
    if (!fileOk) return false;
    final ruleOk = ruleFilter == 'all' || issue.rule == ruleFilter;
    if (!ruleOk) return false;
    if (normalizedQuery.isEmpty) return true;

    final blob =
        '${issue.message} ${issue.file} ${issue.rule} ${issue.field} ${issue.value ?? ''} ${issue.row ?? ''} ${issue.severity}'
            .toLowerCase();
    return blob.contains(normalizedQuery);
  }).toList(growable: false);

  int compare(ValidationIssue a, ValidationIssue b) {
    switch (sortField) {
      case IssueSortField.severity:
        final rankDelta = _severityRank(a.severity).compareTo(
          _severityRank(b.severity),
        );
        if (rankDelta != 0) {
          return rankDelta;
        }
        return a.severity.compareTo(b.severity);
      case IssueSortField.file:
        return a.file.compareTo(b.file);
      case IssueSortField.rule:
        return a.rule.compareTo(b.rule);
      case IssueSortField.field:
        return a.field.compareTo(b.field);
      case IssueSortField.row:
        return _rowSortValue(a.row).compareTo(_rowSortValue(b.row));
      case IssueSortField.value:
        return _valueSortValue(a.value).compareTo(_valueSortValue(b.value));
      case IssueSortField.message:
        return a.message.compareTo(b.message);
    }
  }

  filtered.sort((a, b) {
    var result = compare(a, b);
    if (result == 0) {
      result = _rowSortValue(a.row).compareTo(_rowSortValue(b.row));
    }
    if (result == 0) {
      result = a.message.compareTo(b.message);
    }
    return ascending ? result : -result;
  });
  return filtered;
}

int _rowSortValue(int? row) => row ?? 1 << 30;

String _valueSortValue(String? value) {
  if (value == null || value.isEmpty) {
    return '\uffff';
  }
  return value;
}

int _severityRank(String severity) {
  switch (severity.toLowerCase()) {
    case 'critical':
    case 'fatal':
      return 0;
    case 'error':
      return 1;
    case 'warning':
      return 2;
    case 'info':
      return 3;
    default:
      return 4;
  }
}

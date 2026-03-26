class ValidationIssue {
  final String file;
  final String rule;
  final String field;
  final int? row;
  final String? value;
  final String message;
  final String severity;

  ValidationIssue({
    required this.file,
    required this.rule,
    required this.field,
    required this.row,
    required this.value,
    required this.message,
    required this.severity,
  });

  factory ValidationIssue.fromJson(Map<String, dynamic> json) {
    return ValidationIssue(
      file: json['file']?.toString() ?? '',
      rule: json['rule']?.toString() ?? '',
      field: json['field']?.toString() ?? '',
      row: (json['row'] as num?)?.toInt(),
      value: json['value']?.toString(),
      message: json['message']?.toString() ?? '',
      severity: json['severity']?.toString() ?? 'error',
    );
  }
}

class ValidationSummary {
  final int filesChecked;
  final int passed;
  final int failed;
  final int durationMs;

  ValidationSummary({
    required this.filesChecked,
    required this.passed,
    required this.failed,
    required this.durationMs,
  });

  factory ValidationSummary.fromJson(Map<String, dynamic>? json) {
    final j = json ?? const <String, dynamic>{};
    return ValidationSummary(
      filesChecked: (j['files_checked'] as num?)?.toInt() ?? 0,
      passed: (j['passed'] as num?)?.toInt() ?? 0,
      failed: (j['failed'] as num?)?.toInt() ?? 0,
      durationMs: (j['duration_ms'] as num?)?.toInt() ?? 0,
    );
  }
}

class ValidationReport {
  final String schemaVersion;
  final ValidationSummary summary;
  final List<ValidationIssue> issues;

  ValidationReport({
    required this.schemaVersion,
    required this.summary,
    required this.issues,
  });

  factory ValidationReport.fromJson(Map<String, dynamic> json) {
    final rawIssues = (json['issues'] as List<dynamic>? ?? const <dynamic>[])
        .cast<Map<String, dynamic>>();

    return ValidationReport(
      schemaVersion: json['schema_version']?.toString() ?? '',
      summary: ValidationSummary.fromJson(json['summary'] as Map<String, dynamic>?),
      issues: rawIssues.map(ValidationIssue.fromJson).toList(growable: false),
    );
  }
}

class ValidationResult {
  final int exitCode;
  final ValidationReport? report;
  final String stdoutText;
  final String stderrText;

  const ValidationResult({
    required this.exitCode,
    required this.report,
    required this.stdoutText,
    required this.stderrText,
  });

  bool get isPass => exitCode == 0;
  bool get hasValidationFailure => exitCode == 1;
  bool get isRuntimeError => exitCode >= 2;
}

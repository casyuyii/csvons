export type IssueSortField =
  | 'severity'
  | 'file'
  | 'rule'
  | 'field'
  | 'row'
  | 'value'
  | 'message';

export interface ValidationIssue {
  file?: string;
  rule?: string;
  field?: string;
  row?: number;
  value?: string;
  message: string;
  severity: string;
}

function severityRank(severity: string): number {
  switch (severity.toLowerCase()) {
    case 'critical':
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

export function filterAndSortIssues(opts: {
  issues: ValidationIssue[];
  query?: string;
  severityFilter?: string;
  fileFilter?: string;
  ruleFilter?: string;
  sortField?: IssueSortField;
  ascending?: boolean;
}): ValidationIssue[] {
  const {
    issues,
    query = '',
    severityFilter = 'all',
    fileFilter = 'all',
    ruleFilter = 'all',
    sortField = 'severity',
    ascending = true,
  } = opts;

  let filtered = issues;

  if (severityFilter !== 'all') {
    filtered = filtered.filter(
      (i) => i.severity.toLowerCase() === severityFilter.toLowerCase(),
    );
  }

  if (fileFilter !== 'all') {
    filtered = filtered.filter((i) => i.file === fileFilter);
  }

  if (ruleFilter !== 'all') {
    filtered = filtered.filter((i) => i.rule === ruleFilter);
  }

  if (query.trim()) {
    const q = query.toLowerCase();
    filtered = filtered.filter((i) => {
      const blob = [
        i.message,
        i.file ?? '',
        i.rule ?? '',
        i.field ?? '',
        i.value ?? '',
        i.row != null ? String(i.row) : '',
        i.severity,
      ]
        .join(' ')
        .toLowerCase();
      return blob.includes(q);
    });
  }

  const sorted = [...filtered].sort((a, b) => {
    let cmp = 0;

    switch (sortField) {
      case 'severity':
        cmp = severityRank(a.severity) - severityRank(b.severity);
        if (cmp === 0) cmp = a.severity.localeCompare(b.severity);
        break;
      case 'row':
        cmp = (a.row ?? 1 << 30) - (b.row ?? 1 << 30);
        break;
      case 'file':
        cmp = (a.file ?? '').localeCompare(b.file ?? '');
        break;
      case 'rule':
        cmp = (a.rule ?? '').localeCompare(b.rule ?? '');
        break;
      case 'field':
        cmp = (a.field ?? '').localeCompare(b.field ?? '');
        break;
      case 'value':
        cmp = (a.value ?? '\uffff').localeCompare(b.value ?? '\uffff');
        break;
      case 'message':
        cmp = a.message.localeCompare(b.message);
        break;
    }

    if (cmp === 0) cmp = (a.row ?? 1 << 30) - (b.row ?? 1 << 30);
    if (cmp === 0) cmp = a.message.localeCompare(b.message);

    return ascending ? cmp : -cmp;
  });

  return sorted;
}

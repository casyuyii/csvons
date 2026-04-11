import { useState, useMemo } from 'react';
import {
  filterAndSortIssues,
  type ValidationIssue,
  type IssueSortField,
} from '../lib/issue-filters.ts';

interface Props {
  issues: ValidationIssue[];
}

const columns: { key: IssueSortField; label: string }[] = [
  { key: 'severity', label: 'Severity' },
  { key: 'file', label: 'File' },
  { key: 'rule', label: 'Rule' },
  { key: 'field', label: 'Field' },
  { key: 'row', label: 'Row' },
  { key: 'value', label: 'Value' },
  { key: 'message', label: 'Message' },
];

export function IssuesTable({ issues }: Props) {
  const [query, setQuery] = useState('');
  const [severityFilter, setSeverityFilter] = useState('all');
  const [fileFilter, setFileFilter] = useState('all');
  const [ruleFilter, setRuleFilter] = useState('all');
  const [sortField, setSortField] = useState<IssueSortField>('severity');
  const [ascending, setAscending] = useState(true);

  const uniqueFiles = useMemo(
    () => [...new Set(issues.map((i) => i.file).filter(Boolean))],
    [issues],
  );
  const uniqueRules = useMemo(
    () => [...new Set(issues.map((i) => i.rule).filter(Boolean))],
    [issues],
  );

  const filtered = useMemo(
    () =>
      filterAndSortIssues({
        issues,
        query,
        severityFilter,
        fileFilter,
        ruleFilter,
        sortField,
        ascending,
      }),
    [issues, query, severityFilter, fileFilter, ruleFilter, sortField, ascending],
  );

  // Severity counts from filtered issues
  const severityCounts = useMemo(() => {
    const counts = new Map<string, number>();
    for (const issue of filtered) {
      counts.set(issue.severity, (counts.get(issue.severity) ?? 0) + 1);
    }
    return counts;
  }, [filtered]);

  function handleSort(field: IssueSortField) {
    if (sortField === field) {
      setAscending(!ascending);
    } else {
      setSortField(field);
      setAscending(true);
    }
  }

  function resetFilters() {
    setQuery('');
    setSeverityFilter('all');
    setFileFilter('all');
    setRuleFilter('all');
  }

  const hasActiveFilters =
    query || severityFilter !== 'all' || fileFilter !== 'all' || ruleFilter !== 'all';

  return (
    <div>
      {/* Filter bar */}
      <div
        style={{
          display: 'flex',
          gap: 8,
          alignItems: 'center',
          flexWrap: 'wrap',
          marginBottom: 8,
        }}
      >
        <input
          type="text"
          value={query}
          onChange={(e) => setQuery(e.target.value)}
          placeholder="Search issues..."
          style={{
            padding: '6px 10px',
            border: '1px solid #ccc',
            borderRadius: 4,
            fontSize: 13,
            width: 200,
          }}
        />
        <select
          value={fileFilter}
          onChange={(e) => setFileFilter(e.target.value)}
          style={selectStyle}
        >
          <option value="all">All files</option>
          {uniqueFiles.map((f) => (
            <option key={f} value={f}>
              {f}
            </option>
          ))}
        </select>
        <select
          value={ruleFilter}
          onChange={(e) => setRuleFilter(e.target.value)}
          style={selectStyle}
        >
          <option value="all">All rules</option>
          {uniqueRules.map((r) => (
            <option key={r} value={r}>
              {r}
            </option>
          ))}
        </select>
        {hasActiveFilters && (
          <button onClick={resetFilters} style={resetBtnStyle}>
            Reset filters
          </button>
        )}
      </div>

      {/* Severity chips */}
      <div style={{ display: 'flex', gap: 4, marginBottom: 8 }}>
        <button
          onClick={() => setSeverityFilter('all')}
          style={{
            ...chipBtn,
            background: severityFilter === 'all' ? '#3f51b5' : '#e0e0e0',
            color: severityFilter === 'all' ? '#fff' : '#333',
          }}
        >
          All ({filtered.length})
        </button>
        {[...severityCounts.entries()]
          .sort(
            ([a], [b]) => severityRank(a) - severityRank(b),
          )
          .map(([sev, count]) => (
            <button
              key={sev}
              onClick={() =>
                setSeverityFilter(
                  severityFilter === sev ? 'all' : sev,
                )
              }
              style={{
                ...chipBtn,
                background:
                  severityFilter === sev ? severityColor(sev) : '#e0e0e0',
                color: severityFilter === sev ? '#fff' : '#333',
              }}
            >
              {titleCase(sev)} ({count})
            </button>
          ))}
      </div>

      {/* Count */}
      <div style={{ fontSize: 13, color: '#666', marginBottom: 8 }}>
        Showing {filtered.length} of {issues.length} issues
        {hasActiveFilters ? ' (filtered)' : ''}
      </div>

      {/* Table */}
      <div style={{ overflowX: 'auto' }}>
        <table
          style={{
            width: '100%',
            borderCollapse: 'collapse',
            fontSize: 13,
          }}
        >
          <thead>
            <tr>
              {columns.map((col) => (
                <th
                  key={col.key}
                  onClick={() => handleSort(col.key)}
                  style={{
                    ...thStyle,
                    cursor: 'pointer',
                    userSelect: 'none',
                  }}
                >
                  {col.label}
                  {sortField === col.key && (ascending ? ' ▲' : ' ▼')}
                </th>
              ))}
            </tr>
          </thead>
          <tbody>
            {filtered.map((issue, i) => (
              <tr
                key={i}
                style={{ background: i % 2 === 0 ? '#fff' : '#fafafa' }}
              >
                <td style={tdStyle}>
                  <span
                    style={{
                      padding: '2px 6px',
                      borderRadius: 3,
                      background: severityColor(issue.severity),
                      color: '#fff',
                      fontSize: 11,
                      fontWeight: 600,
                    }}
                  >
                    {issue.severity}
                  </span>
                </td>
                <td style={tdStyle}>{issue.file ?? ''}</td>
                <td style={tdStyle}>{issue.rule ?? ''}</td>
                <td style={tdStyle}>{issue.field ?? ''}</td>
                <td style={tdStyle}>{issue.row ?? ''}</td>
                <td style={{ ...tdStyle, maxWidth: 320, overflow: 'hidden', textOverflow: 'ellipsis' }}>
                  {issue.value ?? ''}
                </td>
                <td style={tdStyle}>{issue.message}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
}

function severityRank(s: string): number {
  switch (s.toLowerCase()) {
    case 'critical': return 0;
    case 'error': return 1;
    case 'warning': return 2;
    case 'info': return 3;
    default: return 4;
  }
}

function severityColor(s: string): string {
  switch (s.toLowerCase()) {
    case 'critical': return '#b71c1c';
    case 'error': return '#c62828';
    case 'warning': return '#f57f17';
    case 'info': return '#1565c0';
    default: return '#757575';
  }
}

function titleCase(s: string): string {
  return s.replace(/\b\w/g, (c) => c.toUpperCase());
}

const thStyle: React.CSSProperties = {
  padding: '8px 10px',
  textAlign: 'left',
  borderBottom: '2px solid #ddd',
  fontWeight: 600,
  whiteSpace: 'nowrap',
};

const tdStyle: React.CSSProperties = {
  padding: '6px 10px',
  borderBottom: '1px solid #eee',
  whiteSpace: 'nowrap',
};

const selectStyle: React.CSSProperties = {
  padding: '6px 10px',
  border: '1px solid #ccc',
  borderRadius: 4,
  fontSize: 13,
  background: '#fff',
};

const chipBtn: React.CSSProperties = {
  padding: '4px 10px',
  border: 'none',
  borderRadius: 12,
  cursor: 'pointer',
  fontSize: 12,
  fontWeight: 500,
};

const resetBtnStyle: React.CSSProperties = {
  padding: '6px 10px',
  border: 'none',
  background: 'transparent',
  color: '#3f51b5',
  cursor: 'pointer',
  fontSize: 13,
  textDecoration: 'underline',
};

import { useState, useEffect } from 'react';
import { IssuesTable } from './IssuesTable.tsx';

declare global {
  interface Window {
    csvonsAPI: any;
  }
}

interface ValidationReport {
  schema_version: string;
  summary: {
    files_checked: number;
    passed: number;
    failed: number;
    duration_ms: number;
  };
  issues: any[];
}

export function ValidatePage() {
  const [rulerPath, setRulerPath] = useState('');
  const [exportPath, setExportPath] = useState('');
  const [running, setRunning] = useState(false);
  const [exporting, setExporting] = useState(false);
  const [report, setReport] = useState<ValidationReport | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [exportMessage, setExportMessage] = useState<string | null>(null);
  const [exportIsError, setExportIsError] = useState(false);
  const [recentRulerPaths, setRecentRulerPaths] = useState<string[]>([]);
  const [recentExportPaths, setRecentExportPaths] = useState<string[]>([]);

  useEffect(() => {
    loadState();
  }, []);

  async function loadState() {
    try {
      const state = await window.csvonsAPI.loadState();
      setRecentRulerPaths(state.recentRulerPaths ?? []);
      setRecentExportPaths(state.recentExportPaths ?? []);
    } catch {}
  }

  async function handleSelectRuler() {
    const path = await window.csvonsAPI.selectFile({
      filters: [{ name: 'JSON', extensions: ['json'] }],
    });
    if (path) setRulerPath(path);
  }

  async function handleRun() {
    if (!rulerPath) {
      setError('Please select a ruler.json file');
      return;
    }

    setRunning(true);
    setError(null);
    setReport(null);
    setExportMessage(null);

    try {
      // Save recent path
      const newRecent = [
        rulerPath,
        ...recentRulerPaths.filter((p) => p !== rulerPath),
      ].slice(0, 8);
      setRecentRulerPaths(newRecent);
      await window.csvonsAPI.saveState({ recentRulerPaths: newRecent });

      const result = await window.csvonsAPI.validate(rulerPath);
      setReport(result);
    } catch (err: any) {
      setError(err.message ?? 'Validation failed');
    } finally {
      setRunning(false);
    }
  }

  async function handleExport(format: 'json' | 'markdown') {
    if (!report || !exportPath) return;
    setExporting(true);
    setExportMessage(null);

    try {
      const newRecent = [
        exportPath,
        ...recentExportPaths.filter((p) => p !== exportPath),
      ].slice(0, 8);
      setRecentExportPaths(newRecent);
      await window.csvonsAPI.saveState({ recentExportPaths: newRecent });

      const writtenPath = await window.csvonsAPI.exportReport(
        report,
        exportPath,
        format,
      );
      setExportMessage(`Exported to ${writtenPath}`);
      setExportIsError(false);
    } catch (err: any) {
      setExportMessage(err.message ?? 'Export failed');
      setExportIsError(true);
    } finally {
      setExporting(false);
    }
  }

  return (
    <div style={{ maxWidth: 1100, margin: '0 auto' }}>
      <h2 style={{ marginBottom: 16 }}>Validate CSV Files</h2>

      {/* Ruler path */}
      <div style={{ display: 'flex', gap: 8, marginBottom: 8 }}>
        <input
          type="text"
          value={rulerPath}
          onChange={(e) => setRulerPath(e.target.value)}
          placeholder="Path to ruler.json"
          style={{
            flex: 1,
            padding: '8px 12px',
            border: '1px solid #ccc',
            borderRadius: 4,
            fontSize: 14,
          }}
        />
        <button onClick={handleSelectRuler} style={btnOutlined}>
          Browse...
        </button>
      </div>

      {/* Recent ruler paths */}
      {recentRulerPaths.length > 0 && (
        <div style={{ display: 'flex', gap: 4, flexWrap: 'wrap', marginBottom: 12 }}>
          {recentRulerPaths.slice(0, 3).map((p) => (
            <button
              key={p}
              onClick={() => setRulerPath(p)}
              style={chipStyle}
            >
              {p.split('/').pop()}
            </button>
          ))}
        </div>
      )}

      {/* Run button */}
      <button
        onClick={handleRun}
        disabled={running}
        style={{
          ...btnFilled,
          opacity: running ? 0.7 : 1,
          marginBottom: 16,
        }}
      >
        {running ? 'Running...' : 'Run Validation'}
      </button>

      {/* Error */}
      {error && (
        <div style={{ color: '#c62828', marginBottom: 12, fontSize: 14 }}>
          {error}
        </div>
      )}

      {/* Results */}
      {report && (
        <div>
          {/* Summary */}
          <div
            style={{
              display: 'flex',
              gap: 16,
              alignItems: 'center',
              marginBottom: 12,
              padding: '8px 12px',
              background: report.summary.failed === 0 ? '#e8f5e9' : '#ffebee',
              borderRadius: 6,
              fontSize: 14,
            }}
          >
            <span
              style={{
                width: 10,
                height: 10,
                borderRadius: '50%',
                background: report.summary.failed === 0 ? '#4caf50' : '#f44336',
                display: 'inline-block',
              }}
            />
            <span>
              Files: {report.summary.files_checked} | Passed:{' '}
              {report.summary.passed} | Failed: {report.summary.failed} |{' '}
              {report.summary.duration_ms}ms
            </span>
          </div>

          {/* Issues table */}
          {report.issues.length > 0 && (
            <IssuesTable issues={report.issues} />
          )}

          {report.issues.length === 0 && (
            <div
              style={{
                padding: 24,
                textAlign: 'center',
                color: '#4caf50',
                fontSize: 16,
              }}
            >
              All validations passed!
            </div>
          )}

          {/* Export section */}
          <div
            style={{
              marginTop: 16,
              padding: 12,
              border: '1px solid #e0e0e0',
              borderRadius: 6,
            }}
          >
            <h4 style={{ marginBottom: 8 }}>Export Report</h4>
            <div style={{ display: 'flex', gap: 8, marginBottom: 8 }}>
              <input
                type="text"
                value={exportPath}
                onChange={(e) => setExportPath(e.target.value)}
                placeholder="Export file path"
                style={{
                  flex: 1,
                  padding: '8px 12px',
                  border: '1px solid #ccc',
                  borderRadius: 4,
                  fontSize: 14,
                }}
              />
            </div>

            {recentExportPaths.length > 0 && (
              <div
                style={{
                  display: 'flex',
                  gap: 4,
                  flexWrap: 'wrap',
                  marginBottom: 8,
                }}
              >
                {recentExportPaths.slice(0, 3).map((p) => (
                  <button
                    key={p}
                    onClick={() => setExportPath(p)}
                    style={chipStyle}
                  >
                    {p.split('/').pop()}
                  </button>
                ))}
              </div>
            )}

            <div style={{ display: 'flex', gap: 8 }}>
              <button
                onClick={() => handleExport('json')}
                disabled={exporting || !exportPath}
                style={btnOutlined}
              >
                Export JSON
              </button>
              <button
                onClick={() => handleExport('markdown')}
                disabled={exporting || !exportPath}
                style={btnOutlined}
              >
                Export Markdown
              </button>
            </div>

            {exportMessage && (
              <div
                style={{
                  marginTop: 8,
                  fontSize: 13,
                  color: exportIsError ? '#c62828' : '#2e7d32',
                }}
              >
                {exportMessage}
              </div>
            )}
          </div>
        </div>
      )}
    </div>
  );
}

const btnFilled: React.CSSProperties = {
  padding: '10px 20px',
  background: '#3f51b5',
  color: '#fff',
  border: 'none',
  borderRadius: 4,
  cursor: 'pointer',
  fontSize: 14,
  fontWeight: 500,
};

const btnOutlined: React.CSSProperties = {
  padding: '8px 16px',
  background: 'transparent',
  color: '#3f51b5',
  border: '1px solid #3f51b5',
  borderRadius: 4,
  cursor: 'pointer',
  fontSize: 14,
};

const chipStyle: React.CSSProperties = {
  padding: '4px 10px',
  background: '#e8eaf6',
  border: 'none',
  borderRadius: 12,
  cursor: 'pointer',
  fontSize: 12,
  color: '#3f51b5',
};

interface Props {
  headers: string[];
  rows: string[][];
  totalRows: number;
  filePath: string;
}

export function CsvPreviewCard({ headers, rows, totalRows, filePath }: Props) {
  return (
    <div
      style={{
        border: '1px solid #e0e0e0',
        borderRadius: 6,
        padding: 12,
      }}
    >
      {/* Header chips */}
      <div style={{ marginBottom: 12 }}>
        <div
          style={{
            fontSize: 12,
            color: '#666',
            marginBottom: 4,
          }}
        >
          {headers.length} columns | {totalRows} rows
        </div>
        <div style={{ display: 'flex', gap: 4, flexWrap: 'wrap' }}>
          {headers.map((h, i) => (
            <span
              key={i}
              style={{
                padding: '2px 8px',
                background: '#e8eaf6',
                borderRadius: 10,
                fontSize: 12,
                color: '#3f51b5',
              }}
            >
              {h}
            </span>
          ))}
        </div>
      </div>

      {/* Sample rows */}
      <div style={{ overflowX: 'auto' }}>
        <table
          style={{
            width: '100%',
            borderCollapse: 'collapse',
            fontSize: 12,
          }}
        >
          <thead>
            <tr>
              {headers.map((h, i) => (
                <th
                  key={i}
                  style={{
                    padding: '6px 8px',
                    textAlign: 'left',
                    borderBottom: '2px solid #ddd',
                    fontWeight: 600,
                    whiteSpace: 'nowrap',
                  }}
                >
                  {h}
                </th>
              ))}
            </tr>
          </thead>
          <tbody>
            {rows.map((row, ri) => (
              <tr
                key={ri}
                style={{ background: ri % 2 === 0 ? '#fff' : '#fafafa' }}
              >
                {row.map((cell, ci) => (
                  <td
                    key={ci}
                    style={{
                      padding: '4px 8px',
                      borderBottom: '1px solid #eee',
                      whiteSpace: 'nowrap',
                      maxWidth: 200,
                      overflow: 'hidden',
                      textOverflow: 'ellipsis',
                    }}
                  >
                    {cell}
                  </td>
                ))}
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      {totalRows > rows.length && (
        <div
          style={{
            marginTop: 8,
            fontSize: 12,
            color: '#999',
            textAlign: 'center',
          }}
        >
          Showing {rows.length} of {totalRows} rows
        </div>
      )}
    </div>
  );
}

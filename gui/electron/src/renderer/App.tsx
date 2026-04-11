import { useState } from 'react';
import { ValidatePage } from './components/ValidatePage.tsx';
import { WorkspacePage } from './components/WorkspacePage.tsx';

const tabs = [
  { label: 'Validate', icon: '▶' },
  { label: 'Workspace', icon: '📁' },
] as const;

export function App() {
  const [activeTab, setActiveTab] = useState(0);

  return (
    <div style={{ display: 'flex', flexDirection: 'column', height: '100vh' }}>
      <main style={{ flex: 1, overflow: 'auto', padding: '16px' }}>
        {activeTab === 0 ? <ValidatePage /> : <WorkspacePage />}
      </main>
      <nav
        style={{
          display: 'flex',
          borderTop: '1px solid #ddd',
          background: '#fff',
        }}
      >
        {tabs.map((tab, i) => (
          <button
            key={tab.label}
            onClick={() => setActiveTab(i)}
            style={{
              flex: 1,
              padding: '12px',
              border: 'none',
              background: activeTab === i ? '#e8eaf6' : 'transparent',
              cursor: 'pointer',
              fontWeight: activeTab === i ? 600 : 400,
              fontSize: '14px',
              color: activeTab === i ? '#3f51b5' : '#666',
            }}
          >
            {tab.icon} {tab.label}
          </button>
        ))}
      </nav>
    </div>
  );
}

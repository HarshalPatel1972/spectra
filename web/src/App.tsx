import { useState, useEffect } from 'react'
import './index.css'

interface Finding {
  id: string;
  algorithm: string;
  source: string;
  file_path: string;
  line_number: number;
  key_size?: number;
  qrs: number;
  risk_band: string;
  migration_effort: string;
}

interface ScanResult {
  scan_root: string;
  started_at: string;
  completed_at: string;
  total_files: number;
  scanned_files: number;
  findings: Finding[];
  aggregate_qrs: number;
  findings_by_band: Record<string, number>;
}

function App() {
  const [data, setData] = useState<ScanResult | null>(null);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    // In production, this will hit the Go server API.
    // In dev, we can mock it or proxy it.
    fetch('/api/findings')
      .then(res => {
        if (!res.ok) throw new Error('Failed to fetch data. Is the Spectra dashboard running?');
        return res.json();
      })
      .then(json => setData(json))
      .catch(err => setError(err.message));
  }, []);

  if (error) {
    return (
      <div className="app">
        <div className="card">
          <h2 style={{ color: 'var(--critical)' }}>Error</h2>
          <p>{error}</p>
        </div>
      </div>
    );
  }

  if (!data) {
    return <div className="loading">Loading Spectra Results...</div>;
  }

  // Derive highest band for the aggregate
  let overallBand = 'SAFE';
  if (data.findings_by_band['CRITICAL'] > 0) overallBand = 'CRITICAL';
  else if (data.findings_by_band['HIGH'] > 0) overallBand = 'HIGH';
  else if (data.findings_by_band['MEDIUM'] > 0) overallBand = 'MEDIUM';
  else if (data.findings_by_band['LOW'] > 0) overallBand = 'LOW';

  return (
    <div className="app">
      <header className="header">
        <div className="logo">SPECTRA</div>
        <div className="meta">
          Scanned {data.scanned_files} files in {data.scan_root}
        </div>
      </header>

      <div className="dashboard-grid">
        <div className="card">
          <div className="card-title">Aggregate QRS</div>
          <div className={`stat-value ${overallBand}`}>{data.aggregate_qrs}</div>
          <span className={`badge ${overallBand}`}>{overallBand}</span>
        </div>
        
        <div className="card">
          <div className="card-title">Risk Distribution</div>
          <div style={{ display: 'flex', gap: '1rem', marginTop: '1rem', flexWrap: 'wrap' }}>
            {['CRITICAL', 'HIGH', 'MEDIUM', 'LOW', 'SAFE'].map(band => (
              <div key={band} style={{ textAlign: 'center' }}>
                <div style={{ fontSize: '1.5rem', fontWeight: '700', color: `var(--${band.toLowerCase()})` }}>
                  {data.findings_by_band[band] || 0}
                </div>
                <div style={{ fontSize: '0.75rem', color: 'var(--text-dim)' }}>{band}</div>
              </div>
            ))}
          </div>
        </div>

        <div className="card">
          <div className="card-title">Total Findings</div>
          <div className="stat-value" style={{ color: 'var(--text)' }}>
            {data.findings?.length || 0}
          </div>
        </div>
      </div>

      <div className="card" style={{ padding: '0', overflowX: 'auto' }}>
        <table>
          <thead>
            <tr>
              <th>Algorithm</th>
              <th>Source</th>
              <th>File</th>
              <th>Line</th>
              <th>QRS</th>
              <th>Band</th>
              <th>Effort</th>
            </tr>
          </thead>
          <tbody>
            {data.findings?.map(f => (
              <tr key={f.id}>
                <td><strong>{f.algorithm}</strong>{f.key_size ? ` (${f.key_size})` : ''}</td>
                <td><span className="badge" style={{ backgroundColor: 'rgba(255,255,255,0.1)' }}>{f.source}</span></td>
                <td><code>{f.file_path}</code></td>
                <td>{f.line_number > 0 ? f.line_number : '-'}</td>
                <td>{f.qrs}</td>
                <td><span className={`badge ${f.risk_band}`}>{f.risk_band}</span></td>
                <td>{f.migration_effort}</td>
              </tr>
            ))}
            {(!data.findings || data.findings.length === 0) && (
              <tr>
                <td colSpan={7} style={{ textAlign: 'center', color: 'var(--text-dim)' }}>
                  No cryptographic assets found.
                </td>
              </tr>
            )}
          </tbody>
        </table>
      </div>
    </div>
  )
}

export default App

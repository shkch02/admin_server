import React, { useState, useEffect } from 'react'
import { getAlerts } from '../api/alerts'
import './Alerts.css'

function Alerts() {
  const [alerts, setAlerts] = useState([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState(null)
  const [limit, setLimit] = useState(50)

  useEffect(() => {
    loadAlerts()
    // Auto-refresh every 5 seconds
    const interval = setInterval(loadAlerts, 5000)
    return () => clearInterval(interval)
  }, [limit])

  const loadAlerts = async () => {
    try {
      setError(null)
      const data = await getAlerts(limit)
      setAlerts(data.alerts || [])
    } catch (err) {
      setError(err.message || '알림을 불러오지 못했습니다.')
    } finally {
      setLoading(false)
    }
  }

  const formatTimestamp = (timestamp) => {
    try {
      const date = new Date(timestamp)
      return date.toLocaleString()
    } catch {
      return timestamp
    }
  }

  const getSeverityColor = (severity) => {
    switch (severity?.toLowerCase()) {
      case 'high':
        return '#e74c3c'
      case 'medium':
        return '#f39c12'
      case 'low':
        return '#3498db'
      default:
        return '#7f8c8d'
    }
  }

  if (loading && alerts.length === 0) {
    return <div className="loading">알림을 불러오는 중...</div>
  }

  return (
    <div>
      <div className="card">
        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
          <h2>보안 알림</h2>
          <div>
            <label style={{ marginRight: '10px' }}>
              표시 개수:
              <input
                type="number"
                value={limit}
                onChange={(e) => setLimit(parseInt(e.target.value) || 50)}
                style={{ width: '80px', marginLeft: '5px', padding: '5px' }}
                min="1"
                max="1000"
              />
            </label>
            <button className="button" onClick={loadAlerts}>
              새로고침
            </button>
          </div>
        </div>
        <p>총 알림 수: {alerts.length}</p>
      </div>

      {error && (
        <div className="error">
          {error}
        </div>
      )}

      {alerts.length === 0 ? (
        <div className="card">
          <p style={{ textAlign: 'center', color: '#7f8c8d' }}>알림이 없습니다.</p>
        </div>
      ) : (
        alerts.map((alert, index) => (
          <div key={index} className="card alert-item">
            <div className="alert-header">
              <div>
                <span className="alert-id">{alert.alert_id}</span>
                <span
                  style={{
                    marginLeft: '10px',
                    padding: '2px 8px',
                    backgroundColor: getSeverityColor(alert.severity),
                    color: 'white',
                    borderRadius: '4px',
                    fontSize: '12px',
                  }}
                >
                  {alert.severity || '미확인'}
                </span>
              </div>
              <span className="alert-timestamp">{formatTimestamp(alert.timestamp)}</span>
            </div>
            <div className="alert-rule">
              <strong>룰:</strong> {alert.rule_id}
            </div>
            <div style={{ marginTop: '5px', fontSize: '14px', color: '#7f8c8d' }}>
              {alert.rule_description}
            </div>
            <div className="alert-pod">
              <strong>Pod:</strong> {alert.pod_name} ({alert.namespace})
            </div>
            {alert.syscall_log && (
              <details style={{ marginTop: '10px' }}>
                <summary style={{ cursor: 'pointer', color: '#3498db' }}>시스템콜 로그 보기</summary>
                <pre style={{
                  marginTop: '10px',
                  padding: '10px',
                  backgroundColor: '#f8f9fa',
                  borderRadius: '4px',
                  fontSize: '12px',
                  overflow: 'auto'
                }}>
                  {JSON.stringify(alert.syscall_log, null, 2)}
                </pre>
              </details>
            )}
          </div>
        ))
      )}
    </div>
  )
}

export default Alerts


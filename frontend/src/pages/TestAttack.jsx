import React, { useState, useEffect } from 'react'
import { triggerTest } from '../api/tests'
import { getRules } from '../api/rules'
import './TestAttack.css'

function TestAttack() {
  const [rules, setRules] = useState([])
  const [loading, setLoading] = useState(false)
  const [triggering, setTriggering] = useState(false)
  const [error, setError] = useState(null)
  const [success, setSuccess] = useState(null)

  useEffect(() => {
    loadRules()
  }, [])

  const loadRules = async () => {
    try {
      setLoading(true)
      const data = await getRules()
      setRules(data.rules || [])
    } catch (err) {
      setError('Failed to load rules: ' + err.message)
    } finally {
      setLoading(false)
    }
  }

  const handleTriggerTest = async (testType) => {
    try {
      setTriggering(true)
      setError(null)
      setSuccess(null)

      const response = await triggerTest(testType)
      setSuccess(`Test triggered successfully! Job: ${response.job_name}`)
    } catch (err) {
      setError(err.response?.data?.error || err.message || 'Failed to trigger test')
    } finally {
      setTriggering(false)
    }
  }

  if (loading) {
    return <div className="loading">Loading rules...</div>
  }

  return (
    <div>
      <div className="card">
        <h2>Test Attack</h2>
        <p>Trigger a test attack to verify that the security rules are working correctly.</p>
      </div>

      {error && (
        <div className="error">
          {error}
        </div>
      )}

      {success && (
        <div className="success-message">
          {success}
        </div>
      )}

      <div className="card">
        <h3>Available Test Attacks</h3>
        <p style={{ marginBottom: '20px', color: '#7f8c8d' }}>
          Select a rule to test. This will create a K8s Job that attempts to perform the prohibited action.
        </p>

        {rules.length === 0 ? (
          <p style={{ color: '#7f8c8d' }}>No rules available for testing</p>
        ) : (
          <div style={{ display: 'grid', gap: '15px' }}>
            {rules.map((rule, index) => (
              <div key={index} className="card" style={{ border: '1px solid #ddd' }}>
                <div style={{ marginBottom: '10px' }}>
                  <strong style={{ fontSize: '16px', color: '#2c3e50' }}>{rule.rule_id}</strong>
                </div>
                <div style={{ marginBottom: '15px', color: '#7f8c8d', fontSize: '14px' }}>
                  {rule.description}
                </div>
                <button
                  className="button danger"
                  onClick={() => handleTriggerTest(rule.rule_id)}
                  disabled={triggering}
                >
                  {triggering ? 'Triggering...' : 'Trigger Test Attack'}
                </button>
              </div>
            ))}
          </div>
        )}
      </div>

      <div className="card">
        <h3>How it works</h3>
        <ul style={{ lineHeight: '1.8' }}>
          <li>Clicking a "Trigger Test Attack" button creates a Kubernetes Job</li>
          <li>The Job attempts to perform the action that the rule is designed to detect</li>
          <li>The eBPF monitor should detect the syscall and trigger an alert</li>
          <li>Check the Alerts page to see if the attack was detected</li>
        </ul>
      </div>
    </div>
  )
}

export default TestAttack


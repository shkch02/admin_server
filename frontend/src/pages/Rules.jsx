import React, { useState, useEffect } from 'react'
import { getRules } from '../api/rules'
import './Rules.css'

function Rules() {
  const [rules, setRules] = useState(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState(null)
  const [showYaml, setShowYaml] = useState(false)

  useEffect(() => {
    loadRules()
  }, [])

  const loadRules = async () => {
    try {
      setLoading(true)
      setError(null)
      const data = await getRules()
      setRules(data)
    } catch (err) {
      setError(err.message || 'Failed to load rules')
    } finally {
      setLoading(false)
    }
  }

  const formatValue = (value) => {
    if (Array.isArray(value)) {
      return value.join(', ')
    }
    return String(value)
  }

  if (loading) {
    return <div className="loading">Loading rules...</div>
  }

  if (error) {
    return <div className="error">Error: {error}</div>
  }

  if (!rules) {
    return <div className="error">No rules found</div>
  }

  return (
    <div>
      <div className="card">
        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
          <h2>Security Rules</h2>
          <button className="button" onClick={loadRules}>
            Refresh
          </button>
        </div>
        <p><strong>Version:</strong> {rules.ruleset_version}</p>
        <p><strong>Description:</strong> {rules.description}</p>
        <p><strong>Total Rules:</strong> {rules.rules.length}</p>
      </div>

      {rules.rules.map((rule, index) => (
        <div key={index} className="card rule-item">
          <div className="rule-id">{rule.rule_id}</div>
          <div className="rule-description">{rule.description}</div>
          <div style={{ marginTop: '15px' }}>
            <strong>Conditions:</strong>
            {rule.conditions.map((condition, condIndex) => (
              <div key={condIndex} className="condition">
                <strong>{condition.field}</strong> {condition.operator} {formatValue(condition.value)}
              </div>
            ))}
          </div>
        </div>
      ))}

      <div className="card">
        <button
          className="button"
          onClick={() => setShowYaml(!showYaml)}
        >
          {showYaml ? 'Hide' : 'Show'} YAML Source
        </button>
        {showYaml && (
          <pre style={{ marginTop: '15px', padding: '15px', backgroundColor: '#f8f9fa', borderRadius: '4px', overflow: 'auto' }}>
            {JSON.stringify(rules, null, 2)}
          </pre>
        )}
      </div>
    </div>
  )
}

export default Rules


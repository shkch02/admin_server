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
      setError(err.message || '룰 정보를 불러오지 못했습니다.')
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
    return <div className="loading">룰 정보를 불러오는 중...</div>
  }

  if (error) {
    return <div className="error">오류: {error}</div>
  }

  if (!rules) {
    return <div className="error">등록된 룰이 없습니다.</div>
  }

  return (
    <div>
      <div className="card">
        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
          <h2>보안 룰 현황</h2>
          <button className="button" onClick={loadRules}>
            새로고침
          </button>
        </div>
        <p><strong>버전:</strong> {rules.ruleset_version}</p>
        <p><strong>설명:</strong> {rules.description}</p>
        <p><strong>총 룰 수:</strong> {rules.rules.length}</p>
      </div>

      {rules.rules.map((rule, index) => (
        <div key={index} className="card rule-item">
          <div className="rule-id">{rule.rule_id}</div>
          <div className="rule-description">{rule.description}</div>
          <div style={{ marginTop: '15px' }}>
            <strong>조건:</strong>
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
          {showYaml ? 'YAML 원본 숨기기' : 'YAML 원본 보기'}
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


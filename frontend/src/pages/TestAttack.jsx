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
      setError('룰을 불러오지 못했습니다: ' + err.message)
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
      setSuccess(`공격 테스트가 실행되었습니다! 생성된 Job: ${response.job_name}`)
    } catch (err) {
      setError(err.response?.data?.error || err.message || '공격 테스트 실행에 실패했습니다.')
    } finally {
      setTriggering(false)
    }
  }

  if (loading) {
    return <div className="loading">룰 정보를 불러오는 중...</div>
  }

  return (
    <div>
      <div className="card">
        <h2>공격 테스트</h2>
        <p>보안 룰이 정상 동작하는지 확인하기 위해 테스트 공격을 실행합니다.</p>
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
        <h3>테스트 가능한 공격</h3>
        <p style={{ marginBottom: '20px', color: '#7f8c8d' }}>
          테스트할 룰을 선택하면 해당 룰을 위반하는 동작을 수행하는 Kubernetes Job이 생성됩니다.
        </p>

        {rules.length === 0 ? (
          <p style={{ color: '#7f8c8d' }}>테스트할 수 있는 룰이 없습니다.</p>
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
                  {triggering ? '실행 중...' : '공격 테스트 실행'}
                </button>
              </div>
            ))}
          </div>
        )}
      </div>

      <div className="card">
        <h3>동작 방식</h3>
        <ul style={{ lineHeight: '1.8' }}>
          <li>"공격 테스트 실행" 버튼을 누르면 Kubernetes Job이 생성됩니다.</li>
          <li>생성된 Job은 해당 룰이 탐지해야 하는 금지 동작을 시도합니다.</li>
          <li>eBPF 모니터가 시스템콜을 감지하여 알림을 발생시켜야 합니다.</li>
          <li>알림 페이지에서 공격이 탐지되었는지 확인하세요.</li>
        </ul>
      </div>
    </div>
  )
}

export default TestAttack


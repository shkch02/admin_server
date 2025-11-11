import React, { useState, useEffect } from 'react'
import { getRules, updateRules } from '../api/rules'
import Editor from '@monaco-editor/react'
import './UpdateRules.css'

function UpdateRules() {
  const [rules, setRules] = useState(null)
  const [yamlContent, setYamlContent] = useState('')
  const [loading, setLoading] = useState(true)
  const [saving, setSaving] = useState(false)
  const [error, setError] = useState(null)
  const [success, setSuccess] = useState(null)

  useEffect(() => {
    loadRules()
  }, [])

  const loadRules = async () => {
    try {
      setLoading(true)
      setError(null)
      const data = await getRules()
      setRules(data)
      // Convert to YAML-like JSON for editing
      setYamlContent(JSON.stringify(data, null, 2))
    } catch (err) {
      setError(err.message || '룰 정보를 불러오지 못했습니다.')
    } finally {
      setLoading(false)
    }
  }

  const handleSave = async () => {
    try {
      setSaving(true)
      setError(null)
      setSuccess(null)

      // Parse JSON
      let parsedRules
      try {
        parsedRules = JSON.parse(yamlContent)
      } catch (parseError) {
        setError('JSON 형식이 올바르지 않습니다: ' + parseError.message)
        setSaving(false)
        return
      }

      // Validate structure
      if (!parsedRules.ruleset_version || !parsedRules.rules) {
        setError('룰 구조가 올바르지 않습니다: ruleset_version 또는 rules 필드가 없습니다.')
        setSaving(false)
        return
      }

      const response = await updateRules(parsedRules)
      setSuccess(`룰이 성공적으로 업데이트되었습니다! 새 버전: ${response.new_version || '미정'}`)
      setRules(parsedRules)
    } catch (err) {
      setError(err.response?.data?.error || err.message || '룰 업데이트에 실패했습니다.')
    } finally {
      setSaving(false)
    }
  }

  const handleAddRule = () => {
    try {
      const parsed = JSON.parse(yamlContent)
      const newRule = {
        rule_id: `RULE_NEW_${Date.now()}`,
        description: '새로운 룰 설명',
        conditions: [
          {
            field: 'syscall_name',
            operator: 'equals',
            value: 'openat'
          }
        ]
      }
      parsed.rules.push(newRule)
      setYamlContent(JSON.stringify(parsed, null, 2))
    } catch (err) {
      setError('룰 추가에 실패했습니다: ' + err.message)
    }
  }

  const handleDeleteRule = (index) => {
    try {
      const parsed = JSON.parse(yamlContent)
      parsed.rules.splice(index, 1)
      setYamlContent(JSON.stringify(parsed, null, 2))
    } catch (err) {
      setError('룰 삭제에 실패했습니다: ' + err.message)
    }
  }

  if (loading) {
    return <div className="loading">룰 정보를 불러오는 중...</div>
  }

  return (
    <div>
      <div className="card">
        <h2>룰 수정</h2>
        <p>아래에서 룰 설정을 수정하면 ConfigMap에 반영됩니다.</p>
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
        <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: '15px' }}>
          <div>
            <button className="button" onClick={loadRules}>
              다시 불러오기
            </button>
            <button className="button success" onClick={handleSave} disabled={saving}>
              {saving ? '저장 중...' : '변경 사항 저장'}
            </button>
          </div>
        </div>

        <div style={{ marginBottom: '15px' }}>
          <label>
            <strong>룰 설정 (JSON 형식)</strong>
          </label>
          <div style={{ border: '1px solid #ddd', borderRadius: '4px', marginTop: '5px' }}>
            <Editor
              height="500px"
              defaultLanguage="json"
              value={yamlContent}
              onChange={(value) => setYamlContent(value || '')}
              theme="vs-light"
              options={{
                minimap: { enabled: false },
                fontSize: 14,
                wordWrap: 'on',
              }}
            />
          </div>
        </div>

        <div style={{ marginTop: '20px' }}>
          <h3>빠른 작업</h3>
          <button className="button" onClick={handleAddRule} style={{ marginRight: '10px' }}>
            새 룰 추가
          </button>
        </div>
      </div>

      {rules && (
        <div className="card">
          <h3>현재 룰 미리보기</h3>
          {rules.rules.map((rule, index) => (
            <div key={index} className="rule-item" style={{ marginBottom: '10px' }}>
              <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                <div>
                  <strong>{rule.rule_id}</strong>
                  <div style={{ fontSize: '14px', color: '#7f8c8d' }}>{rule.description}</div>
                </div>
                <button
                  className="button danger"
                  onClick={() => handleDeleteRule(index)}
                  style={{ fontSize: '12px', padding: '5px 10px' }}
                >
                  삭제
                </button>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  )
}

export default UpdateRules


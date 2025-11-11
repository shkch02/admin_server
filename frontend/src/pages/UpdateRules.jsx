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
      setError(err.message || 'Failed to load rules')
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
        setError('Invalid JSON format: ' + parseError.message)
        setSaving(false)
        return
      }

      // Validate structure
      if (!parsedRules.ruleset_version || !parsedRules.rules) {
        setError('Invalid rule structure: missing ruleset_version or rules')
        setSaving(false)
        return
      }

      const response = await updateRules(parsedRules)
      setSuccess(`Rules updated successfully! New version: ${response.new_version || 'N/A'}`)
      setRules(parsedRules)
    } catch (err) {
      setError(err.response?.data?.error || err.message || 'Failed to update rules')
    } finally {
      setSaving(false)
    }
  }

  const handleAddRule = () => {
    try {
      const parsed = JSON.parse(yamlContent)
      const newRule = {
        rule_id: `RULE_NEW_${Date.now()}`,
        description: 'New rule description',
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
      setError('Failed to add rule: ' + err.message)
    }
  }

  const handleDeleteRule = (index) => {
    try {
      const parsed = JSON.parse(yamlContent)
      parsed.rules.splice(index, 1)
      setYamlContent(JSON.stringify(parsed, null, 2))
    } catch (err) {
      setError('Failed to delete rule: ' + err.message)
    }
  }

  if (loading) {
    return <div className="loading">Loading rules...</div>
  }

  return (
    <div>
      <div className="card">
        <h2>Update Rules</h2>
        <p>Edit the rules configuration below. Changes will be applied to the ConfigMap.</p>
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
              Reload
            </button>
            <button className="button success" onClick={handleSave} disabled={saving}>
              {saving ? 'Saving...' : 'Save Changes'}
            </button>
          </div>
        </div>

        <div style={{ marginBottom: '15px' }}>
          <label>
            <strong>Rules Configuration (JSON format)</strong>
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
          <h3>Quick Actions</h3>
          <button className="button" onClick={handleAddRule} style={{ marginRight: '10px' }}>
            Add New Rule
          </button>
        </div>
      </div>

      {rules && (
        <div className="card">
          <h3>Current Rules Preview</h3>
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
                  Delete
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


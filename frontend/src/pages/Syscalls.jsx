import React, { useState, useEffect } from 'react'
import { getCallableSyscalls } from '../api/syscalls'
import './Syscalls.css'

function Syscalls() {
  const [syscalls, setSyscalls] = useState([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState(null)
  const [searchTerm, setSearchTerm] = useState('')

  useEffect(() => {
    loadSyscalls()
  }, [])

  const loadSyscalls = async () => {
    try {
      setLoading(true)
      setError(null)
      const data = await getCallableSyscalls()
      setSyscalls(data.syscalls || [])
    } catch (err) {
      setError(err.message || 'Failed to load syscalls')
    } finally {
      setLoading(false)
    }
  }

  const filteredSyscalls = syscalls.filter(syscall =>
    syscall.name.toLowerCase().includes(searchTerm.toLowerCase())
  )

  if (loading) {
    return <div className="loading">Loading syscalls...</div>
  }

  return (
    <div>
      <div className="card">
        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
          <h2>Callable System Calls</h2>
          <button className="button" onClick={loadSyscalls}>
            Refresh
          </button>
        </div>
        <p>Total syscalls: {syscalls.length}</p>
        <div style={{ marginTop: '15px' }}>
          <input
            type="text"
            placeholder="Search syscalls..."
            value={searchTerm}
            onChange={(e) => setSearchTerm(e.target.value)}
            style={{
              width: '100%',
              padding: '8px',
              border: '1px solid #ddd',
              borderRadius: '4px',
            }}
          />
        </div>
      </div>

      {error && (
        <div className="error">
          {error}
        </div>
      )}

      {filteredSyscalls.length === 0 ? (
        <div className="card">
          <p style={{ textAlign: 'center', color: '#7f8c8d' }}>
            {searchTerm ? 'No syscalls found matching your search' : 'No syscalls found'}
          </p>
        </div>
      ) : (
        filteredSyscalls.map((syscall, index) => (
          <div key={index} className="card">
            <div style={{ marginBottom: '10px' }}>
              <strong style={{ fontSize: '18px', color: '#2c3e50' }}>{syscall.name}</strong>
            </div>
            {syscall.args && syscall.args.length > 0 ? (
              <div>
                <strong>Arguments:</strong>
                <ul style={{ marginTop: '10px', marginLeft: '20px' }}>
                  {syscall.args.map((arg, argIndex) => (
                    <li key={argIndex} style={{ marginBottom: '5px', fontFamily: 'monospace' }}>
                      <span style={{ color: '#3498db' }}>{arg.type}</span> <strong>{arg.name}</strong>
                    </li>
                  ))}
                </ul>
              </div>
            ) : (
              <p style={{ color: '#7f8c8d' }}>No arguments</p>
            )}
          </div>
        ))
      )}
    </div>
  )
}

export default Syscalls


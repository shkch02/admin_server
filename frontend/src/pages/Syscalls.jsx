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
      setError(err.message || '시스템콜 목록을 불러오지 못했습니다.')
    } finally {
      setLoading(false)
    }
  }

  const filteredSyscalls = syscalls.filter(syscall =>
    syscall.name.toLowerCase().includes(searchTerm.toLowerCase())
  )

  if (loading) {
    return <div className="loading">시스템콜 목록을 불러오는 중...</div>
  }

  return (
    <div>
      <div className="card">
        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
          <h2>호출 가능한 시스템콜</h2>
          <button className="button" onClick={loadSyscalls}>
            새로고침
          </button>
        </div>
        <p>총 시스템콜 수: {syscalls.length}</p>
        <div style={{ marginTop: '15px' }}>
          <input
            type="text"
            placeholder="시스템콜 검색..."
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
            {searchTerm ? '검색 조건에 맞는 시스템콜이 없습니다.' : '시스템콜 정보가 없습니다.'}
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
                <strong>인자:</strong>
                <ul style={{ marginTop: '10px', marginLeft: '20px' }}>
                  {syscall.args.map((arg, argIndex) => (
                    <li key={argIndex} style={{ marginBottom: '5px', fontFamily: 'monospace' }}>
                      <span style={{ color: '#3498db' }}>{arg.type}</span> <strong>{arg.name}</strong>
                    </li>
                  ))}
                </ul>
              </div>
            ) : (
              <p style={{ color: '#7f8c8d' }}>인자 없음</p>
            )}
          </div>
        ))
      )}
    </div>
  )
}

export default Syscalls


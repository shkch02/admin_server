import React from 'react'
import { BrowserRouter as Router, Routes, Route, Link, useLocation } from 'react-router-dom'
import Rules from './pages/Rules'
import UpdateRules from './pages/UpdateRules'
import ClusterDiagram from './pages/ClusterDiagram'
import Syscalls from './pages/Syscalls'
import TestAttack from './pages/TestAttack'
import './App.css'

function Navigation() {
  const location = useLocation()
  
  return (
    <nav className="nav">
      <Link to="/" className={location.pathname === '/' ? 'active' : ''}>
        룰 현황
      </Link>
      <Link to="/update-rules" className={location.pathname === '/update-rules' ? 'active' : ''}>
        룰 수정
      </Link>
      <Link to="/diagram" className={location.pathname === '/diagram' ? 'active' : ''}>
        클러스터 구조도
      </Link>
      <Link to="/syscalls" className={location.pathname === '/syscalls' ? 'active' : ''}>
        시스템콜 목록
      </Link>
      <Link to="/test-attack" className={location.pathname === '/test-attack' ? 'active' : ''}>
        공격 테스트 
      </Link>
    </nav>
  )
}

function App() {
  return (
    <Router>
      <div className="container">
        <div className="header">
          <h1>관리자 페이지</h1>
          <p className="subtitle">eBPF 시스템콜 기반 보안 위협 탐지</p>
          <div className="author-info">
          <span>2025 홍익대학교 컴퓨터공학부 졸업전시</span>
          <span className="divider">|</span>
          <span>C182023 신경철</span>
          </div>
          <Navigation />
        </div>
        <Routes>
          <Route path="/" element={<Rules />} />
          <Route path="/update-rules" element={<UpdateRules />} />
          <Route path="/diagram" element={<ClusterDiagram />} />
          <Route path="/syscalls" element={<Syscalls />} />
          <Route path="/test-attack" element={<TestAttack />} />
        </Routes>
        <p className="version-info">버전 (v2.0.0)</p>
      </div>
    </Router>
  )
}

export default App


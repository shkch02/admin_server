import React from 'react'
import { BrowserRouter as Router, Routes, Route, Link, useLocation } from 'react-router-dom'
import Rules from './pages/Rules'
import UpdateRules from './pages/UpdateRules'
import Alerts from './pages/Alerts'
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
      <Link to="/alerts" className={location.pathname === '/alerts' ? 'active' : ''}>
        알림 로그
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
          <h1>IPS 관리자 콘솔</h1>
          <p>eBPF 시스템콜 기반 보안 위협 탐지, CICD테스트</p>
          <Navigation />
        </div>
        <Routes>
          <Route path="/" element={<Rules />} />
          <Route path="/update-rules" element={<UpdateRules />} />
          <Route path="/alerts" element={<Alerts />} />
          <Route path="/syscalls" element={<Syscalls />} />
          <Route path="/test-attack" element={<TestAttack />} />
        </Routes>
      </div>
    </Router>
  )
}

export default App


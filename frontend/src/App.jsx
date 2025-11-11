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
        Rules
      </Link>
      <Link to="/update-rules" className={location.pathname === '/update-rules' ? 'active' : ''}>
        Update Rules
      </Link>
      <Link to="/alerts" className={location.pathname === '/alerts' ? 'active' : ''}>
        Alerts
      </Link>
      <Link to="/syscalls" className={location.pathname === '/syscalls' ? 'active' : ''}>
        Syscalls
      </Link>
      <Link to="/test-attack" className={location.pathname === '/test-attack' ? 'active' : ''}>
        Test Attack
      </Link>
    </nav>
  )
}

function App() {
  return (
    <Router>
      <div className="container">
        <div className="header">
          <h1>IPS Admin Server</h1>
          <p>eBPF System Call based Security Violation Detection</p>
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


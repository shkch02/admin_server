import React from 'react'

function ClusterDiagram() {
  return (
    <div>
      <div className="card">
        <h2>클러스터 구조도</h2>
        <p>현재 시스템의 전체 아키텍처 다이어그램입니다.</p>
      </div>

      <div className="card" style={{ textAlign: 'center', overflow: 'auto' }}>
        {/* public 폴더에 있는 이미지는 /파일명 으로 바로 접근 가능합니다 */}
        <img 
          src="/cluster_diagram.png" 
          alt="Cluster Architecture" 
          style={{ maxWidth: '100%', height: 'auto', borderRadius: '8px', border: '1px solid #ddd' }}
        />
      </div>
    </div>
  )
}

export default ClusterDiagram
# IPS Admin Server

관리 웹서버 프론트엔드(React)와 백엔드(Golang) 구현


## 프로젝트 구조

```
admin_server/
├── backend/
│   ├── main.go
│   └── internal/
│       ├── config/
│       │   └── config.go
│       ├── handlers/
│       │   ├── rule_handler.go
│       │   ├── syscall_handler.go
│       │   ├── alert_handler.go
│       │   └── test_handler.go
│       ├── models/
│       │   └── models.go
│       └── services/
│           ├── rule_service.go
│           ├── syscall_service.go
│           ├── alert_service.go
│           └── test_service.go
├── frontend/
│   ├── src/
│   │   ├── api/
│   │   │   ├── client.js
│   │   │   ├── rules.js
│   │   │   ├── alerts.js
│   │   │   ├── syscalls.js
│   │   │   └── tests.js
│   │   ├── pages/
│   │   │   ├── Rules.jsx
│   │   │   ├── UpdateRules.jsx
│   │   │   ├── Alerts.jsx
│   │   │   ├── Syscalls.jsx
│   │   │   └── TestAttack.jsx
│   │   ├── App.jsx
│   │   ├── App.css
│   │   ├── main.jsx
│   │   └── index.css
│   ├── package.json
│   ├── vite.config.js
│   └── index.html
└── go.mod
```

## API 엔드포인트

### 1. Rules
- `GET /api/v1/rules` - 현재 룰 조회
- `PUT /api/v1/rules` - 룰 업데이트

### 2. Syscalls
- `GET /api/v1/syscalls/callable` - 클러스터가 호출 가능한 syscall 목록 조회

### 3. Alerts
- `GET /api/v1/alerts` - 알림 로그 조회
- `POST /api/v1/alerts/webhook` - 웹훅으로 알림 수신 (내부 API)

### 4. Tests
- `POST /api/v1/tests/trigger` - 테스트 공격 트리거

## 실행 방법

재부팅시에는 마지막 명령어만

### 백엔드 실행

```bash
cd backend
go mod download
go run main.go
```

백엔드는 기본적으로 `:8080` 포트에서 실행됩니다.

### 프론트엔드 실행

```bash
cd frontend
npm install
npm run dev
```

프론트엔드는 기본적으로 `:3000` 포트에서 실행되며, `/api` 요청은 자동으로 백엔드로 프록시됩니다.

## 환경 변수

백엔드 환경 변수:

- `PORT` - 서버 포트 (기본값: 8080)
- `KUBE_CONFIG_PATH` - Kubernetes 설정 파일 경로
- `NAMESPACE` - Kubernetes 네임스페이스 (기본값: default)
- `CONFIG_MAP_NAME` - ConfigMap 이름 (기본값: rule-yaml)
- `REDIS_HOST` - Redis 호스트 (기본값: localhost)
- `REDIS_PORT` - Redis 포트 (기본값: 6379)
- `REDIS_PASSWORD` - Redis 비밀번호
- `CLUSTER_SYSCALLS_REDIS_HOST` - 클러스터 syscalls Redis 호스트
- `CLUSTER_SYSCALLS_REDIS_PORT` - 클러스터 syscalls Redis 포트
- `CLUSTER_SYSCALLS_REDIS_PASSWORD` - 클러스터 syscalls Redis 비밀번호

## 구현 상태

현재 구현된 기능:
- ✅ API 엔드포인트 기본 구조
- ✅ 프론트엔드 페이지 (Rules, UpdateRules, Alerts, Syscalls, TestAttack)
- ✅ 프론트-백엔드 연결
- ✅ 목업 데이터로 기본 동작 확인 가능  ㅁㄴㅇ 

TODO (실제 구현 필요):
- ⏳ Kubernetes ConfigMap 실제 읽기/쓰기
- ⏳ Redis 연결 및 데이터 조회
- ⏳ Kafka 컨슈머로 알림 수신
- ⏳ Kubernetes Job 생성
- ⏳ 룰 엔진 및 eBPF generator 트리거

pv 마운트해서 각종 다이어그램도 볼수있는 기능 추가하면 좋을듯

## 주의사항

현재는 목업 데이터를 사용하고 있습니다. 실제 Kubernetes 및 Redis 연결은 각 서비스 파일의 TODO 주석 부분을 구현해야 합니다.ㅋ


## 의존성

### 프론트

### 백엔드
 go get k8s.io/client-go@v0.28.0
go get k8s.io/apimachinery@v0.28.0


## TODO

클러스터에 있는 rule-configmap.yaml을 참조하는게 아니라 이 프로젝트 디렉토리에 있는 rule-configmap을 바라보는거같은데,,확인 필요 
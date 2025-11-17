# --- 1. 빌드 스테이지 ---
FROM golang:1.24-alpine AS builder

WORKDIR /app

# go.mod와 go.sum을 먼저 복사하여 종속성 캐싱 활용
COPY go.mod go.sum ./
RUN go mod download

# 소스 코드 복사
COPY . .

# 애플리케이션 빌드
# CGO_ENABLED=0는 정적 바이너리를 생성하여 'scratch' 같은 최소 이미지에서도 작동하게 함
RUN CGO_ENABLED=0 go build -o /app/admin-server ./main.go

# --- 2. 최종 스테이지 ---
FROM alpine:latest

WORKDIR /app

# 빌드 스테이지에서 컴파일된 바이너리 복사# --- 1. 빌드 스테이지 ---
FROM golang:1.24-alpine AS builder

WORKDIR /app

# go.mod와 go.sum을 먼저 복사하여 종속성 캐싱 활용
# 이 단계는 캐시 레이어를 설정합니다.
COPY go.mod go.sum ./

# 소스 코드 복사 (Redis를 import 하는 syscall_service.go 파일 포함)
COPY . .

# [수정된 부분]
# 1. go.mod를 정리하여 필요한 종속성을 정확히 파악합니다.
# 2. 종속성을 다운로드하여 캐시합니다.
# 이 명령은 COPY . . 이후에 위치해야 최신 코드를 기반으로 종속성을 파악합니다.
RUN go mod tidy && go mod download

# 애플리케이션 빌드
# CGO_ENABLED=0는 정적 바이너리를 생성하여 'scratch' 같은 최소 이미지에서도 작동하게 함
RUN CGO_ENABLED=0 go build -o /app/admin-server ./main.go

# --- 2. 최종 스테이지 ---
FROM alpine:latest

WORKDIR /app

# 빌드 스테이지에서 컴파일된 바이너리 복사
COPY --from=builder /app/admin-server .

# 애플리케이션 실행 포트 (main.go에서 8080으로 설정됨)
EXPOSE 8080

# 컨테이너 시작 명령어
CMD ["/app/admin-server"]
COPY --from=builder /app/admin-server .

# 애플리케이션 실행 포트 (main.go에서 8080으로 설정됨)
EXPOSE 8080

# 컨테이너 시작 명령어
CMD ["/app/admin-server"]
# --- 1. 빌드 스테이지 ---
FROM golang:1.21-alpine AS builder

WORKDIR /app

# go.mod와 go.sum을 먼저 복사하여 종속성 캐싱 활용
COPY ../go.mod ../go.sum ./
RUN go mod download

# 소스 코드 복사
COPY . .

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
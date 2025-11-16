# --- 1. 빌드 스테이지 ---
# package.json에 명시된 React 버전에 맞는 Node.js 이미지 사용
FROM node:18-alpine AS builder

WORKDIR /app

# 종속성 설치
COPY package.json package-lock.json ./
RUN npm install

# 소스 코드 복사
COPY . .

# package.json의 "build" 스크립트 실행
RUN npm run build

# --- 2. 최종 스테이지 ---
FROM nginx:1.25-alpine

# 1. 빌드 스테이지에서 생성된 정적 파일(dist 폴더)을 Nginx 웹 서버 루트로 복사
COPY --from=builder /app/dist /usr/share/nginx/html

# 2. 1단계에서 작성한 커스텀 Nginx 설정을 Nginx 설정 디렉터리로 복사
COPY nginx.conf /etc/nginx/conf.d/default.conf

# Nginx는 기본적으로 80 포트를 사용
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
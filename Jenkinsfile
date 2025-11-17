pipeline {
    agent any // kubectl, docker, kustomize가 설치된 에이전트 사용

    environment {
        // 1. Harbor 및 이미지 정보
        HARBOR_URL       = "shkch.duckdns.org"
        HARBOR_PROJECT   = "webserver"
        HARBOR_CREDS_ID  = "harbor-creds"
        
        // 2. 백엔드/프론트엔드 이미지 이름
        BACKEND_IMAGE_NAME  = "admin-server"
        FRONTEND_IMAGE_NAME = "admin-frontend"
        
        // 3. Kubeconfig 자격 증명 ID
        KUBE_CREDS_ID = "kubeconfig-creds" // Secret File 타입 사용
    }

    stages {
        // 1단계: Git 저장소에서 코드 가져오기
        stage('Checkout') {
            steps {
                // SCM 방식이므로 'checkout scm'이 맞습니다.
                checkout scm
            }
        }

        // 2단계: 이미지 태그 정의 (Git Commit 해시 사용)
        stage('Define Image Tag') {
            steps {
                script {
                    // Git Commit 해시의 앞 8자리를 이미지 태그로 사용
                    env.IMAGE_TAG = sh(returnStdout: true, script: 'git rev-parse --short=8 HEAD').trim()
                    echo "Using Image Tag: ${env.IMAGE_TAG}"
                }
            }
        }

        // 3단계: 이미지 빌드 및 푸시 (순차 실행으로 변경)
        stage('Build & Push Images') {
            steps {
                // Harbor 로그인 (파이프라인 시작 시 한 번만)
                withCredentials([usernamePassword(credentialsId: env.HARBOR_CREDS_ID, usernameVariable: 'HARBOR_USER', passwordVariable: 'HARBOR_PASS')]) {
                    sh "docker login ${env.HARBOR_URL} -u ${HARBOR_USER} -p '${HARBOR_PASS}'"
                }

                // 1. 백엔드 빌드/푸시 (먼저 실행)
                echo "Building Backend Image..."
                dir('backend') { // 'backend' 디렉터리로 이동
                    sh "docker build -t ${env.HARBOR_URL}/${env.HARBOR_PROJECT}/${env.BACKEND_IMAGE_NAME}:${env.IMAGE_TAG} ."
                    sh "docker push ${env.HARBOR_URL}/${env.HARBOR_PROJECT}/${env.BACKEND_IMAGE_NAME}:${env.IMAGE_TAG}"
                }

                // 2. 프론트엔드 빌드/푸시 (백엔드 완료 후 실행)
                echo "Building Frontend Image..."
                dir('frontend') { // 'frontend' 디렉터리로 이동
                    sh "docker build -t ${env.HARBOR_URL}/${env.HARBOR_PROJECT}/${env.FRONTEND_IMAGE_NAME}:${env.IMAGE_TAG} ."
                    sh "docker push ${env.HARBOR_URL}/${env.HARBOR_PROJECT}/${env.FRONTEND_IMAGE_NAME}:${env.IMAGE_TAG}"
                }
            }
        }

        // 4단계: Kubernetes에 배포
        stage('Deploy to Kubernetes') {
            steps {
                withCredentials([string(credentialsId: 'saToken', variable: 'K8S_TOKEN')]) {

                    dir('k8s') {
                        echo "Deploying with image tag: ${env.IMAGE_TAG}"

                        // 1) 이미지 태그 업데이트
                        sh """
                            kustomize edit set image ${env.HARBOR_URL}/${env.HARBOR_PROJECT}/${env.BACKEND_IMAGE_NAME}=${env.HARBOR_URL}/${env.HARBOR_PROJECT}/${env.BACKEND_IMAGE_NAME}:${env.IMAGE_TAG}
                            kustomize edit set image ${env.HARBOR_URL}/${env.HARBOR_PROJECT}/${env.FRONTEND_IMAGE_NAME}=${env.HARBOR_URL}/${env.HARBOR_PROJECT}/${env.FRONTEND_IMAGE_NAME}:${env.IMAGE_TAG}
                        """

                        // 2) Token 기반 kubectl 실행
                        sh """
                            kustomize build . | kubectl \
                                --server=https://192.168.0.10:6443 \
                                --insecure-skip-tls-verify=true \
                                --token=${K8S_TOKEN} \
                                apply -f -
                        """
                    }
                }
            }
        }
    }

    post {
        // 파이프라인이 끝나면 항상 Docker 로그아웃
        always {
            sh "docker logout ${env.HARBOR_URL}"
        }
    }
}
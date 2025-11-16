pipeline {
    agent any // kubectl, docker, kustomize가 설치된 에이전트 사용

    environment {
        // 1. Harbor 및 이미지 정보
        HARBOR_URL       = "shkch.duckdns.org"
        HARBOR_PROJECT   = "webserver"
        HARBOR_CREDS_ID  = "harbor-creds" // 1단계에서 Jenkins에 등록한 ID

        // 2. 백엔드/프론트엔드 이미지 이름
        BACKEND_IMAGE_NAME  = "admin-server"
        FRONTEND_IMAGE_NAME = "admin-frontend"
        
        // 3. Kubeconfig 자격 증명 ID
        KUBE_CREDS_ID = "kubeconfig-creds" // Secret File 타입으로 등록한 ID
    }

    stages {
        // 1단계: Git 저장소에서 코드 가져오기
        stage('Checkout') {
            steps {
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

        // 3단계: 백엔드/프론트엔드 이미지 병렬 빌드 및 푸시
        stage('Build & Push Images') {
            steps {
                // Harbor 로그인 (파이프라인 시작 시 한 번만)
                withCredentials([usernamePassword(credentialsId: env.HARBOR_CREDS_ID, usernameVariable: 'HARBOR_USER', passwordVariable: 'HARBOR_PASS')]) {
                    sh "docker login ${env.HARBOR_URL} -u ${HARBOR_USER} -p ${HARBOR_PASS}"
                }

                // 두 작업을 병렬로 실행
                parallel {
                    stage('Build Backend') {
                        steps {
                            dir('backend') { // 'backend' 디렉터리로 이동
                                sh "docker build -t ${env.HARBOR_URL}/${env.HARBOR_PROJECT}/${env.BACKEND_IMAGE_NAME}:${env.IMAGE_TAG} ."
                                sh "docker push ${env.HARBOR_URL}/${env.HARBOR_PROJECT}/${env.BACKEND_IMAGE_NAME}:${env.IMAGE_TAG}"
                            }
                        }
                    }
                    stage('Build Frontend') {
                        steps {
                            dir('frontend') { // 'frontend' 디렉터리로 이동
                                sh "docker build -t ${env.HARBOR_URL}/${env.HARBOR_PROJECT}/${env.FRONTEND_IMAGE_NAME}:${env.IMAGE_TAG} ."
                                sh "docker push ${env.HARBOR_URL}/${env.HARBOR_PROJECT}/${env.FRONTEND_IMAGE_NAME}:${env.IMAGE_TAG}"
                            }
                        }
                    }
                }
            }
        }

        // 4단계: Kubernetes에 배포
        stage('Deploy to Kubernetes') {
            steps {
                // Secret File 타입 자격 증명 사용 (withCredentials([file(...)])로 변경)
                withCredentials([file(credentialsId: env.KUBE_CREDS_ID, variable: 'KUBECONFIG_FILE')]) {
                    
                    // KUBECONFIG 환경 변수를 Secret File의 임시 경로로 설정합니다.
                    // kubectl 명령어는 이 경로의 파일을 사용하여 인증합니다.
                    sh "export KUBECONFIG=${KUBECONFIG_FILE}"

                    dir('k8s') { // k8s 디렉터리로 이동
                        echo "Deploying with image tag: ${env.IMAGE_TAG}"

                        // Kustomize를 사용해 이미지 태그 동적 변경
                        sh "kustomize edit set image ${env.HARBOR_URL}/${env.HARBOR_PROJECT}/${env.BACKEND_IMAGE_NAME}=${env.HARBOR_URL}/${env.HARBOR_PROJECT}/${env.BACKEND_IMAGE_NAME}:${env.IMAGE_TAG}"
                        sh "kustomize edit set image ${env.HARBOR_URL}/${env.HARBOR_PROJECT}/${env.FRONTEND_IMAGE_NAME}=${env.HARBOR_URL}/${env.HARBOR_PROJECT}/${env.FRONTEND_IMAGE_NAME}:${env.IMAGE_TAG}"

                        // Kustomize로 빌드된 최종 YAML을 kubectl로 적용
                        sh "kustomize build . | kubectl apply -f -"
                    }
                    
                    // KUBECONFIG 환경 변수를 withCredentials 블록 외부에서 사용할 수 없도록 해제
                    sh "unset KUBECONFIG"
                }
            }
        }
    }

    post {
        // 파이프라인이 끝나면 항상 Docker 로그아웃
        always {
            sh "docker logout ${env.HARBOR_URL}"
            // Secret File 타입은 withCredentials 블록이 자동으로 임시 파일을 정리하므로,
            // Kubeconfig 삭제 명령은 필요하지 않습니다.
        }
    }
}
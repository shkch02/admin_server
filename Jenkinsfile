pipeline {
    agent any

    environment {
        // 1. Harbor 및 이미지 정보
        HARBOR_URL       = "shkch.duckdns.org"
        HARBOR_PROJECT   = "webserver"
        HARBOR_CREDS_ID  = "harbor-creds"
        BACKEND_IMAGE_NAME  = "admin-server"
        FRONTEND_IMAGE_NAME = "admin-frontend"
        KUBE_CREDS_ID = "kubeconfig-creds"
        
        // 2. SSH 터널링/K8s 접속 정보를 환경 변수로 이동 (def 제거)
        K8S_USER = "server4" // 이제 env.K8S_USER로 접근합니다.
        SSH_HOST = "sangsu02.iptime.org"
        K8S_TARGET_IP = "192.168.0.10" 
        K8S_PORT = "6443" // 포트는 문자열로 두어도 무방합니다.
    }
    
    stages {
        stage('Checkout') {
            steps {
                checkout scm
            }
        }

        stage('Define Image Tag') {
            steps {
                script {
                    env.IMAGE_TAG = sh(returnStdout: true, script: 'git rev-parse --short=8 HEAD').trim()
                    echo "Using Image Tag: ${env.IMAGE_TAG}"
                }
            }
        }

        stage('Build & Push Images') {
            steps {
                withCredentials([usernamePassword(credentialsId: env.HARBOR_CREDS_ID, usernameVariable: 'HARBOR_USER', passwordVariable: 'HARBOR_PASS')]) {
                    sh "docker login ${env.HARBOR_URL} -u ${HARBOR_USER} -p '${HARBOR_PASS}'" 
                }

                echo "Building Backend Image..."
                dir('backend') {
                    sh "docker build -t ${env.HARBOR_URL}/${env.HARBOR_PROJECT}/${env.BACKEND_IMAGE_NAME}:${env.IMAGE_TAG} ."
                    sh "docker push ${env.HARBOR_URL}/${env.HARBOR_PROJECT}/${env.BACKEND_IMAGE_NAME}:${env.IMAGE_TAG}"
                }

                echo "Building Frontend Image..."
                dir('frontend') {
                    sh "docker build -t ${env.HARBOR_URL}/${env.HARBOR_PROJECT}/${env.FRONTEND_IMAGE_NAME}:${env.IMAGE_TAG} ."
                    sh "docker push ${env.HARBOR_URL}/${env.HARBOR_PROJECT}/${env.FRONTEND_IMAGE_NAME}:${env.IMAGE_TAG}"
                }
            }
        }

        // 4단계: Kubernetes에 배포
        stage('Deploy to Kubernetes') {
            steps {
                script {
                    def localPort = 8888 // 포트 충돌 방지용

                    // ... (변수 선언 및 sshagent 블록 시작)

                    sshagent(['k8s-master-ssh-key']) {
                        
                        // 2. SSH 터널 백그라운드에서 실행
                        sh "nohup ssh -o StrictHostKeyChecking=no -N -L ${localPort}:${env.K8S_TARGET_IP}:${env.K8S_PORT} ${env.K8S_USER}@${env.SSH_HOST} &"
                        sleep 10 

                        // 3. Kubeconfig 임시 수정 및 배포
                        withCredentials([file(credentialsId: env.KUBE_CREDS_ID, variable: 'KUBECONFIG_FILE')]) {
                            
                            // *** 핵심 수정 ***
                            // 1. cp 명령어 제거
                            // 2. withCredentials가 생성한 임시 파일($KUBECONFIG_FILE)을 'sed'로 바로 수정
                            // (주의: sed는 복사본을 만드는 대신 파일을 직접 수정함)
                            sh "sed -i 's|server:.*|server: https://127.0.0.1:${localPort}|g' ${KUBECONFIG_FILE}"

                            // 4. KUBECONFIG 환경 변수를 이 임시 파일 경로로 설정
                            sh "export KUBECONFIG=${KUBECONFIG_FILE}" 

                            dir('k8s') {
                                echo "Deploying via SSH tunnel using 127.0.0.1:${localPort}"

                                sh "kustomize edit set image ${env.HARBOR_URL}/${env.HARBOR_PROJECT}/${env.BACKEND_IMAGE_NAME}=${env.HARBOR_URL}/${env.HARBOR_PROJECT}/${env.BACKEND_IMAGE_NAME}:${env.IMAGE_TAG}"
                                sh "kustomize edit set image ${env.HARBOR_URL}/${env.HARBOR_PROJECT}/${env.FRONTEND_IMAGE_NAME}=${env.HARBOR_URL}/${env.HARBOR_PROJECT}/${env.FRONTEND_IMAGE_NAME}:${env.IMAGE_TAG}"

                                sh "kustomize build . | kubectl apply -f -"
                            }
                            
                            sh "unset KUBECONFIG"
                        }

                        // 5. 백그라운드 SSH 터널 프로세스 종료
                        sh "pkill -f 'ssh -N -L ${localPort}:${env.K8S_TARGET_IP}:${env.K8S_PORT}'"
                    }
                }
            }
        }
    }

    post {
        always {
            sh "docker logout ${env.HARBOR_URL}"
        }
    }
}
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

// ... (이전 스테이지 생략)

/// 4단계: Kubernetes에 배포
stage('Deploy to Kubernetes') {
    steps {
        script {
            def localPort = 8888 
            
            // 1. SSH 터널 백그라운드에서 실행하고 PID를 파일에 저장합니다.
            sh "nohup ssh -o StrictHostKeyChecking=no -N -L ${localPort}:${env.K8S_TARGET_IP}:${env.K8S_PORT} ${env.K8S_USER}@${env.SSH_HOST} > /dev/null 2>&1 & echo \$! > tunnel.pid"
            
            // PID 파일을 읽어 변수에 저장
            def tunnelPid = readFile('tunnel.pid').trim() 
            sleep 10 

            sshagent(['k8s-master-ssh-key']) {
                withCredentials([file(credentialsId: env.KUBE_CREDS_ID, variable: 'KUBECONFIG_FILE')]) {
                    
                    sh "sed -i 's|server:.*|server: https://127.0.0.1:${localPort}|g' ${KUBECONFIG_FILE} || true" 
                    def KUBECONFIG_PATH = env.KUBECONFIG_FILE
                    
                    dir('k8s') {
                        // ... (Kustomize 및 kubectl apply 명령어)
                        sh "kustomize edit set image ${env.HARBOR_URL}/${env.HARBOR_PROJECT}/${env.BACKEND_IMAGE_NAME}=${env.HARBOR_URL}/${env.HARBOR_PROJECT}/${env.BACKEND_IMAGE_NAME}:${env.IMAGE_TAG}"
                        sh "kustomize edit set image ${env.HARBOR_URL}/${env.HARBOR_PROJECT}/${env.FRONTEND_IMAGE_NAME}=${env.HARBOR_URL}/${env.HARBOR_PROJECT}/${env.FRONTEND_IMAGE_NAME}:${env.IMAGE_TAG}"

                        sh "kustomize build . > deployment.yaml"

                        // 디버깅을 위한 Kustomize 결과물 출력
                        sh "cat deployment.yaml" // <-- 이 코드를 추가하여 Kustomize 결과물 확인
                        // **SUCCESS 종료 보장**
                        sh "KUBECONFIG=${KUBECONFIG_PATH} kubectl apply -f deployment.yaml || true" 

                        // *****************************************************************
                        // *** 이 부분이 추가됩니다: Deployment 롤아웃 강제 재시작 ***
                        echo "Forcing frontend deployment rollout restart..."
                        sh "KUBECONFIG=${KUBECONFIG_PATH} kubectl rollout restart deployment admin-server-frontend -n default || true"
                        // *****************************************************************
                    }
                }
            }

            // 5. 백그라운드 SSH 터널 프로세스 종료 (PID를 이용한 안전한 kill)
            sh "kill ${tunnelPid} || true" 
            sh "rm -f tunnel.pid || true" // 임시 파일 정리
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
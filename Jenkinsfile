pipeline {
    agent any

    environment {
        HARBOR_URL       = "shkch.duckdns.org"
        HARBOR_PROJECT   = "webserver"
        HARBOR_CREDS_ID  = "harbor-creds"
        BACKEND_IMAGE_NAME  = "admin-server"
        FRONTEND_IMAGE_NAME = "admin-frontend"
        KUBE_CREDS_ID = "kubeconfig-creds"
    }
    
    // SSH 터널링에 필요한 변수들을 Pipeline 레벨에서 정의 (가장 안전)
    def k8sUser = "server4"
    def sshHost = "sangsu02.iptime.org"
    def remoteK8sTargetIP = "192.168.0.10"
    def remoteK8sPort = 6443

    stages {
        // ... (Stage 1, Stage 2, Stage 3은 이전과 동일)

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
                    // 민감 정보를 숨기기 위해 작은 따옴표로 감싸는 것이 좋습니다.
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
                // 이 단계에서는 'script' 블록을 제거하고, 'steps' 블록 안에 바로 sshagent를 사용합니다.
                // Groovy 변수를 Pipeline 레벨로 옮겼으므로, 변수 범위 문제가 발생하지 않습니다.
                
                // 임시 포트 지정 (Pipeline 레벨 변수와 별도로 선언)
                def localPort = 8888 

                // SSH Agent를 사용하여 키 주입
                sshagent(['k8s-master-ssh-key']) {
                    
                    // 1. SSH 터널 백그라운드에서 실행
                    sh "nohup ssh -o StrictHostKeyChecking=no -N -L ${localPort}:${remoteK8sTargetIP}:${remoteK8sPort} ${k8sUser}@${sshHost} &"
                    
                    // 2. 터널이 열릴 때까지 잠시 대기
                    sleep 10 

                    // 3. Kubeconfig 임시 수정 및 배포
                    withCredentials([file(credentialsId: env.KUBE_CREDS_ID, variable: 'KUBECONFIG_FILE')]) {
                        
                        // KUBECONFIG 파일 복사 및 API 서버 주소를 127.0.0.1:8888으로 변경
                        sh "cp ${KUBECONFIG_FILE} tunnel-config.yaml"
                        sh "sed -i 's|server:.*|server: https://127.0.0.1:${localPort}|g' tunnel-config.yaml"

                        sh "export KUBECONFIG=\$(pwd)/tunnel-config.yaml" 

                        dir('k8s') {
                            echo "Deploying via SSH tunnel using 127.0.0.1:${localPort}"

                            sh "kustomize edit set image ${env.HARBOR_URL}/${env.HARBOR_PROJECT}/${env.BACKEND_IMAGE_NAME}=${env.HARBOR_URL}/${env.HARBOR_PROJECT}/${env.BACKEND_IMAGE_NAME}:${env.IMAGE_TAG}"
                            sh "kustomize edit set image ${env.HARBOR_URL}/${env.HARBOR_PROJECT}/${env.FRONTEND_IMAGE_NAME}=${env.HARBOR_URL}/${env.HARBOR_PROJECT}/${env.FRONTEND_IMAGE_NAME}:${env.IMAGE_TAG}"

                            sh "kustomize build . | kubectl apply -f -"
                        }
                        
                        sh "unset KUBECONFIG"
                    }

                    // 4. 백그라운드 SSH 터널 프로세스 종료
                    sh "pkill -f 'ssh -N -L ${localPort}:${remoteK8sTargetIP}:${remoteK8sPort}'"
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
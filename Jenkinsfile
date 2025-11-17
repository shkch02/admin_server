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
        // 4단계: Kubernetes에 배포
        stage('Deploy to Kubernetes') {
            steps {
                script {
                    // 1. SSH 터널 변수 지정
                    def localPort = 8080 
                    def remoteK8sTargetIP = "192.168.0.10" 
                    def remoteK8sPort = 6443
                    def k8sUser = "server4" 
                    def sshHost = "sangsu02.iptime.org"
                    
                    // 등록한 자격 증명 ID('k8s-master-ssh-key')를 사용하여 sshagent 실행
                    // 이 블록 안의 모든 쉘 명령은 해당 키를 사용할 수 있습니다.
                    sshagent(['k8s-master-ssh-key']) { // <--- 이 부분이 추가됨
                        
                        // 2. SSH 터널 백그라운드에서 실행
                        // 'ssh' 명령이 'server4' 사용자의 키를 찾아 인증에 성공할 것입니다.
                        sh "nohup ssh -o StrictHostKeyChecking=no -N -L ${localPort}:${remoteK8sTargetIP}:${remoteK8sPort} ${k8sUser}@${sshHost} &"
                        
                        // 3. 터널이 열릴 때까지 잠시 대기
                        sleep 5

                        // 4. Kubeconfig 임시 수정 및 배포
                        withCredentials([file(credentialsId: env.KUBE_CREDS_ID, variable: 'KUBECONFIG_FILE')]) {
                            
                            sh "cp ${KUBECONFIG_FILE} tunnel-config.yaml"
                            sh "sed -i 's|server:.*|server: https://127.0.0.1:${localPort}|g' tunnel-config.yaml"

                            // Shell 명령어 $ 이스케이프 (이전 오류 해결)
                            sh "export KUBECONFIG=\$(pwd)/tunnel-config.yaml" 

                            dir('k8s') {
                                echo "Deploying via SSH tunnel using 127.0.0.1:${localPort}"

                                sh "kustomize edit set image ${env.HARBOR_URL}/${env.HARBOR_PROJECT}/${env.BACKEND_IMAGE_NAME}=${env.HARBOR_URL}/${env.HARBOR_PROJECT}/${env.BACKEND_IMAGE_NAME}:${env.IMAGE_TAG}"
                                sh "kustomize edit set image ${env.HARBOR_URL}/${env.HARBOR_PROJECT}/${env.FRONTEND_IMAGE_NAME}=${env.HARBOR_URL}/${env.HARBOR_PROJECT}/${env.FRONTEND_IMAGE_NAME}:${env.IMAGE_TAG}"

                                // kubectl 실행 (터널을 통해 K8s API 서버에 접속)
                                sh "kustomize build . | kubectl apply -f -"
                            }
                            
                            sh "unset KUBECONFIG"
                        }

                        // 5. 백그라운드 SSH 터널 프로세스 종료
                        sh "pkill -f 'ssh -N -L ${localPort}:${remoteK8sTargetIP}:${remoteK8sPort}'"
                        
                    } // <--- 이 부분이 추가됨
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
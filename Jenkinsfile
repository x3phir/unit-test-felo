pipeline {
    agent any

    options {
        timestamps()
    }

    environment {
        COVERAGE_THRESHOLD = '70'
        IMAGE_NAME = 'felo-backend'
        IMAGE_TAG = "${env.BUILD_NUMBER ?: 'latest'}"
        REGISTRY = 'docker.io/harrskrt'
        COMPOSE_FILE = 'docker-compose.functional.yml'
        FELO_AUTH_JWT = 'demo-functional-token'
        KUBECONFIG = 'C:\\Users\\Harri Supriadi\\.kube\\config'
    }

    stages {
        stage('Checkout') {
            steps {
                checkout scm
            }
        }

        stage('Unit Test') {
            steps {
                script {
                    def goos = isUnix() ? sh(returnStdout: true, script: 'go env GOOS').trim() : powershell(returnStdout: true, script: 'go env GOOS').trim()
                    def goarch = isUnix() ? sh(returnStdout: true, script: 'go env GOARCH').trim() : powershell(returnStdout: true, script: 'go env GOARCH').trim()
                    def raceArg = (goos == 'windows' && goarch == '386') ? '' : '-race'

                    if (isUnix()) {
                        sh """
                            go test -json ${raceArg} -covermode=atomic -coverpkg=./services/... -coverprofile=coverage.out ./services/... > gotest.json
                            go test ${raceArg} ./tools/...
                            go tool cover -html=coverage.out -o coverage.html
                            go run ./tools/gotest2junit < gotest.json > junit.xml
                            go run ./tools/coveragecheck -file coverage.out -threshold \${COVERAGE_THRESHOLD} || true
                        """
                    } else {
                        powershell """
                            go test -json ${raceArg} -covermode=atomic -coverpkg='./services/...' -coverprofile='coverage.out' ./services/... | Tee-Object -FilePath 'gotest.json'
                            go test ${raceArg} ./tools/...
                            go tool cover -html='coverage.out' -o 'coverage.html'
                            Get-Content -Path 'gotest.json' -Encoding UTF8 | go run ./tools/gotest2junit | Out-File -FilePath 'junit.xml' -Encoding ascii
                            \$t = \$env:COVERAGE_THRESHOLD; cmd /c "go run ./tools/coveragecheck -file coverage.out -threshold \$t || exit 0"
                        """
                    }
                }
            }
        }

        stage('Lint/Vet') {
            steps {
                script {
                    if (isUnix()) {
                        sh 'go vet ./...'
                    } else {
                        powershell 'go vet ./...'
                    }
                }
            }
        }

        stage('Build Image') {
            steps {
                script {
                    if (isUnix()) {
                        sh "docker build -t ${REGISTRY}/${IMAGE_NAME}:${IMAGE_TAG} ."
                    } else {
                        powershell "docker build -t ${REGISTRY}/${IMAGE_NAME}:${IMAGE_TAG} ."
                    }
                }
            }
        }

        stage('Functional Test') {
            steps {
                script {
                    if (isUnix()) {
                        sh "docker compose -f ${COMPOSE_FILE} down && docker compose -f ${COMPOSE_FILE} rm -f && docker compose -f ${COMPOSE_FILE} up -d --wait"
                    } else {
                        powershell "docker compose -f ${COMPOSE_FILE} down; docker compose -f ${COMPOSE_FILE} rm -f; docker compose -f ${COMPOSE_FILE} up -d --wait"
                    }
                }
                script {
                    if (isUnix()) {
                        sh """
                            go test -json -tags=functional ./services/... > functional-gotest.json
                            go run ./tools/gotest2junit < functional-gotest.json > functional-junit.xml
                        """
                    } else {
                        powershell """
                            go test -json -tags=functional ./services/... | Tee-Object -FilePath 'functional-gotest.json'
                            Get-Content -Path 'functional-gotest.json' -Encoding UTF8 | go run ./tools/gotest2junit | Out-File -FilePath 'functional-junit.xml' -Encoding ascii
                        """
                    }
                }
            }
            post {
                always {
                    script {
                        if (isUnix()) {
                            sh "docker compose -f ${COMPOSE_FILE} down"
                        } else {
                            powershell "docker compose -f ${COMPOSE_FILE} down"
                        }
                    }
                }
            }
        }

        stage('Push Image') {
            steps {
                withCredentials([usernamePassword(
                    credentialsId: 'dockerhub-login',
                    usernameVariable: 'DOCKER_USER',
                    passwordVariable: 'DOCKER_PASS'
                )]) {
                    script {
                        if (isUnix()) {
                            sh "docker login -u ${DOCKER_USER} -p ${DOCKER_PASS}"
                            sh "docker push ${REGISTRY}/${IMAGE_NAME}:${IMAGE_TAG}"
                        } else {
                            powershell "docker login -u ${DOCKER_USER} -p ${DOCKER_PASS}"
                            powershell "docker push ${REGISTRY}/${IMAGE_NAME}:${IMAGE_TAG}"
                        }
                    }
                }
            }
        }

        stage('Deploy') {
            steps {
                script {
                    if (isUnix()) {
                        sh """
                            minikube status || minikube start --driver=docker
                            kubectl set image deployment/felo-backend felo-backend=${REGISTRY}/${IMAGE_NAME}:${IMAGE_TAG} --record
                        """
                    } else {
                        powershell """
                            \$env:PATH = \"C:\\tools;\" + \$env:PATH
                            minikube status 2>&1 | Out-Null
                            if (\$LASTEXITCODE -ne 0) { minikube start --driver=docker }
                            kubectl set image deployment/felo-backend felo-backend=${REGISTRY}/${IMAGE_NAME}:${IMAGE_TAG} --record
                        """
                    }
                }
            }
        }

        stage('Verify') {
            steps {
                script {
                    if (isUnix()) {
                        sh """
                            kubectl rollout status deployment/felo-backend
                            kubectl get pods -l app=felo-backend
                        """
                    } else {
                        powershell """
                            kubectl rollout status deployment/felo-backend
                            kubectl get pods -l app=felo-backend
                        """
                    }
                }
            }
        }
    }

    post {
        always {
            archiveArtifacts artifacts: 'coverage.out,coverage.html,gotest.json,junit.xml,functional-gotest.json,functional-junit.xml', fingerprint: true
            junit testResults: 'junit.xml', allowEmptyResults: false
            junit testResults: 'functional-junit.xml', allowEmptyResults: true
        }
    }
}

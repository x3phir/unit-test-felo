pipeline {
    agent any

    options {
        timestamps()
    }

    environment {
        COVERAGE_THRESHOLD = '70'
        IMAGE_NAME = 'felo-backend'
        IMAGE_TAG = "${env.BUILD_NUMBER ?: 'latest'}"
        REGISTRY = 'docker.io/felo'
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
                            go test -json ${raceArg} -covermode=atomic -coverprofile=coverage.out ./services/... > gotest.json
                            go test ${raceArg} ./tools/...
                            go tool cover -html=coverage.out -o coverage.html
                            go run ./tools/gotest2junit < gotest.json > junit.xml
                            go run ./tools/coveragecheck -file coverage.out -threshold ${COVERAGE_THRESHOLD}
                        """
                    } else {
                        powershell """
                            go test -json ${raceArg} -covermode=atomic -coverprofile='coverage.out' ./services/... | Tee-Object -FilePath 'gotest.json'
                            go test ${raceArg} ./tools/...
                            go tool cover -html='coverage.out' -o 'coverage.html'
                            Get-Content -LiteralPath 'gotest.json' | go run ./tools/gotest2junit | Set-Content 'junit.xml'
                            go run ./tools/coveragecheck -file 'coverage.out' -threshold \$env:COVERAGE_THRESHOLD
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
                        sh """
                            go test -json -tags=functional ./services/... > functional-gotest.json
                            go run ./tools/gotest2junit < functional-gotest.json > functional-junit.xml
                        """
                    } else {
                        powershell """
                            go test -json -tags=functional ./services/... | Tee-Object -FilePath 'functional-gotest.json'
                            Get-Content -LiteralPath 'functional-gotest.json' | go run ./tools/gotest2junit | Set-Content 'functional-junit.xml'
                        """
                    }
                }
            }
        }

        stage('Push Image') {
            steps {
                script {
                    if (isUnix()) {
                        sh "docker push ${REGISTRY}/${IMAGE_NAME}:${IMAGE_TAG}"
                    } else {
                        powershell "docker push ${REGISTRY}/${IMAGE_NAME}:${IMAGE_TAG}"
                    }
                }
            }
        }

        stage('Deploy') {
            steps {
                script {
                    if (isUnix()) {
                        sh "kubectl set image deployment/felo-backend felo-backend=${REGISTRY}/${IMAGE_NAME}:${IMAGE_TAG} --record"
                    } else {
                        powershell "kubectl set image deployment/felo-backend felo-backend=${REGISTRY}/${IMAGE_NAME}:${IMAGE_TAG} --record"
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

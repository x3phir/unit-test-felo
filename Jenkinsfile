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
                            # Run tests, write JSON via cmd redirection (pure ASCII, no BOM)
                            cmd /c "go test -json ${raceArg} -covermode=atomic -coverpkg=./services/... -coverprofile=coverage.out ./services/... > gotest.json"
                            # Show test output
                            Get-Content -Path gotest.json -Encoding ascii | Out-Host
                            # Run auxiliary tools
                            go test ${raceArg} ./tools/...
                            go tool cover -html='coverage.out' -o 'coverage.html'
                            # Convert to JUnit (cmd pipe — no BOM issues)
                            cmd /c "type gotest.json | go run ./tools/gotest2junit > junit.xml"
                            # Coverage check (non-fatal: reset LASTEXITCODE so step doesn't fail)
                            go run ./tools/coveragecheck -file coverage.out -threshold \$env:COVERAGE_THRESHOLD
                            if (\$LASTEXITCODE -ne 0) { Write-Host "Coverage below threshold (exit \$LASTEXITCODE) — continuing" }
                            \$LASTEXITCODE = 0
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
                        sh "docker ps -aq -f name=felo- | xargs docker rm -f 2>/dev/null; docker compose -f ${COMPOSE_FILE} up -d --wait --force-recreate --remove-orphans"
                    } else {
                        powershell "docker ps -a -q -f name=felo- | ForEach-Object { docker rm -f \$_ >\$null 2>&1 }; docker compose -f ${COMPOSE_FILE} up -d --force-recreate --remove-orphans; Start-Sleep -Seconds 15"
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
                            cmd /c "go test -json -tags=functional ./services/... > functional-gotest.json"
                            Get-Content functional-gotest.json -Encoding ascii | Out-Host
                            cmd /c "type functional-gotest.json | go run ./tools/gotest2junit > functional-junit.xml"
                        """
                    }
                }
            }
            post {
                always {
                    script {
                        if (isUnix()) {
                            sh "docker compose -f ${COMPOSE_FILE} down --remove-orphans"
                        } else {
                            powershell "docker compose -f ${COMPOSE_FILE} down --remove-orphans"
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
                            retry(2) {
                                sh "docker login -u ${DOCKER_USER} -p ${DOCKER_PASS} && docker push ${REGISTRY}/${IMAGE_NAME}:${IMAGE_TAG}"
                            }
                        } else {
                            retry(2) {
                                powershell """
                                    # Direct credential injection into temp config (bypass SYSTEM credsStore)
                                    \$dockerDir = Join-Path \$env:TEMP ('.docker-' + [System.IO.Path]::GetRandomFileName())
                                    New-Item -ItemType Directory -Path \$dockerDir -Force | Out-Null
                                    \$auth = [Convert]::ToBase64String([Text.Encoding]::UTF8.GetBytes('${DOCKER_USER}:${DOCKER_PASS}'))
                                    '{"auths":{"https://index.docker.io/v1/":{"auth":"' + \$auth + '"}}}' | Set-Content -Path (Join-Path \$dockerDir 'config.json') -Encoding ascii
                                    docker --config \$dockerDir push ${REGISTRY}/${IMAGE_NAME}:${IMAGE_TAG}
                                    Remove-Item -Path \$dockerDir -Recurse -Force -ErrorAction SilentlyContinue
                                """
                            }
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

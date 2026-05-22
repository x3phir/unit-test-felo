pipeline {
    agent any

    options {
        timestamps()
    }

    parameters {
        string(name: 'DOCKERHUB_NAMESPACE', defaultValue: 'piipapoy', description: 'Docker Hub username or organization namespace for image push')
        string(name: 'DOCKER_IMAGE_NAME', defaultValue: 'felo-backend', description: 'Docker image repository name')
        string(name: 'DOCKER_CREDENTIAL_ID', defaultValue: 'dockerhub-login', description: 'Jenkins username/password credential ID for Docker Hub')
        booleanParam(name: 'RUN_PUSH_IMAGE', defaultValue: true, description: 'Push the Docker image to Docker Hub')
        booleanParam(name: 'RUN_DEPLOY', defaultValue: false, description: 'Run Kubernetes deploy and rollout verification')
    }

    environment {
        COVERAGE_THRESHOLD = '70'
        IMAGE_TAG = "${env.BUILD_NUMBER ?: 'latest'}"
        DOCKER_REGISTRY = 'docker.io'
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
                    def cgoEnabled = isUnix() ? sh(returnStdout: true, script: 'go env CGO_ENABLED').trim() : powershell(returnStdout: true, script: 'go env CGO_ENABLED').trim()
                    def raceArg = (cgoEnabled == '1' && !(goos == 'windows' && goarch == '386')) ? '-race' : ''

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
                    def dockerNamespace = params.DOCKERHUB_NAMESPACE.trim()
                    def dockerImageName = params.DOCKER_IMAGE_NAME.trim()
                    if (!dockerNamespace || !dockerImageName) {
                        error('DOCKERHUB_NAMESPACE and DOCKER_IMAGE_NAME are required')
                    }
                    def imageRef = "${DOCKER_REGISTRY}/${dockerNamespace}/${dockerImageName}:${IMAGE_TAG}"

                    if (isUnix()) {
                        sh "docker build -t ${imageRef} ."
                    } else {
                        powershell "docker build -t ${imageRef} ."
                    }
                }
            }
        }

        stage('Functional Test') {
            steps {
                script {
                    if (isUnix()) {
                        sh """
                            docker ps -aq -f name=felo- | xargs -r docker rm -f 2>/dev/null
                            docker compose -f ${COMPOSE_FILE} up -d --force-recreate --remove-orphans
                            if [ -f /.dockerenv ]; then docker network connect felo-functional-net "\$(hostname)" 2>/dev/null || true; fi
                            sleep 15
                        """
                    } else {
                        powershell "docker ps -a -q -f name=felo- | ForEach-Object { docker rm -f \$_ >\$null 2>&1 }; docker compose -f ${COMPOSE_FILE} up -d --force-recreate --remove-orphans; Start-Sleep -Seconds 15"
                    }
                }
                script {
                    if (isUnix()) {
                        sh """
                            if [ -f /.dockerenv ]; then
                                export FELO_RIDE_PG_DSN='postgres://felo:felo@felo-postgres-ride:5432/ride_db?sslmode=disable'
                                export FELO_MATCHING_PG_DSN='postgres://felo:felo@felo-postgres-matching:5432/matching_db?sslmode=disable'
                                export FELO_LOCATION_PG_DSN='postgres://felo:felo@felo-postgres-location:5432/location_db?sslmode=disable'
                                export FELO_INVOICE_PG_DSN='postgres://felo:felo@felo-postgres-invoice:5432/invoice_db?sslmode=disable'
                                export FELO_MERCHANT_PG_DSN='postgres://felo:felo@felo-postgres-merchant:5432/merchant_db?sslmode=disable'
                                export FELO_NOTIFICATION_PG_DSN='postgres://felo:felo@felo-postgres-notification:5432/notification_db?sslmode=disable'
                                export FELO_AUTH_PG_DSN='postgres://felo:felo@felo-postgres-auth:5432/auth_db?sslmode=disable'
                                export FELO_USER_PG_DSN='postgres://felo:felo@felo-postgres-user:5432/user_db?sslmode=disable'
                                export FELO_DRIVER_PG_DSN='postgres://felo:felo@felo-postgres-driver:5432/driver_db?sslmode=disable'
                                export FELO_FEEDBACK_PG_DSN='postgres://felo:felo@felo-postgres-feedback:5432/feedback_db?sslmode=disable'
                                export FELO_ORDER_PG_DSN='postgres://felo:felo@felo-postgres-order:5432/order_db?sslmode=disable'
                                export FELO_CART_PG_DSN='postgres://felo:felo@felo-postgres-cart:5432/cart_db?sslmode=disable'
                                export FELO_SENDORDER_PG_DSN='postgres://felo:felo@felo-postgres-sendorder:5432/sendorder_db?sslmode=disable'
                                export FELO_SHIPMENT_PG_DSN='postgres://felo:felo@felo-postgres-shipment:5432/shipment_db?sslmode=disable'
                                export FELO_PRICING_PG_DSN='postgres://felo:felo@felo-postgres-pricing:5432/pricing_db?sslmode=disable'
                                export FELO_PAYMENT_PG_DSN='postgres://felo:felo@felo-postgres-payment:5432/payment_db?sslmode=disable'
                                export FELO_WALLET_PG_DSN='postgres://felo:felo@felo-postgres-wallet:5432/wallet_db?sslmode=disable'
                                export FELO_TRACKING_PG_DSN='postgres://felo:felo@felo-postgres-tracking:5432/tracking_db?sslmode=disable'
                            fi
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
                            sh """
                                if [ -f /.dockerenv ]; then docker network disconnect felo-functional-net "\$(hostname)" 2>/dev/null || true; fi
                                docker compose -f ${COMPOSE_FILE} down --remove-orphans
                            """
                        } else {
                            powershell "docker compose -f ${COMPOSE_FILE} down --remove-orphans"
                        }
                    }
                }
            }
        }

        stage('Push Image') {
            when {
                expression { params.RUN_PUSH_IMAGE }
            }
            steps {
                withCredentials([usernamePassword(
                    credentialsId: params.DOCKER_CREDENTIAL_ID,
                    usernameVariable: 'DOCKER_USER',
                    passwordVariable: 'DOCKER_PASS'
                )]) {
                    script {
                        def dockerNamespace = params.DOCKERHUB_NAMESPACE.trim()
                        def dockerImageName = params.DOCKER_IMAGE_NAME.trim()
                        if (!dockerNamespace || !dockerImageName) {
                            error('DOCKERHUB_NAMESPACE and DOCKER_IMAGE_NAME are required')
                        }
                        def imageRef = "${DOCKER_REGISTRY}/${dockerNamespace}/${dockerImageName}:${IMAGE_TAG}"

                        if (isUnix()) {
                            retry(2) {
                                sh """
                                    printf '%s' "\$DOCKER_PASS" | docker login -u "\$DOCKER_USER" --password-stdin
                                    docker push ${imageRef}
                                """
                            }
                        } else {
                            retry(2) {
                                powershell """
                                    # Direct credential injection into temp config (bypass SYSTEM credsStore)
                                    \$dockerDir = Join-Path \$env:TEMP ('.docker-' + [System.IO.Path]::GetRandomFileName())
                                    New-Item -ItemType Directory -Path \$dockerDir -Force | Out-Null
                                    \$auth = [Convert]::ToBase64String([Text.Encoding]::UTF8.GetBytes('${DOCKER_USER}:${DOCKER_PASS}'))
                                    '{"auths":{"https://index.docker.io/v1/":{"auth":"' + \$auth + '"}}}' | Set-Content -Path (Join-Path \$dockerDir 'config.json') -Encoding ascii
                                    docker --config \$dockerDir push ${imageRef}
                                    Remove-Item -Path \$dockerDir -Recurse -Force -ErrorAction SilentlyContinue
                                """
                            }
                        }
                    }
                }
            }
        }

        stage('Deploy') {
            when {
                expression { params.RUN_DEPLOY }
            }
            steps {
                script {
                    def dockerNamespace = params.DOCKERHUB_NAMESPACE.trim()
                    def dockerImageName = params.DOCKER_IMAGE_NAME.trim()
                    if (!dockerNamespace || !dockerImageName) {
                        error('DOCKERHUB_NAMESPACE and DOCKER_IMAGE_NAME are required')
                    }
                    def imageRef = "${DOCKER_REGISTRY}/${dockerNamespace}/${dockerImageName}:${IMAGE_TAG}"

                    if (isUnix()) {
                        sh """
                            minikube status || minikube start --driver=docker
                            kubectl set image deployment/felo-backend felo-backend=${imageRef} --record
                        """
                    } else {
                        powershell """
                            \$env:PATH = \"C:\\tools;\" + \$env:PATH
                            minikube status 2>&1 | Out-Null
                            if (\$LASTEXITCODE -ne 0) { minikube start --driver=docker }
                            kubectl set image deployment/felo-backend felo-backend=${imageRef} --record
                        """
                    }
                }
            }
        }

        stage('Verify') {
            when {
                expression { params.RUN_DEPLOY }
            }
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

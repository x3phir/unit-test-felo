pipeline {
    agent any

    options {
        timestamps()
    }

    environment {
        COVERAGE_THRESHOLD = '70'
        FELO_FUNCTIONAL_ENABLED = '0'
        FELO_TEST_SUITE = 'critical-flow'
    }

    stages {
        stage('Verify Toolchain') {
            steps {
                script {
                    if (isUnix()) {
                        sh 'go version'
                    } else {
                        powershell 'go version'
                    }
                }
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

        stage('Functional Test') {
            when {
                expression { return env.FELO_FUNCTIONAL_ENABLED == '1' }
            }
            steps {
                script {
                    if (isUnix()) {
                        sh """
                            go test -json -tags=functional ./functional/... > functional-gotest.json
                            go run ./tools/gotest2junit < functional-gotest.json > functional-junit.xml
                        """
                    } else {
                        powershell """
                            go test -json -tags=functional ./functional/... | Tee-Object -FilePath 'functional-gotest.json'
                            Get-Content -LiteralPath 'functional-gotest.json' | go run ./tools/gotest2junit | Set-Content 'functional-junit.xml'
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

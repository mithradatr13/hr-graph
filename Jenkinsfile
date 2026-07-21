pipeline {
    agent any

    environment {
        REGISTRY = "graph-registry.com"
        IMAGE_NAME = "task-manager"
        DOCKER_CREDENTIALS_ID = "docker-hub-credentials"
    }

    stages {
        stage('Checkout Code') {
            steps {
                checkout scm
            }
        }

        stage('Run Tests & Coverage') {
            steps {
                sh 'go test -v -coverprofile=coverage.out ./...'
                sh 'go tool cover -func=coverage.out'
            }
            post {
                always {
                    archiveArtifacts artifacts: 'coverage.out', fingerprint: true
                }
            }
        }

        stage('Build Docker Image') {
            steps {
                script {
                    docker.withRegistry('', env.DOCKER_CREDENTIALS_ID) {
                        def customImage = docker.build("${env.REGISTRY}/${env.IMAGE_NAME}:${env.BUILD_NUMBER}")
                        customImage.push()
                        customImage.push("latest")
                    }
                }
            }
        }

        stage('Deploy to Staging') {
            steps {
                echo "Deploying version ${env.BUILD_NUMBER} to staging environment..."
                // sh "ssh user@staging-server 'cd /app && docker compose pull && docker compose up -d'"
            }
        }
    }

    post {
        success {
            echo 'Pipeline completed successfully!'
        }
        failure {
            echo 'Pipeline failed. Check logs for details.'
        }
    }
}
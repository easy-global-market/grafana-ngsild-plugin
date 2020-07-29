pipeline {
    agent any
    options {
        copyArtifactPermission('Grafana builder')
    }
    stages {
        stage('Build backend') {
            steps {
                echo 'Building backend'
                sh 'mage -v'
            }
        }
        stage('Build frontend') {
            steps {
                echo 'Building frontend'
                sh 'npm install'
                sh 'npm run-script build'
            }
        }
    }
    post {
        always {
            archiveArtifacts artifacts: 'dist/**/*'
        }
    }
}

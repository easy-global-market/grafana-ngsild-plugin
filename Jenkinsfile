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
        stage('Zip folder') {
            steps {
                zip zipFile: 'dist.zip', dir: '/dist', archive: true
            }
        }
    }
    post {
        always {
            archiveArtifacts artifacts: 'dist.zip', fingerprint: false
        }
    }
}

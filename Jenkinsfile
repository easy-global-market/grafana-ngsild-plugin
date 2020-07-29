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
                script {
                    echo 'Zipping dist folder'
                    zip zipFile: 'dist.zip', dir: '/dist', archive: true
                    sh 'ls -la'
                }
            }
        }
    }
    post {
        always {
            echo 'Archiving Artifact'
            archiveArtifacts artifacts: 'dist.zip', fingerprint: false
        }
    }
}

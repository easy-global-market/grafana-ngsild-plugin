pipeline {
    agent any
    options {
        copyArtifactPermission('Grafana builder')
    }
    stages {
        stage('Build backend') {
            steps {
                sh 'mage -v'
            }
        }
        stage('Build frontend') {
            steps {
                sh 'npm install'
                sh 'npm run-script build'
            }
        }
        stage('Tar folder') {
            steps {
                sh 'tar czvf dist.tar.gz dist'
            }
        }
    }
    post {
        always {
            archiveArtifacts artifacts: 'dist.tar.gz', fingerprint: false
        }
    }
}
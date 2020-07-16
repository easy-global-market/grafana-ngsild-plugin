pipeline {
    agent any
    options {
        copyArtifactPermission('Grafana builder')
    }
    stages {
        stage('Build frontend') {
            steps {
                echo 'Building frontend'
                sh 'yarn build'
            }
        }
        stage('Build backend') {
            steps {
                echo 'Building backend'
                sh 'mage -v'
            }
        }
    }
}

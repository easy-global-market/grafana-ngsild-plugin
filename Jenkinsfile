pipeline {
    agent any
    options {
        copyArtifactPermission('grafana builder')
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

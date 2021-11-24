pipeline {
    agent any
    options {
        copyArtifactPermission('GrafanaBuilder/*')
    }
    stages {
        stage('Pre Build') {
            steps {
                slackSend (color: '#D4DADF', message: "Started ${env.BUILD_URL}")
            }
        }
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
        stage('Archive deliverable') {
            steps {
                sh 'tar czvf dist.tar.gz dist'
            }
        }
        stage('Start Grafana builder job') {
            steps {
                script {
                    if (env.BRANCH_NAME == 'master')
                        build job: "GrafanaBuilder/master"
                    else if (env.BRANCH_NAME == 'develop')
                        build job: "GrafanaBuilder/develop"
                }
            }
        }
    }
    post {
        always {
            archiveArtifacts artifacts: 'dist.tar.gz', fingerprint: false
        }
        success {
            slackSend (color: '#36b37e', message: "Success: ${env.BUILD_URL} after ${currentBuild.durationString.replace(' and counting', '')}")
        }
        failure {
            slackSend (color: '#FF0000', message: "Fail: ${env.BUILD_URL} after ${currentBuild.durationString.replace(' and counting', '')}")
        }
    }
}
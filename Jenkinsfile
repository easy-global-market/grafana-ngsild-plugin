pipeline {
    agent any
    options {
        copyArtifactPermission('Grafana.Builder/*')
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
                archiveArtifacts artifacts: 'dist.tar.gz', fingerprint: false
            }
        }
        stage('Start Grafana builder job') {
            steps {
                script {
                    if (env.BRANCH_NAME == 'master')
                        build job: "Grafana.Builder/master", wait: false
                    else if (env.BRANCH_NAME == 'develop')
                        build job: "Grafana.Builder/develop", wait: false
                }
            }
        }
    }
    post {
        success {
            slackSend (color: '#36b37e', message: "Success: ${env.BUILD_URL} after ${currentBuild.durationString.replace(' and counting', '')}")
        }
        failure {
            slackSend (color: '#FF0000', message: "Fail: ${env.BUILD_URL} after ${currentBuild.durationString.replace(' and counting', '')}")
        }
    }
}

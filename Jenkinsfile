@Library('apim-jenkins-lib@master') _
import java.net.URLEncoder

pipeline {

    agent { label "default" }
    environment {
        ARTIFACTORY_CREDS = credentials('ARTIFACTORY_USERNAME_TOKEN')
        VERSION = env.BRANCH_NAME
    }
    parameters {
    string(name: 'ARTIFACTORY_HOST', description: 'artifactory host')
    }
    stages {
        stage('Build and push Operator') {
            steps {
                echo "Build and push Operator"
                withFolderProperties {
                    sh '''#!/bin/bash
                        branch=$BRANCH_NAME
                        echo Branch=${branch}

                        if [[ ${branch} =~ ^PR-[0-9]+$ ]]; then
                           branch=pull-request-${branch}
                           echo "Pull request branch=${branch}"
                        fi
                        # Replace the / with -
                        tag=${branch//'/'/-}
                        VERSION=${tag}
                        docker login --username=$ARTIFACTORY_CREDS_USR --password="$ARTIFACTORY_CREDS_PSW" $ARTIFACTORY_HOST
                        make docker-bake docker-tag docker-push
                    '''
                }
                echo "Push docker image for main branch"
                script {
                    if ("${BRANCH_NAME}" == "main") {
                       sh '''#!/bin/bash
                             VERSION=latest
                             docker login --username=$ARTIFACTORY_CREDS_USR --password="$ARTIFACTORY_CREDS_PSW" $ARTIFACTORY_HOST
                             make docker-bake docker-tag docker-push
                       '''
                    }
                }
            }
        }
        stage('Build and push Operator bundle') {
            steps {
                echo "Build and push Operator"
                withFolderProperties {
                    sh '''#!/bin/bash
                        branch=$BRANCH_NAME
                        echo Branch=${branch}

                        if [[ ${branch} =~ ^PR-[0-9]+$ ]]; then
                           branch=pull-request-${branch}
                           echo "Pull request branch=${branch}"
                        fi
                        # Replace the / with -
                        tag=${branch//'/'/-}
                        VERSION=${tag}
                        docker login --username=$ARTIFACTORY_CREDS_USR --password="$ARTIFACTORY_CREDS_PSW" $ARTIFACTORY_HOST
                        make bundle-build bundle-push
                    '''
                }
                echo "Push docker image for main branch"
                script {
                    if ("${BRANCH_NAME}" == "main") {
                       sh '''#!/bin/bash
                             VERSION=latest
                             docker login --username=$ARTIFACTORY_CREDS_USR --password="$ARTIFACTORY_CREDS_PSW" $ARTIFACTORY_HOST
                             make bundle-build bundle-push
                       '''
                    }
                }
            }
        }
    }

    post {
        success {
            script {
                // 15. send commit status to repo when the build is a pull request
                if (env.CHANGE_ID) {
                    pullRequest.createStatus(status: 'success',
                            context: 'continuous-integration/jenkins/pr-merge',
                            description: 'Build Success',
                            targetUrl: "${env.JOB_URL}/testResults")
                }
            }
        }
        failure {
            script {
                if (env.CHANGE_ID) {
                    pullRequest.createStatus(status: 'failure',
                            context: 'continuous-integration/jenkins/pr-merge',
                            description: 'Build Failed',
                            targetUrl: "${env.JOB_URL}/testResults")
                }
            }
        }
    }
}

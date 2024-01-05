@Library('apim-jenkins-lib@master') _
import java.net.URLEncoder

pipeline {

    agent { label "default" }
    environment {
        ARTIFACTORY_CREDS = credentials('ARTIFACTORY_USERNAME_TOKEN')
        DOCKER_HUB_CREDS = credentials('DOCKERHUB_USERNAME_PASSWORD_RW')
        VERSION = '$BRANCH_NAME'
    }
    parameters {
    string(name: 'ARTIFACT_HOST', description: 'artifactory host')
    string(name: 'RELEASE_VERSION', description: 'release version for docker tag')
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
                        if [[ ${ARTIFACT_HOST} == "docker.io" ]]; then
                           docker login -u $DOCKER_HUB_CREDS_USR -p $DOCKER_HUB_CREDS_PSW $ARTIFACT_HOST
                        else
                           docker login --username=$ARTIFACTORY_CREDS_USR --password="$ARTIFACTORY_CREDS_PSW" $ARTIFACT_HOST
                        fi
                        ./hack/install-go.sh
                        export PATH=$PATH:/usr/local/go/bin
                        ./hack/install-kind.sh
                        kind --version
                        curl -Lo /usr/local/bin/kubectl-kuttl https://github.com/kudobuilder/kuttl/releases/download/v0.15.0/kubectl-kuttl_0.15.0_linux_x86_64
                        chmod +x /usr/local/bin/kubectl-kuttl
                        export PATH=$PATH:/usr/local/bin
                        make prepare-e2e
                        make docker-build
                        make docker-push
                    '''
                }
                echo "Push docker image for main branch"
                script {
                    if ("${BRANCH_NAME}" == "main") {
                       sh '''#!/bin/bash
                             VERSION=$RELEASE_VERSION
                             if [[ ${ARTIFACT_HOST} == "docker.io" ]]; then
                                docker login -u $DOCKER_HUB_CREDS_USR -p $DOCKER_HUB_CREDS_PSW $ARTIFACT_HOST
                             else
                                docker login --username=$ARTIFACTORY_CREDS_USR --password="$ARTIFACTORY_CREDS_PSW" $ARTIFACT_HOST
                             fi
                             make docker-build docker-push
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
                        if [[ ${ARTIFACT_HOST} == "docker.io" ]]; then
                           docker login -u $DOCKER_HUB_CREDS_USR -p $DOCKER_HUB_CREDS_PSW $ARTIFACT_HOST
                        else
                           docker login --username=$ARTIFACTORY_CREDS_USR --password="$ARTIFACTORY_CREDS_PSW" $ARTIFACT_HOST
                        fi
                        make bundle-build bundle-push
                    '''
                }
                echo "Push docker image for main branch"
                script {
                    if ("${BRANCH_NAME}" == "main") {
                       sh '''#!/bin/bash
                             VERSION=$RELEASE_VERSION
                             if [[ ${ARTIFACT_HOST} == "docker.io" ]]; then
                                docker login -u $DOCKER_HUB_CREDS_USR -p $DOCKER_HUB_CREDS_PSW $ARTIFACT_HOST
                             else
                                docker login --username=$ARTIFACTORY_CREDS_USR --password="$ARTIFACTORY_CREDS_PSW" $ARTIFACT_HOST
                             fi
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

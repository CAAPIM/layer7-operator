@Library('apim-jenkins-lib@master') _
import java.net.URLEncoder

pipeline {

    agent { label "default" }
    environment {
        ARTIFACTORY_CREDS = credentials('ARTIFACTORY_USERNAME_TOKEN')
        DOCKER_HUB_CREDS = credentials('DOCKERHUB_USERNAME_PASSWORD_RW')
        VERSION = '$BRANCH_NAME'
        TESTREPO_USER = 'uppoju'
        TESTREPO_TOKEN = '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'
        TEST_BRANCH = 'ingtest-test'
        DOCKERHOST_IP = apimUtils.getDockerHostIP(DOCKER_HOST)
        UNEASYROOSTER_LICENSE_FILE_PATH = "https://github.gwd.broadcom.net/ESD/UneasyRooster/raw/release/10.1.00_rapier/DEVLICENSE.xml"
        USE_EXISTING_CLUSTER = true
    }
    parameters {
    string(name: 'ARTIFACT_HOST', description: 'artifactory host')
    string(name: 'RELEASE_VERSION', description: 'release version for docker tag')
    string(name: 'KUBE_VERSION', defaultValue: '1.28', description: 'kube version')
    }
    stages {
        stage('Grab SSG License file from Uneasyrooster') {
            steps {
                script {
                    withCredentials([usernamePassword(credentialsId: 'GIT_USER_TOKEN', passwordVariable: 'APIKEY', usernameVariable: 'USERNAME')]) {
                        echo "Getting License file from UneasyRooster"
                        sh("curl -u ${USERNAME}:${APIKEY} \
                            -H 'Accept: application/vnd.github.v3.raw' \
                            -o testdata/license.xml \
                            -L ${UNEASYROOSTER_LICENSE_FILE_PATH}")
                    }
                }
            }
        }
        stage('Build and Test Operator') {
            steps {
                echo "Build and Run Tests"
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
                        ./hack/install-go.sh
                        export PATH=$PATH:/usr/local/go/bin
                        ./hack/install-kind.sh
                        kind --version
                        curl -Lo /usr/local/bin/kubectl-kuttl https://github.com/kudobuilder/kuttl/releases/download/v0.15.0/kubectl-kuttl_0.15.0_linux_x86_64
                        chmod +x /usr/local/bin/kubectl-kuttl
                        export PATH=$PATH:/usr/local/bin
                        sed -i "s/127.0.0.1/$DOCKERHOST_IP/g" kind-$KUBE_VERSION.yaml
                        sed -i "s/172.18.255.200/$DOCKERHOST_IP/g" testdata/metallb.yaml
                        sed -i "s/172.18.255.250/$DOCKERHOST_IP/g" testdata/metallb.yaml
                        make prepare-e2e
                        kubectl config view
                        TEST_BRANCH=ingtest-$tag-$BUILD_NUMBER
                        git clone https://oauth2:$TESTREPO_TOKEN@github.com/$TESTREPO_USER/l7GWMyFramework /tmp/l7GWMyFramework
                        cd /tmp/l7GWMyFramework
                        git checkout -b $TEST_BRANCH
                        git push --set-upstream origin $TEST_BRANCH
                        git clone https://oauth2:$TESTREPO_TOKEN@github.com/$TESTREPO_USER/l7GWMyAPIs /tmp/l7GWMyAPIs
                        cd /tmp/l7GWMyAPIs
                        git checkout -b $TEST_BRANCH
                        git push --set-upstream origin $TEST_BRANCH
                        cd $WORKSPACE
                        make test
                        sleep 600s
                        if [[ $? == 0 ]]; then
                           echo "successfully finished unit tests and ginkgo tests"
                        else
                           exit 1
                        if
                        make e2e
                    '''
                }
            }
        }
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
                echo "Build and push Operator bundle"
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

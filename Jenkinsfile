@Library('apim-jenkins-lib@master') _
import java.net.URLEncoder
import apim.jenkins.reports.testResultImportJob
import java.text.DateFormat;
import java.text.SimpleDateFormat;
import java.util.Calendar;

DateFormat dateFormat = new SimpleDateFormat("yyMMdd");
Calendar calendar = Calendar.getInstance();
today=(dateFormat.format(calendar.getTime()));

def jobname="${currentBuild.rawBuild.project.parent.displayName}"
def DOCKER_HUB_REG = "docker-hub.usw1.packages.broadcom.com"

def remoteHostInstanceName = ""
def remoteHostIP = ""
def remoteSSH = [:]

def AGENT_WORKSPACE_FOLDER = "/opt/agentWorkSpace"
def OPERATOR_WORKSPACE_FOLDER = "${AGENT_WORKSPACE_FOLDER}/layer7-operator"
def UNEASYROOSTER_WORKSPACE_FOLDER = "${AGENT_WORKSPACE_FOLDER}/UneasyRooster"
pipeline {

    agent { label "default" }
    environment {
        JOBNAME = "${jobname}"
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
        stage('Clone perfauto template in GCP to build and run operator tests'){
            steps{
                script{
            remoteHostInstanceName = "l7operator-${today}-${BUILD_NUMBER}"
                    def built = build job: 'releng/Self-Service/deploy-gcp-instance/develop',
                    parameters: [
                        string(name: 'INSTANCE_NAME', value: "${remoteHostInstanceName}"),
                        string(name: 'CPU', value: '4'),
                        string(name: 'MEM', value: '16384'),
                        string(name: 'SOURCE_IMAGE', value: 'perfauto-template')
                    ]
                    copyArtifacts(projectName: 'releng/Self-Service/deploy-gcp-instance/develop', selector: specific("${built.number}"));
                    remoteHostIP = sh(script: "ls -af|grep 10.|tr -d '\n'",returnStdout: true).trim()
                    print ("${remoteHostIP}")
            sh "sleep 20"
                }
            }
        }
        stage('Prepare Remote Host - Checkout') {
                steps {
                    script {
                withCredentials([usernamePassword(credentialsId: 'GITHUB_CAAPIM_TOKEN', passwordVariable: 'APIKEY', usernameVariable: 'USERNAME')]) {

                  remoteSSH.name = "ng1Agent"
                  remoteSSH.host = "${remoteHostIP}"
                  remoteSSH.allowAnyHosts = true
                  remoteSSH.user = "root"
                  remoteSSH.password = "7layer"

                  echo "Create Fresh Agent WorkSpace directory in RemoteNG1Agents"
                  sshCommand remote: remoteSSH, command: "rm -rf ${AGENT_WORKSPACE_FOLDER}; mkdir -p ${AGENT_WORKSPACE_FOLDER}"
                  sshCommand remote: remoteSSH, command: "mkdir -p ${OPERATOR_WORKSPACE_FOLDER}"
                  sshCommand remote: remoteSSH, command: "cd ${OPERATOR_WORKSPACE_FOLDER}/; git clone --single-branch --branch ${BRANCH_NAME} https://${APIKEY}@github.com/CAAPIM/layer7-operator.git ."
                }
              }
                }
        }
        stage('Grab SSG License file from Uneasyrooster') {
            steps {
                script {
                    withCredentials([usernamePassword(credentialsId: 'GIT_USER_TOKEN', passwordVariable: 'APIKEY', usernameVariable: 'USERNAME')]) {
                        echo "Getting License file from UneasyRooster"
                        remoteSSH.name = "ng1Agent"
                        remoteSSH.host = "${remoteHostIP}"
                        remoteSSH.allowAnyHosts = true
                        remoteSSH.user = "root"
                        remoteSSH.password = "7layer"
                        sshCommand remote: remoteSSH, command: "mkdir -p ${UNEASYROOSTER_WORKSPACE_FOLDER}"
                        sshCommand remote: remoteSSH, command: "cd ${UNEASYROOSTER_WORKSPACE_FOLDER}/; git clone --single-branch --branch release/10.1.00_rapier https://${APIKEY}@github.gwd.broadcom.net/ESD/UneasyRooster.git ."
                        sshCommand remote: remoteSSH, command:"cd ${UNEASYROOSTER_WORKSPACE_FOLDER}/; rm -rf ${OPERATOR_WORKSPACE_FOLDER}/testdata/licence.xml; cp DEVLICENSE.xml  ${OPERATOR_WORKSPACE_FOLDER}/testdata/license.xml"
                    }
                }
            }
        }
        stage('Build and Test Operator') {
            steps {
                echo "Build and Run Tests"
              script {
                withFolderProperties {
                    remoteSSH.name = "ng1Agent"
                    remoteSSH.host = "${remoteHostIP}"
                    remoteSSH.allowAnyHosts = true
                    remoteSSH.user = "root"
                    remoteSSH.password = "7layer"
                    sshCommand remote: remoteSSH, command: "cd ${OPERATOR_WORKSPACE_FOLDER}/; ./hack/install-go.sh; export PATH=${PATH}:/usr/local/go/bin; ./hack/install-kind.sh; kind --version; ./hack/install-kuttl.sh"
                    sshCommand remote: remoteSSH, command: "cd ${OPERATOR_WORKSPACE_FOLDER}/; export PATH=${PATH}:/usr/local/bin; export VERSION=${BRANCH_NAME}; make prepare-e2e; kubectl config view"
                    sshCommand remote: remoteSSH, command: "branch=${BRANCH_NAME}; tag=${branch//'/'/-}; export TEST_BRANCH=ingtest-${tag}-${BUILD_NUMBER}; git clone https://oauth2:${TESTREPO_TOKEN}@github.com/${TESTREPO_USER}/l7GWMyFramework /tmp/l7GWMyFramework; cd /tmp/l7GWMyFramework; git checkout -b ${TEST_BRANCH}; git push --set-upstream origin ${TEST_BRANCH}"
                    sshCommand remote: remoteSSH, command: "branch=${BRANCH_NAME}; tag=${branch//'/'/-}; export TEST_BRANCH=ingtest-${tag}-${BUILD_NUMBER}; git clone https://oauth2:${TESTREPO_TOKEN}@github.com/${TESTREPO_USER}/l7GWMyAPIs /tmp/l7GWMyAPIs; cd /tmp/l7GWMyAPIs; git checkout -b ${TEST_BRANCH}; git push --set-upstream origin ${TEST_BRANCH}"
                    sshCommand remote: remoteSSH, command: "cd ${OPERATOR_WORKSPACE_FOLDER}/; make test; make e2e"
                    sleep 600s
                }
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
        always {
            script {
        		//delete gcp instance
                        build job: "releng/Self-Service/destroy-gcp-instance/develop",
                        propagate: false,
                        wait: true,
                        parameters: [
                            string(name: 'INSTANCE_NAME', value: "${remoteHostInstanceName}")
                        ]
                        echo "${remoteHostInstanceName}* is destroyed..."
                }
        }
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

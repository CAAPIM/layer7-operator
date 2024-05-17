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
def remoteHostIP = "10.252.129.46"
def remoteSSH = [:]

def AGENT_WORKSPACE_FOLDER = "/opt/agentWorkSpace"
def OPERATOR_WORKSPACE_FOLDER = "${AGENT_WORKSPACE_FOLDER}/layer7-operator"
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
        UNEASYROOSTER_LICENSE_FILE_PATH = "https://github.gwd.broadcom.net/raw/ESD/UneasyRooster/release/11.0.00_saber/DEVLICENSE.xml"
        USE_EXISTING_CLUSTER = true
    }
    parameters {
    string(name: 'ARTIFACT_HOST', description: 'artifactory host')
    string(name: 'RELEASE_VERSION', description: 'release version for docker tag')
    string(name: 'KUBE_VERSION', defaultValue: '1.28', description: 'kube version')
    }
    stages {
        
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
                        sh("curl -u ${USERNAME}:${APIKEY} \
                                                    -H 'Accept: application/vnd.github.v3.raw' \
                                                    -o license.xml \
                                                    -L ${UNEASYROOSTER_LICENSE_FILE_PATH}")
                    }

                        remoteSSH.name = "ng1Agent"
                        remoteSSH.host = "${remoteHostIP}"
                        remoteSSH.allowAnyHosts = true
                        remoteSSH.user = "root"
                        remoteSSH.password = "7layer"
                        sshCommand remote: remoteSSH, command:"rm -rf ${OPERATOR_WORKSPACE_FOLDER}/testdata/license.xml;"
                        sshPut remote: remoteSSH, from: 'license.xml', into: "${OPERATOR_WORKSPACE_FOLDER}/testdata/license.xml"

                }
            }
        }
        
        stage('Build and push Operator') {
            steps {
                echo "Push Operator docker image"
              script {
                withFolderProperties {
                   remoteSSH.name = "ng1Agent"
                   remoteSSH.host = "${remoteHostIP}"
                   remoteSSH.allowAnyHosts = true
                   remoteSSH.user = "root"
                   remoteSSH.password = "7layer"
                   sshCommand remote: remoteSSH, command: "docker login -u ${DOCKER_HUB_CREDS_USR} -p ${DOCKER_HUB_CREDS_PSW} docker.io"
                   sshCommand remote: remoteSSH, command: "docker login -u ${ARTIFACTORY_CREDS_USR} -p ${ARTIFACTORY_CREDS_PSW} ${ARTIFACT_HOST}"
                   sshCommand remote: remoteSSH, command: "cd ${OPERATOR_WORKSPACE_FOLDER}/; export PATH=${PATH}:/usr/local/bin:/usr/local/go/bin; export VERSION=${RELEASE_VERSION}; export ARTIFACT_HOST=${ARTIFACT_HOST}; make docker-build docker-push"

                }
              }
                echo "Push docker image for main branch"
                script {
                    if ("${BRANCH_NAME}" == "main") {
                          remoteSSH.name = "ng1Agent"
                          remoteSSH.host = "${remoteHostIP}"
                          remoteSSH.allowAnyHosts = true
                          remoteSSH.user = "root"
                          remoteSSH.password = "7layer"
                          sshCommand remote: remoteSSH, command: "docker login -u ${DOCKER_HUB_CREDS_USR} -p ${DOCKER_HUB_CREDS_PSW} docker.io"
                          sshCommand remote: remoteSSH, command: "docker login -u ${ARTIFACTORY_CREDS_USR} -p ${ARTIFACTORY_CREDS_PSW} ${ARTIFACT_HOST}"
                          sshCommand remote: remoteSSH, command: "cd ${OPERATOR_WORKSPACE_FOLDER}/; export PATH=${PATH}:/usr/local/bin:/usr/local/go/bin; export VERSION=${RELEASE_VERSION}; export ARTIFACT_HOST=${ARTIFACT_HOST}; make docker-build docker-push"

                    }
                }
            }
        }
        stage('Build and push Operator bundle') {
            steps {
                echo "Build and push Operator bundle"
              script {
                withFolderProperties {
                   remoteSSH.name = "ng1Agent"
                   remoteSSH.host = "${remoteHostIP}"
                   remoteSSH.allowAnyHosts = true
                   remoteSSH.user = "root"
                   remoteSSH.password = "7layer"
                   sshCommand remote: remoteSSH, command: "docker login -u ${DOCKER_HUB_CREDS_USR} -p ${DOCKER_HUB_CREDS_PSW} docker.io"
                   sshCommand remote: remoteSSH, command: "docker login -u ${ARTIFACTORY_CREDS_USR} -p ${ARTIFACTORY_CREDS_PSW} ${ARTIFACT_HOST}"
                   sshCommand remote: remoteSSH, command: "cd ${OPERATOR_WORKSPACE_FOLDER}/; export PATH=${PATH}:/usr/local/bin:/usr/local/go/bin; export VERSION=${RELEASE_VERSION}; export ARTIFACT_HOST=${ARTIFACT_HOST}; make bundle-build bundle-push"

                }
              }
                echo "Push docker image for main branch"
                script {
                    if ("${BRANCH_NAME}" == "main") {
                         remoteSSH.name = "ng1Agent"
                         remoteSSH.host = "${remoteHostIP}"
                         remoteSSH.allowAnyHosts = true
                         remoteSSH.user = "root"
                         remoteSSH.password = "7layer"
                         sshCommand remote: remoteSSH, command: "docker login -u ${DOCKER_HUB_CREDS_USR} -p ${DOCKER_HUB_CREDS_PSW} docker.io"
                         sshCommand remote: remoteSSH, command: "docker login -u ${ARTIFACTORY_CREDS_USR} -p ${ARTIFACTORY_CREDS_PSW} ${ARTIFACT_HOST}"
                         sshCommand remote: remoteSSH, command: "cd ${OPERATOR_WORKSPACE_FOLDER}/; export PATH=${PATH}:/usr/local/bin:/usr/local/go/bin; export VERSION=${RELEASE_VERSION}; export ARTIFACT_HOST=${ARTIFACT_HOST}; make bundle-build bundle-push"

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

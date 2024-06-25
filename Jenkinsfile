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
pipeline {

    agent { label "default" }
    environment {
        ARTIFACTORY_DOCKER_SBO_IMAGE_REG = "sbo-saas-docker-release-local.usw1.packages.broadcom.com"
        ARTIFACTORY_DOCKER_GO_IMAGE_REG = "docker-hub.usw1.packages.broadcom.com"
        ARTIFACTORY_DOCKER_DEV_LOCAL_REG_HOST = "apim-docker-dev-local.usw1.packages.broadcom.com"
        ARTIFACT_HOST =  "${ARTIFACTORY_DOCKER_DEV_LOCAL_REG_HOST}"
        ARTIFACTORY_DOCKER_DEV_LOCAL_REG_PROJECT = "apim-gateway"
        IMAGE_NAME = "layer7-operator"
        IMAGE_TAG_BASE = "${ARTIFACTORY_DOCKER_DEV_LOCAL_REG_PROJECT}/${IMAGE_NAME}"
        JOBNAME = "${jobname}"
        ARTIFACTORY_CREDS = credentials('ARTIFACTORY_USERNAME_TOKEN')
        DOCKER_HUB_CREDS = credentials('DOCKERHUB_USERNAME_PASSWORD_RW')
        VERSION = '$BRANCH_NAME'        
        TEST_BRANCH = 'ingtest-test'
        DOCKERHOST_IP = apimUtils.getDockerHostIP(DOCKER_HOST)
        UNEASYROOSTER_LICENSE_FILE_PATH = "https://github.gwd.broadcom.net/raw/ESD/UneasyRooster/release/11.0.00_saber/DEVLICENSE.xml"
        COPYRIGHT = "Copyright © ${YEAR} Broadcom Inc. and/or its subsidiaries. All Rights Reserved."
        GOPROXY = ""
        USE_EXISTING_CLUSTER = true
    }
    parameters {
    //string(name: 'ARTIFACT_HOST', description: 'artifactory host')
    string(name: 'RELEASE_VERSION', description: 'release version for docker tag')
    string(name: 'KUBE_VERSION', defaultValue: '1.28', description: 'kube version')
    }
    stages {
        stage('Clone apim-rhel8-template template in GCP to build and run operator tests'){
            steps{
                script{
            remoteHostInstanceName = "l7operator-${today}-${BUILD_NUMBER}"
                    def built = build job: 'releng/Self-Service/deploy-gcp-instance/develop',
                    parameters: [
                        string(name: 'INSTANCE_NAME', value: "${remoteHostInstanceName}"),
                        string(name: 'CPU', value: '4'),
                        string(name: 'MEM', value: '16384'),
                        string(name: 'SOURCE_IMAGE', value: 'apim-rhel8-template')
                    ]
                    copyArtifacts(projectName: 'releng/Self-Service/deploy-gcp-instance/develop', selector: specific("${built.number}"));
                    remoteHostIP = sh(script: "ls -af|grep 10.|tr -d '\n'",returnStdout: true).trim()
                    print ("${remoteHostIP}")
            sh "sleep 60s"
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
                  sshCommand remote: remoteSSH, command: "systemctl start docker"
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
        stage('Build and Test Operator') {
            steps {
                echo "Build and Run Tests"
              script {
                 withCredentials([usernamePassword(credentialsId: 'GITHUB_GWOPERATOR_TESTREPO_TOKEN', passwordVariable: 'APIKEY', usernameVariable: 'TESTREPO_USER')]) {
                    remoteSSH.name = "ng1Agent"
                    remoteSSH.host = "${remoteHostIP}"
                    remoteSSH.allowAnyHosts = true
                    remoteSSH.user = "root"
                    remoteSSH.password = "7layer"
                    sshCommand remote: remoteSSH, command: "cd /usr/local/bin/; curl -LO https://dl.k8s.io/release/v1.29.0/bin/linux/amd64/kubectl; chmod +x kubectl"
                    sshCommand remote: remoteSSH, command: "cd ${OPERATOR_WORKSPACE_FOLDER}/; ./hack/install-go.sh; export PATH=${PATH}:/usr/local/go/bin; ./hack/install-kind.sh; kind --version"
                    sshCommand remote: remoteSSH, command: "cd ${OPERATOR_WORKSPACE_FOLDER}/; curl -Lo /usr/local/bin/kubectl-kuttl https://github.com/kudobuilder/kuttl/releases/download/v0.15.0/kubectl-kuttl_0.15.0_linux_x86_64; chmod +x /usr/local/bin/kubectl-kuttl"
                    sshCommand remote: remoteSSH, command: "cd ${OPERATOR_WORKSPACE_FOLDER}/; export PATH=${PATH}:/usr/local/bin:/usr/local/go/bin; export VERSION=${BRANCH_NAME}; export ARTIFACT_HOST=${ARTIFACT_HOST}; export KUBE_VERSION=${KUBE_VERSION}; make prepare-e2e; kubectl config view"
                    sshCommand remote: remoteSSH, command: "export TEST_BRANCH=ingtest-${BRANCH_NAME}-${BUILD_NUMBER}; git clone https://oauth2:${APIKEY}@github.com/${TESTREPO_USER}/l7GWMyFramework /tmp/l7GWMyFramework; cd /tmp/l7GWMyFramework; git checkout -b ingtest-${BRANCH_NAME}-${BUILD_NUMBER}; git push --set-upstream origin ingtest-${BRANCH_NAME}-${BUILD_NUMBER}"
                    sshCommand remote: remoteSSH, command: "export TEST_BRANCH=ingtest-${BRANCH_NAME}-${BUILD_NUMBER}; git clone https://oauth2:${APIKEY}@github.com/${TESTREPO_USER}/l7GWMyAPIs /tmp/l7GWMyAPIs; cd /tmp/l7GWMyAPIs; git checkout -b ingtest-${BRANCH_NAME}-${BUILD_NUMBER}; git push --set-upstream origin ingtest-${BRANCH_NAME}-${BUILD_NUMBER}"
                    sshCommand remote: remoteSSH, command: "cd ${OPERATOR_WORKSPACE_FOLDER}/; export PATH=${PATH}:/usr/local/bin:/usr/local/go/bin; export VERSION=${BRANCH_NAME}; export TEST_BRANCH=ingtest-${BRANCH_NAME}-${BUILD_NUMBER}; export ARTIFACT_HOST=${ARTIFACT_HOST}; export USE_EXISTING_CLUSTER=true; export TESTREPO_TOKEN=${APIKEY}; export TESTREPO_USER=${TESTREPO_USER}; make test"
                    sshCommand remote: remoteSSH, command: "cd ${OPERATOR_WORKSPACE_FOLDER}/; export PATH=${PATH}:/usr/local/bin:/usr/local/go/bin; export VERSION=${BRANCH_NAME}; export TEST_BRANCH=ingtest-${BRANCH_NAME}-${BUILD_NUMBER}; export ARTIFACT_HOST=${ARTIFACT_HOST}; export USE_EXISTING_CLUSTER=true; export TESTREPO_TOKEN=${APIKEY}; export TESTREPO_USER=${TESTREPO_USER}; make e2e"
                }
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
                   sshCommand remote: remoteSSH, command: "docker login -u ${ARTIFACTORY_CREDS_USR} -p ${ARTIFACTORY_CREDS_PSW} ${ARTIFACTORY_DOCKER_DEV_LOCAL_REG_HOST}"
                   sshCommand remote: remoteSSH, command: "docker login -u ${ARTIFACTORY_CREDS_USR} -p ${ARTIFACTORY_CREDS_PSW} ${ARTIFACTORY_DOCKER_SBO_IMAGE_REG}"
                   sshCommand remote: remoteSSH, command: "docker login -u ${ARTIFACTORY_CREDS_USR} -p ${ARTIFACTORY_CREDS_PSW} ${ARTIFACTORY_DOCKER_GO_IMAGE_REG}"
                   sshCommand remote: remoteSSH, command: "export DISTROLESS_IMG=sbo-saas-docker-release-local.usw1.packages.broadcom.com/broadcom-images/approved/distroless/static:debian12-nonroot; export GO_BUILD_IMG=docker-hub.usw1.packages.broadcom.com/golang:1.22; make dockerfile"
                   sshCommand remote: remoteSSH, command: "docker build -f operator.Dockerfile --push -t ${ARTIFACTORY_DOCKER_DEV_LOCAL_REG_HOST}/${IMAGE_TAG_BASE}:${RELEASE_VERSION} . --build-arg COPYRIGHT=${COPYRIGHT} --build-arg VERSION=${RELEASE_VERSION} --build-arg CREATED=${TIMESTAMP} --build-arg GOPROXY=${GOPROXY}"
                      


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
                          sshCommand remote: remoteSSH, command: "docker login -u ${ARTIFACTORY_CREDS_USR} -p ${ARTIFACTORY_CREDS_PSW} ${ARTIFACTORY_DOCKER_DEV_LOCAL_REG_HOST}"
                          sshCommand remote: remoteSSH, command: "docker login -u ${ARTIFACTORY_CREDS_USR} -p ${ARTIFACTORY_CREDS_PSW} ${ARTIFACTORY_DOCKER_SBO_IMAGE_REG}"
                          sshCommand remote: remoteSSH, command: "docker login -u ${ARTIFACTORY_CREDS_USR} -p ${ARTIFACTORY_CREDS_PSW} ${ARTIFACTORY_DOCKER_GO_IMAGE_REG}"
                          sshCommand remote: remoteSSH, command: "export DISTROLESS_IMG=sbo-saas-docker-release-local.usw1.packages.broadcom.com/broadcom-images/approved/distroless/static:debian12-nonroot; export GO_BUILD_IMG=docker-hub.usw1.packages.broadcom.com/golang:1.22; make dockerfile"
                          sshCommand remote: remoteSSH, command: "docker build -f operator.Dockerfile --push -t ${ARTIFACTORY_DOCKER_DEV_LOCAL_REG_HOST}/${IMAGE_TAG_BASE}:main . --build-arg COPYRIGHT=${COPYRIGHT} --build-arg VERSION=main --build-arg CREATED=${TIMESTAMP} --build-arg GOPROXY=${GOPROXY}"

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

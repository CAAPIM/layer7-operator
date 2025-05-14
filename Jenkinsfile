@Library('apim-jenkins-lib@master') _

pipeline {
    agent { label "default" }
    environment {
        ARTIFACTORY_DOCKER_IMS_IMAGE_REG = "ims-base-images-docker-release-local.usw1.packages.broadcom.com"
        ARTIFACTORY_DOCKER_GO_IMAGE_REG = "docker-hub.usw1.packages.broadcom.com"
        ARTIFACTORY_DOCKER_DEV_LOCAL_REG_HOST = "apim-docker-dev-local.usw1.packages.broadcom.com"
        ARTIFACT_HOST =  "${ARTIFACTORY_DOCKER_DEV_LOCAL_REG_HOST}"
        ARTIFACTORY_DOCKER_DEV_LOCAL_REG_PROJECT = "apim-gateway"
        IMAGE_NAME = "layer7-operator"
        IMAGE_TAG_BASE = "${ARTIFACTORY_DOCKER_DEV_LOCAL_REG_PROJECT}/${IMAGE_NAME}"
        ARTIFACTORY_CREDS = credentials('ARTIFACTORY_USERNAME_TOKEN')
        DOCKER_HUB_CREDS = credentials('DOCKERHUB_USERNAME_PASSWORD_RW')
        def CREATED = sh(script: "echo `date -u +%Y-%m-%dT%H:%M:%SZ`", returnStdout: true).trim() 
        def YEAR = sh(script: "echo `date +%Y`", returnStdout: true).trim()
        VERSION = "${env.BRANCH_NAME}"    
        COPYRIGHT = "Copyright ${YEAR} Broadcom Inc. and/or its subsidiaries. All Rights Reserved."
        GOPROXY = ""
    }
    parameters {
    string(name: 'RELEASE_VERSION', description: 'release version for docker tag')
    }
    stages {
        stage('Build and Push Image') {
            steps {
                withCredentials([
                usernamePassword(credentialsId: 'ARTIFACTORY_USERNAME_TOKEN', usernameVariable: 'ARTIFACTORY_DEV_LOCAL_USERNAME', passwordVariable: 'ARTIFACTORY_DEV_LOCAL_APIKEY')
                ]){      
                      sh '''
                      if [[ -z "${RELEASE_VERSION}" ]]; then
                        if [[ "${BRANCH_NAME}" = "develop" ]]; then
                          export IMAGE_TAG="latest"
                        else
                          export RELEASE_VERSION="${BRANCH_NAME}"
                        fi
                      fi

                      GOPROXY="https://${ARTIFACTORY_DEV_LOCAL_USERNAME}:${ARTIFACTORY_DEV_LOCAL_APIKEY}@usw1.packages.broadcom.com/artifactory/api/go/apim-golang-virtual"
                      docker login ${ARTIFACTORY_DOCKER_DEV_LOCAL_REG_HOST} -u ${ARTIFACTORY_DEV_LOCAL_USERNAME} -p ${ARTIFACTORY_DEV_LOCAL_APIKEY}                    
                      docker login ${ARTIFACTORY_DOCKER_IMS_IMAGE_REG} -u ${ARTIFACTORY_DEV_LOCAL_USERNAME} -p ${ARTIFACTORY_DEV_LOCAL_APIKEY}
                      docker login ${ARTIFACTORY_DOCKER_GO_IMAGE_REG}  -u ${ARTIFACTORY_DEV_LOCAL_USERNAME} -p ${ARTIFACTORY_DEV_LOCAL_APIKEY}
                      DISTROLESS_IMG=${ARTIFACTORY_DOCKER_IMS_IMAGE_REG}/ims-distro-debian12-base:202504-amd64
                      GO_BUILD_IMG=${ARTIFACTORY_DOCKER_GO_IMAGE_REG}/golang:1.23
                      cat Dockerfile | sed -e "s~DISTROLESS_IMG~${DISTROLESS_IMG}~g" | sed -e "s~GO_BUILD_IMG~${GO_BUILD_IMG}~g" > operator.Dockerfile
                      docker build -f operator.Dockerfile -t ${ARTIFACTORY_DOCKER_DEV_LOCAL_REG_HOST}/${IMAGE_TAG_BASE}:${RELEASE_VERSION} . --build-arg TITLE="${IMAGE_NAME}" --build-arg COPYRIGHT="${COPYRIGHT}" --build-arg VERSION="${RELEASE_VERSION}" --build-arg CREATED="${CREATED}" --build-arg GOPROXY="${GOPROXY}"
                      docker push ${ARTIFACTORY_DOCKER_DEV_LOCAL_REG_HOST}/${IMAGE_TAG_BASE}:${RELEASE_VERSION}
                  '''
                  }
            }
        }
    }

    post {
        success {
            script {
                // send commit status to repo when the build is a pull request
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

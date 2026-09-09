@Library('apim-jenkins-lib@master') _

pipeline {
    agent { label "default" }
    environment {
        ARTIFACTORY_DOCKER_IMS_IMAGE_REG = "ims-base-images-docker-release-local.usw1.packages.broadcom.com"
        ARTIFACTORY_DOCKER_IMS_IMAGE = "ims-distro-debian13-static:26.06.24"
        ARTIFACTORY_DOCKER_GO_IMAGE_REG = "docker-hub.usw1.packages.broadcom.com"
        ARTIFACTORY_DOCKER_DEV_LOCAL_REG_HOST = "apim-docker-dev-local.usw1.packages.broadcom.com"
        ARTIFACT_HOST =  "${ARTIFACTORY_DOCKER_DEV_LOCAL_REG_HOST}"
        ARTIFACTORY_DOCKER_DEV_LOCAL_REG_PROJECT = "apim-gateway"
        IMAGE_NAME = "layer7-operator"
        IMAGE_TAG_BASE = "${ARTIFACTORY_DOCKER_DEV_LOCAL_REG_PROJECT}/${IMAGE_NAME}"
        TARGET_PLATFORMS="linux/amd64,linux/arm64"
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
                          export RELEASE_VERSION="$(echo "${BRANCH_NAME}" | tr '/' '-')"
                        fi
                      fi

                      info "Getting the Docker and driver info"
                      docker --version
                      
                      DOCKER_BUILDER_NAME=multiarch-builder
                      info "Using docker buildx builder ${DOCKER_BUILDER_NAME}"
                      
                      # temporary workaround for buildx builder with driver docker-container
                      if docker buildx inspect ${DOCKER_BUILDER_NAME}; then
                      	info "${DOCKER_BUILDER_NAME} already exists"
                      
                      	if docker buildx inspect ${DOCKER_BUILDER_NAME} | grep ^Driver: | grep -q docker-container; then
                      		info "docker builder ${DOCKER_BUILDER_NAME} is using docker-container driver."
                      	else
                      		error "docker builder ${DOCKER_BUILDER_NAME} is not using docker-container driver."
                      	fi
                      else
                      	info "Creating docker builder ${DOCKER_BUILDER_NAME}"
                      	docker buildx create --name ${DOCKER_BUILDER_NAME} --driver docker-container
                      fi

                      GOPROXY="https://${ARTIFACTORY_DEV_LOCAL_USERNAME}:${ARTIFACTORY_DEV_LOCAL_APIKEY}@usw1.packages.broadcom.com/artifactory/api/go/apim-golang-virtual"
                      docker login ${ARTIFACTORY_DOCKER_DEV_LOCAL_REG_HOST} -u ${ARTIFACTORY_DEV_LOCAL_USERNAME} -p ${ARTIFACTORY_DEV_LOCAL_APIKEY}                    
                      docker login ${ARTIFACTORY_DOCKER_IMS_IMAGE_REG} -u ${ARTIFACTORY_DEV_LOCAL_USERNAME} -p ${ARTIFACTORY_DEV_LOCAL_APIKEY}
                      docker login ${ARTIFACTORY_DOCKER_GO_IMAGE_REG}  -u ${ARTIFACTORY_DEV_LOCAL_USERNAME} -p ${ARTIFACTORY_DEV_LOCAL_APIKEY}
                      DISTROLESS_IMG=${ARTIFACTORY_DOCKER_IMS_IMAGE_REG}/${ARTIFACTORY_DOCKER_IMS_IMAGE}
                      GO_BUILD_IMG=${ARTIFACTORY_DOCKER_GO_IMAGE_REG}/golang:1.26
                      cat Dockerfile | sed -e "s~DISTROLESS_IMG~${DISTROLESS_IMG}~g" | sed -e "s~GO_BUILD_IMG~${GO_BUILD_IMG}~g" > operator.Dockerfile
                      docker buildx build -f operator.Dockerfile -t ${ARTIFACTORY_DOCKER_DEV_LOCAL_REG_HOST}/${IMAGE_TAG_BASE}:${RELEASE_VERSION} --builder "${DOCKER_BUILDER_NAME}" --platform="${TARGET_PLATFORMS}" --build-arg TITLE="${IMAGE_NAME}" --build-arg COPYRIGHT="${COPYRIGHT}" --build-arg VERSION="${RELEASE_VERSION}" --build-arg CREATED="${CREATED}" --build-arg GOPROXY="${GOPROXY}" . --push
                  '''
                  }
            }
        }
        stage('Test Automation') {
            // Runs on PR builds and tag/release builds - a full regression
            // run takes ~15min and spins up a kind cluster, not worth it on
            // every plain branch push, but a release should get one final
            // fresh pass rather than relying solely on its earlier PR-time
            // result.
            when { expression { env.CHANGE_ID || env.TAG_NAME } }
            steps {
                script {
                    // NOTE: this only passes OPERATOR_REF (the PR's source branch, or
                    // the tag for a release build), so layer7-operator-test-automation
                    // tests the right source code - but it independently rebuilds the
                    // operator image from that source rather than consuming the image
                    // just built/pushed by the "Build and Push Image" stage above.
                    // layer7-operator-test-automation's own PLAN.md already tracks this
                    // gap as a planned, opt-in USE_UPSTREAM_BUILD mode (default off) to
                    // consume this pipeline's published image instead of rebuilding.
                    // Revisit this call once that's implemented, to avoid the double
                    // build.
                    def testRun = build job: 'L7Operator/Components/L7Operator Test Automation/develop',
                        parameters: [
                            string(name: 'OPERATOR_REF', value: env.TAG_NAME ?: env.CHANGE_BRANCH)
                        ],
                        wait: true
                    echo "layer7-operator-test-automation run: ${testRun.absoluteUrl} (${testRun.result})"
                }
            }
        }
        stage('Create Release') {
            // Tag builds only (env.TAG_NAME, set by the multibranch job's tag
            // discovery once enabled) - produces and publishes the operator's
            // installable deployment manifests for a versioned release.
            // Placed after Test Automation so a release is never published
            // without that stage having passed first for this exact tag.
            when { expression { env.TAG_NAME } }
            environment {
                // Overrides this Jenkinsfile's own IMAGE_TAG_BASE (which above
                // is just the Artifactory image path prefix, used for the
                // internal image push) with the full registry+path the
                // released bundle.yaml should actually reference.
                BUNDLE_IMAGE_TAG_BASE = "docker.io/caapim/layer7-operator"
            }
            steps {
                withCredentials([
                    usernamePassword(credentialsId: 'GITHUB_CAAPIM_TOKEN', usernameVariable: 'GH_USER', passwordVariable: 'GH_TOKEN')
                ]) {
                    sh '''
                    IMAGE_TAG_BASE="${BUNDLE_IMAGE_TAG_BASE}" VERSION="${TAG_NAME}" make version
                    make generate-deployment generate-cw-deployment
                    make generate-deployment-bundle generate-cw-deployment-bundle

                    if gh release view "${TAG_NAME}" --repo CAAPIM/layer7-operator >/dev/null 2>&1; then
                        gh release upload "${TAG_NAME}" deploy/bundle.yaml deploy/cw-bundle.yaml --repo CAAPIM/layer7-operator --clobber
                    else
                        gh release create "${TAG_NAME}" deploy/bundle.yaml deploy/cw-bundle.yaml --repo CAAPIM/layer7-operator --title "${TAG_NAME}" --generate-notes
                    fi
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

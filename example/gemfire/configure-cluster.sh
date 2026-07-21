#!/bin/bash
#
# One-time gfsh bootstrap for the example GemFire cluster (cluster1.yaml): deploys the
# Layer7 GemFire function/security jars and creates the regions the Gateway's shared-state
# client expects. Safe to re-run: jar deploys are idempotent, and region-already-exists
# failures from `create region` are ignored.

set -euo pipefail

NAMESPACE=${NAMESPACE:-default}
CLUSTER_NAME=${CLUSTER_NAME:-cluster1}

if [ -z "${FUNCTIONS_JAR:-}" ] || [ ! -f "${FUNCTIONS_JAR}" ]; then
	echo "FUNCTIONS_JAR must be set to a valid path to the layer7-gemfire-functions jar"
	exit 1
fi

if [ -z "${SECURITY_JAR:-}" ] || [ ! -f "${SECURITY_JAR}" ]; then
	echo "SECURITY_JAR must be set to a valid path to the layer7-gemfire-security jar"
	exit 1
fi

LOCATOR_HOST="${CLUSTER_NAME}-locator-clusterip"
LOCATOR_PORT=10334

LOCATOR_POD=$(kubectl -n "${NAMESPACE}" get pod -l gemfire.vmware.com/app="${CLUSTER_NAME}"-locator -o jsonpath='{.items[0].metadata.name}')

if [ -z "${LOCATOR_POD}" ]; then
	echo "could not find a locator pod for cluster ${CLUSTER_NAME} in namespace ${NAMESPACE}"
	exit 1
fi

FUNCTIONS_JAR_NAME=$(basename "${FUNCTIONS_JAR}")
SECURITY_JAR_NAME=$(basename "${SECURITY_JAR}")

echo "copying jars to locator pod ${LOCATOR_POD}"
kubectl -n "${NAMESPACE}" cp "${FUNCTIONS_JAR}" "${LOCATOR_POD}:/tmp/${FUNCTIONS_JAR_NAME}"
kubectl -n "${NAMESPACE}" cp "${SECURITY_JAR}" "${LOCATOR_POD}:/tmp/${SECURITY_JAR_NAME}"

echo "running gfsh bootstrap against ${LOCATOR_HOST}[${LOCATOR_PORT}]"
kubectl -n "${NAMESPACE}" exec "${LOCATOR_POD}" -- gfsh \
	-e "connect --locator=${LOCATOR_HOST}[${LOCATOR_PORT}]" \
	-e "deploy --jar=/tmp/${FUNCTIONS_JAR_NAME}" \
	-e "deploy --jar=/tmp/${SECURITY_JAR_NAME}" \
	-e "create region --name=layer7gw_keyvalue --type=PARTITION --redundant-copies=1" \
	-e "create region --name=layer7gw_session --type=PARTITION --redundant-copies=1" \
	-e "create region --name=layer7gw_sortedset --type=PARTITION --redundant-copies=1" \
	-e "create region --name=layer7gw_ratelimiter --type=PARTITION --redundant-copies=1" \
	-e "create region --name=layer7gw_counter --type=PARTITION --redundant-copies=1" \
	|| echo "gfsh reported an error above — if it's only about regions already existing, this is safe to ignore"

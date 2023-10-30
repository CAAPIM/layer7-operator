#!/bin/bash
SERVICE_ACCOUNT_TOKEN=$(cat /var/run/secrets/kubernetes.io/serviceaccount/token)
UPDATE_SERVICE_ACCOUNT_TOKEN=$(sed -e "s~SERVICE_ACCOUNT_TOKEN~${SERVICE_ACCOUNT_TOKEN}~g" /opt/docker/rc.d/base/update-service-account-token.xml)
echo "${UPDATE_SERVICE_ACCOUNT_TOKEN}" > ${GATEWAY_DIR}/node/default/etc/bootstrap/bundle/update_service_account_token.xml.req.bundle
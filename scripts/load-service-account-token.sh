#!/bin/bash
# Substitute in the shell so the token never appears in a process command line.
SERVICE_ACCOUNT_TOKEN=$(< /var/run/secrets/kubernetes.io/serviceaccount/token)
TEMPLATE=$(< /opt/docker/rc.d/base/update-service-account-token.xml)
printf '%s%s%s\n' "${TEMPLATE%%SERVICE_ACCOUNT_TOKEN*}" "${SERVICE_ACCOUNT_TOKEN}" "${TEMPLATE#*SERVICE_ACCOUNT_TOKEN}" > ${GATEWAY_DIR}/node/default/etc/bootstrap/bundle/update_service_account_token.xml.req.bundle

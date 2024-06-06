#!/bin/bash
namespace=$1

if [[ -z ${namespace} ]]; then
echo "please set namespace"
exit 1
fi

# Update portal settings for Redis
CREATE_TS="$(date +%s)"
kubectl -n ${namespace} exec portal-mysql-0 -- mysql -s -uroot -p7layer dev_portal_portal -e "insert into SETTING(UUID, NAME, VALUE, CREATE_TS, MODIFY_TS, CREATED_BY, TENANT_ID) VALUES('4f33be8f-186a-12e6-8d56-000c295530e3', 'REDIS_GROUP_NAME', 'l7GW:KeyValueStore:apikeys', ${CREATE_TS}, '0', 'admin', 'portal');" 2> /dev/null
kubectl -n ${namespace} exec portal-mysql-0 -- mysql -s -uroot -p7layer dev_portal_portal -e "insert into SETTING(UUID, NAME, VALUE, CREATE_TS, MODIFY_TS, CREATED_BY, TENANT_ID) VALUES('4f33be8f-186a-12e6-7c56-000d295530e4', 'REDIS_KEY_STORE', 'standalone-redis-master:6379', ${CREATE_TS}, '0', 'admin', 'portal');" 2> /dev/null

echo "Portal Configured to use Redis"
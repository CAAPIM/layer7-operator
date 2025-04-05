#!/bin/bash
STATUSCODE=$(/bin/curl --silent -k --write-out "%{http_code}" https://localhost:$OTK_HEALTHCHECK_PORT/auth/oauth/health)

if test $STATUSCODE -ne 200; then
	exit 1
fi

exit 0
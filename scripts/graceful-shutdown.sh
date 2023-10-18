#!/bin/bash
PERIOD_SECONDS=$1
CHECK_INTERVAL=$2
EXCLUDE="${@:3}"

PORTS=$(netstat -an | awk '/LISTEN/ { n = split($4, a, ":"); print a[n] }')

if [ -z $PERIOD_SECONDS ]; then
PERIOD_SECONDS=30
fi

TIMEOUT_TS=$(($(date +%s) + $PERIOD_SECONDS))

# Remove excluded ports
for p in $EXCLUDE; do
  PORTS=("${PORTS[@]/$p}")
done

while [ $(date +%s) -lt $TIMEOUT_TS ]; do
    BUSY_PORTS=0
    for p in $PORTS; do
        # Check open connections
        CONNECTIONS=$(netstat -anp | grep ESTABLISHED | grep java | awk '{ print $4 }' | grep :$p | wc -l)
        if [ $CONNECTIONS -gt 0 ]; then
            let BUSY_PORTS++
            echo Port $p has $CONNECTIONS connection open
        fi
    done

    if [ $BUSY_PORTS -eq 0 ]; then
        echo "no open connections"
        exit 0
    fi
    sleep $CHECK_INTERVAL
done
exit 0
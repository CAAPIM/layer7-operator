#!/bin/bash
BASE_CONFIG_DIR="/opt/docker/custom"
GRAPHMAN_CONFIG_DIR="/opt/docker/graphman"
BUNDLE_DIR="$BASE_CONFIG_DIR/bundle"
CUSTOM_ASSERTIONS_DIR="$BASE_CONFIG_DIR/custom-assertions"
MODULAR_ASSERTIONS_DIR="$BASE_CONFIG_DIR/modular-assertions"
EXTERNAL_LIBRARIES_DIR="$BASE_CONFIG_DIR/external-libraries"
CUSTOM_PROPERTIES_DIR="$BASE_CONFIG_DIR/custom-properties"
CUSTOM_HEALTHCHECK_SCRIPTS_DIR="$BASE_CONFIG_DIR/health-checks"
CUSTOM_SHELL_SCRIPTS_DIR="$BASE_CONFIG_DIR/scripts"

BASE_TARGET_DIR="/opt/SecureSpan/Gateway"
GRAPHMAN_BOOTSTRAP_DIR="$BASE_TARGET_DIR/node/default/etc/bootstrap/bundle"
TARGET_CUSTOM_ASSERTIONS_DIR="$BASE_TARGET_DIR/runtime/modules/lib"
TARGET_MODULAR_ASSERTIONS_DIR="$BASE_TARGET_DIR/runtime/modules/assertions"
TARGET_EXTERNAL_LIBRARIES_DIR="$BASE_TARGET_DIR/runtime/lib/ext"
TARGET_BUNDLE_DIR="$BASE_TARGET_DIR/node/default/etc/bootstrap/bundle"
TARGET_CUSTOM_PROPERTIES_DIR="$BASE_TARGET_DIR/node/default/etc/conf"
TARGET_HEALTHCHECK_DIR="/opt/docker/rc.d/diagnostic/health_check"

error() {
    echo "ERROR - ${1}" 1>&2
    exit 1
}
function cleanup() {
    echo "***************************************************************************"
    echo "removing $BASE_CONFIG_DIR"
    echo "***************************************************************************"
    rm -rf $BASE_CONFIG_DIR/*
}

function copy() {
    TYPE=$1
    EXT=$2
    SOURCE_DIR=$3
    TARGET_DIR=$4
    echo "***************************************************************************"
    echo "scanning for $TYPE in $SOURCE_DIR"
    echo "***************************************************************************"
    FILES=$(find $3 -type f -name '*'$2 2>/dev/null)
    for file in $FILES; do
        name=$(basename "$file")
        cp $file $4/$name
        echo -e "$name written to $4/$name"
    done
}


function gunzip() {
    TYPE=$1
    EXT=$2
    SOURCE_DIR=$3
    echo "***************************************************************************"
    echo "scanning for $TYPE in $SOURCE_DIR"
    echo "***************************************************************************"
    FILES=$(find $3 -type f -name '*'$2 2>/dev/null)
    for file in $FILES; do
        fullname=$(basename "$file")
        name="${fullname%.*}"
        cat $file | gzip -d > $GRAPHMAN_BOOTSTRAP_DIR/$name".json"
        echo -e "$name decompressed"
    done
}

function run() {
    TYPE=$1
    EXT=$2
    SOURCE_DIR=$3
    echo "***************************************************************************"
    echo "scanning for $TYPE in $SOURCE_DIR"
    echo "***************************************************************************"
    FILES=$(find $3 -type f -name '*'$2 2>/dev/null)
    for file in $FILES; do
        name=$(basename "$file")
        echo -e "running $name"
        /bin/bash $file
        if [ $? -ne 0 ]; then
            echo "Failed executing the script: $file"
            exit 1
        fi
    done
}

gunzip "graphman bundles" ".gz" $GRAPHMAN_CONFIG_DIR
copy "bundles" ".bundle" $BUNDLE_DIR $TARGET_BUNDLE_DIR
copy "custom assertions" ".jar" $CUSTOM_ASSERTIONS_DIR $TARGET_CUSTOM_ASSERTIONS_DIR
copy "modular assertions" ".aar" $MODULAR_ASSERTIONS_DIR $TARGET_MODULAR_ASSERTIONS_DIR
copy "external libraries" ".jar" $EXTERNAL_LIBRARIES_DIR $TARGET_EXTERNAL_LIBRARIES_DIR
copy "custom properties" ".properties" $CUSTOM_PROPERTIES_DIR $TARGET_CUSTOM_PROPERTIES_DIR
copy "custom health checks" ".sh" $CUSTOM_HEALTHCHECK_SCRIPTS_DIR $TARGET_HEALTHCHECK_DIR
run "custom shell scripts" ".sh" $CUSTOM_SHELL_SCRIPTS_DIR


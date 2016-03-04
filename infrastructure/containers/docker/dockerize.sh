#!/usr/bin/env bash

set -ex
set -o pipefail

CONTAINER_NAME=${2}
PROJECT_DIR="${PWD}"

if [ "${CONTAINER_NAME}" == "corporate" ]; then
    # Build the static things
    declare -a STATIC_COMPONENTS=( "dashboard" "website" )
    for STATIC_COMPONENT in "${STATIC_COMPONENTS[@]}"
    do
        cd ${PROJECT_DIR}/${STATIC_COMPONENT}
        npm run clean
        npm run bundle
    done

    docker build -f ${PROJECT_DIR}/infrastructure/containers/docker/corporate.docker \
        -t ${CONTAINER_NAME}:${CIRCLE_BUILD_NUM} \
        "${PROJECT_DIR}"
    exit 0
fi

if [ "${CONTAINER_NAME}" == "intaker" ]; then
    BINARY_FILE=${1}
    CONFIG_FILE=${3}

    docker build -f ${PROJECT_DIR}/infrastructure/containers/docker/intaker.docker \
        -t ${CONTAINER_NAME}:${CIRCLE_BUILD_NUM} \
        --build-arg BINARY_FILE=${BINARY_FILE} \
        --build-arg CONFIG_FILE=${CONFIG_FILE} \
        "${PROJECT_DIR}"

    exit 0
fi

exit 1
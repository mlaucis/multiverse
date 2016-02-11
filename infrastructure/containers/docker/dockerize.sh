#!/usr/bin/env bash

set -ex
set -o pipefail

CONTAINER_NAME=${1}
BINARY_FILE=${2}
CONFIG_FILE=${3}

PROJECT_DIR="${PWD}"

docker build -f ${PROJECT_DIR}/infrastructure/containers/docker/intaker.docker \
    -t ${1}:${CIRCLE_BUILD_NUM} \
    --build-arg BINARY_FILE=${BINARY_FILE} \
    --build-arg CONFIG_FILE=${CONFIG_FILE} \
    "${PROJECT_DIR}"
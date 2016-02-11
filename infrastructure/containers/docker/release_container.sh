#!/usr/bin/env bash

set -ex
set -o pipefail

CONTAINER_NAME=${1}
BINARY_FILE=${2}
CONFIG_FILE=${3}

PROJECT_DIR="${PWD}"

${PROJECT_DIR}/infrastructure/containers/docker/dockerize.sh "${CONTAINER_NAME}" "${BINARY_FILE}" "${CONFIG_FILE}"

aws ecr get-login --region us-east-1 | source /dev/stdin
docker tag ${CONTAINER_NAME}:${CIRCLE_BUILD_NUM} 775034650473.dkr.ecr.us-east-1.amazonaws.com/${CONTAINER_NAME}:${CIRCLE_BUILD_NUM}

aws ecr get-login --region us-east-1 | source /dev/stdin
docker push 775034650473.dkr.ecr.us-east-1.amazonaws.com/${CONTAINER_NAME}:${CIRCLE_BUILD_NUM}

rm -f /home/ubuntu/.docker/config.json

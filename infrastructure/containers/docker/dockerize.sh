#!/usr/bin/env bash

set -ex
set -o pipefail

CONTAINER_NAME=${1}
PROJECT_DIR="${PWD}"

PROJECT_NAME='github.com/tapglue/multiverse'
VENDOR_DIR='Godeps/_workspace'

REVISION=`git rev-parse HEAD`

CONTAINER_GOPATH='/go'
CONTAINER_PROJECT_DIR="${CONTAINER_GOPATH}/src/${PROJECT_NAME}"
CONTAINER_PROJECT_GOPATH="${CONTAINER_PROJECT_DIR}/${VENDOR_DIR}:${CONTAINER_GOPATH}"

if [ "${CONTAINER_NAME}" == "dashboard" ]; then
    cd ${PROJECT_DIR}/dashboard
    npm run clean
    npm run bundle

    docker build -f ${PROJECT_DIR}/infrastructure/containers/docker/dashboard.docker \
        -t ${CONTAINER_NAME}:${CIRCLE_BUILD_NUM} \
        "${PROJECT_DIR}"
    exit 0
fi

if [ "${CONTAINER_NAME}" == "gateway-http" ]; then
    BINARY_FILE=${2}
    CONFIG_FILE=${3}

    docker run --rm \
        -v ${PROJECT_DIR}:${CONTAINER_PROJECT_DIR} \
        -e GOPATH=${CONTAINER_PROJECT_GOPATH} \
        -e GODEBUG=netdns=go \
        -w "${CONTAINER_PROJECT_DIR}" \
        golang:1.5.3-alpine \
        go build -v -ldflags "-X main.currentRevision=${REVISION}" -tags postgres -o ${BINARY_FILE} cmd/intaker/intaker.go

    docker build -f ${PROJECT_DIR}/infrastructure/containers/docker/gateway-http.docker \
        -t ${CONTAINER_NAME}:${CIRCLE_BUILD_NUM} \
        --build-arg BINARY_FILE=${BINARY_FILE} \
        --build-arg CONFIG_FILE=${CONFIG_FILE} \
        "${PROJECT_DIR}"

    exit 0
fi

if [ "${CONTAINER_NAME}" == "pganalyze" ]; then
  CONFIG_FILE=${3}

  docker build -f ${PROJECT_DIR}/infrastructure/containers/docker/pganalyze.docker \
    -t ${CONTAINER_NAME}:${CIRCLE_BUILD_NUM} \
    --build-arg CONFIG_FILE=${CONFIG_FILE} \
    "${PROJECT_DIR}"

  exit 0
fi

if [ "${CONTAINER_NAME}" == "reporter" ]; then
    BINARY_FILE=${2}

    docker run --rm \
        -v ${PROJECT_DIR}:${CONTAINER_PROJECT_DIR} \
        -e GOPATH=${CONTAINER_PROJECT_GOPATH} \
        -e GODEBUG=netdns=go \
        -w "${CONTAINER_PROJECT_DIR}" \
        golang:1.5.3-alpine \
        go build -v -ldflags "-X main.currentRevision=${REVISION}" -o ${BINARY_FILE} cmd/reporter/reporter.go

    docker build -f ${PROJECT_DIR}/infrastructure/containers/docker/reporter.docker \
        -t ${CONTAINER_NAME}:${CIRCLE_BUILD_NUM} \
        --build-arg BINARY_FILE=${BINARY_FILE} \
        "${PROJECT_DIR}"

    exit 0
fi

exit 1

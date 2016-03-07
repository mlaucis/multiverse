#!/usr/bin/env bash

set -e
set -o pipefail

PROJECT_NAME='github.com/tapglue/multiverse'
PROJECT_DIR="${PWD}/../../.."
VENDOR_DIR='Godeps/_workspace'

REVISION=`git rev-parse HEAD`

CONTAINER_GOPATH='/gopath'
CONTAINER_PROJECT_DIR="${CONTAINER_GOPATH}/src/${PROJECT_NAME}"
CONTAINER_PROJECT_GOPATH="${CONTAINER_PROJECT_DIR}/${VENDOR_DIR}:${CONTAINER_GOPATH}"

docker run --rm \
    -v ${PROJECT_DIR}:${CONTAINER_PROJECT_DIR} \
    -e GOPATH=${CONTAINER_PROJECT_GOPATH} \
    -e GODEBUG=netdns=go \
    -w "${CONTAINER_PROJECT_DIR}" \
    golang:1.5.3-alpine \
    go build -v -ldflags "-X main.currentRevision=${REVISION}" -tags redis -o intaker_redis_${CIRCLE_BUILD_NUM} cmd/intaker/intaker.go

# If we want to optimize for space then we can strip the debugging information out of the binary
# WARNING: this might cause the stacktraces to become useless if the app panics
#strip "${PROJECT_DIR}/intaker_redis_${CIRCLE_BUILD_NUM}"

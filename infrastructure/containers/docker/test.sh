#!/usr/bin/env bash

set -ex
set -o pipefail

TEST_COMPONENT=${1}
TEST_TARGET=${2}

PROJECT_NAME='github.com/tapglue/multiverse'
PROJECT_DIR="${PWD}"
VENDOR_DIR='Godeps/_workspace'

REVISION=`git rev-parse HEAD`

CONTAINER_GOPATH='/go'
CONTAINER_PROJECT_DIR="${CONTAINER_GOPATH}/src/${PROJECT_NAME}"
CONTAINER_PROJECT_GOPATH="${CONTAINER_PROJECT_DIR}/${VENDOR_DIR}:${CONTAINER_GOPATH}"

if [ ${TEST_COMPONENT} == "object" ]
then
    rm -f output.log

    docker run --rm \
        --net="host" \
        -v ${PROJECT_DIR}:${CONTAINER_PROJECT_DIR} \
        -e CI=true \
        -e GODEBUG=netdns=go \
        -e GOPATH=${CONTAINER_PROJECT_GOPATH} \
        -w "${CONTAINER_PROJECT_DIR}" \
        golang:1.5.3 \
        go test -v -race ./controller 2> output.log

    # Check for race conditions as we don't have a proper exit code for them from the tool
    cat output.log | grep -v 'WARNING: DATA RACE'

    rm -f output.log

    docker run --rm \
        --net="host" \
        -v ${PROJECT_DIR}:${CONTAINER_PROJECT_DIR} \
        -e CI=true \
        -e GODEBUG=netdns=go \
        -e GOPATH=${CONTAINER_PROJECT_GOPATH} \
        -w "${CONTAINER_PROJECT_DIR}" \
        golang:1.5.3 \
        go test -v -race -tags integration ./service/object -postgres.url="postgres://ubuntu:unicode@127.0.0.1/circle_test?sslmode-disable" 2> output.log

    # Check for race conditions as we don't have a proper exit code for them from the tool
    cat output.log | grep -v 'WARNING: DATA RACE'

    exit $?
fi

if [ ${TEST_COMPONENT} == "redis" ]
then
    rm -f output.log

    docker run --rm \
        --net="host" \
        -v ${PROJECT_DIR}:${CONTAINER_PROJECT_DIR} \
        -e CI=true \
        -e GODEBUG=netdns=go \
        -e GOPATH=${CONTAINER_PROJECT_GOPATH} \
        -w "${CONTAINER_PROJECT_DIR}/limiter/redis" \
        golang:1.5.3 \
        go test -v -race 2> output.log

    # Check for race conditions as we don't have a proper exit code for them from the tool
    cat output.log | grep -v 'WARNING: DATA RACE'

    exit $?
fi

declare -a VERSIONS=( "v03" "v04" )

if [ ${CIRCLE_BRANCH} == "master" ]
then
    declare -A TEST_MATRIX=( \
        ["intaker_postgres_v03"]=true \
        ["intaker_postgres_v04"]=true \
        ["intaker_redis_v03"]=true \
        ["intaker_redis_v04"]=true \
    )

    declare -A BUILD_MATRIX=( \
        ["intaker_postgres"]=true \
        ["intaker_redis"]=true \
    )
else
    declare -A TEST_MATRIX=( \
        ["intaker_postgres_v03"]=true \
        ["intaker_postgres_v04"]=true \
        ["intaker_redis_v03"]=false \
        ["intaker_redis_v04"]=false \
    )

    declare -A BUILD_MATRIX=( \
        ["intaker_postgres"]=true \
        ["intaker_redis"]=true \
    )
fi

CURRENT_BUILD_KEY="${TEST_COMPONENT}_${TEST_TARGET}"
if [ "${BUILD_MATRIX[${CURRENT_BUILD_KEY}]}" == true ]
then
    docker run --rm \
        -v ${PROJECT_DIR}:${CONTAINER_PROJECT_DIR} \
        -e GOPATH=${CONTAINER_PROJECT_GOPATH} \
        -e GODEBUG=netdns=go \
        -w "${CONTAINER_PROJECT_DIR}" \
        golang:1.5.3-alpine \
        go build -v -ldflags "-X main.currentRevision=${REVISION}" -tags ${TEST_TARGET} -o docker_${TEST_COMPONENT}_${TEST_TARGET}_${CIRCLE_BUILD_NUM} cmd/${TEST_COMPONENT}/${TEST_COMPONENT}.go

    # If we want to optimize for size of the binary  then we can strip the debugging information out to get around 5-6mb less
    # WARNING: this might cause the stacktraces to become useless if the app panics
    #strip "${PROJECT_DIR}/intaker_redis_${CIRCLE_BUILD_NUM}"
fi

for VERSION in "${VERSIONS[@]}"
do
    CURRENT_TEST_KEY="${TEST_COMPONENT}_${TEST_TARGET}_${VERSION}"
    if [ "${TEST_MATRIX[${CURRENT_TEST_KEY}]}" == false ]
    then
        continue
    fi

    rm -f output.log

    echo "Testing github.com/tapglue/multiverse/${VERSION}/server"

    docker run --rm \
        --net="host" \
        -v ${PROJECT_DIR}:${CONTAINER_PROJECT_DIR} \
        -e CI=true \
        -e GODEBUG=netdns=go \
        -e GOPATH=${CONTAINER_PROJECT_GOPATH} \
        -w "${CONTAINER_PROJECT_DIR}/${VERSION}/server" \
        golang:1.5.3 \
        go test -v -check.v -race -tags ${TEST_TARGET} github.com/tapglue/multiverse/${VERSION}/server 2> output.log

    # Check for race conditions as we don't have a proper exit code for them from the tool
    cat output.log | grep -v 'WARNING: DATA RACE'
done

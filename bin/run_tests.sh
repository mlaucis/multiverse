#!/usr/bin/env bash

set -ex
set -o pipefail

TEST_COMPONENT=${1}
TEST_TARGET=${2}

export PATH=/home/ubuntu/.gimme/versions/go1.5.2.linux.amd64/bin:${PATH}
export GOPATH=`godep path`:${GOPATH}
REVISION=`git rev-parse HEAD`
cd /home/ubuntu/.go_workspace/src/github.com/tapglue/multiverse
CWD=`pwd`

if [ ${TEST_COMPONENT} == "controller" ]
then
  go test \
    -v \
    -race \
    ./controller/...

  exit 0
fi

if [ ${TEST_COMPONENT} == "service" ]
then
  go test \
    -v \
    -race \
    -tags integration \
    ./service/... \
    -postgres.url="postgres://ubuntu:unicode@127.0.0.1/circle_test?sslmode-disable"

  exit 0
fi

if [ ${TEST_COMPONENT} == "redis" ]
then
    cd ${CWD}/limiter/redis
    go test
    exit $?
fi

declare -a VERSIONS=( "v03" "v04" )

if [ ${CIRCLE_BRANCH} == "master" ]
then
    declare -A TEST_MATRIX=( \
        ["intaker_postgres_v03"]=true \
        ["intaker_postgres_v04"]=false \
        ["intaker_redis_v03"]=false \
        ["intaker_redis_v04"]=false \
    )

    declare -A BUILD_MATRIX=( \
        ["intaker_postgres"]=false \
        ["intaker_redis"]=false \
    )
else
    declare -A TEST_MATRIX=( \
        ["intaker_postgres_v03"]=true \
        ["intaker_postgres_v04"]=false \
        ["intaker_redis_v03"]=false \
        ["intaker_redis_v04"]=false \
    )

    declare -A BUILD_MATRIX=( \
        ["intaker_postgres"]=false \
        ["intaker_redis"]=false \
    )
fi

CURRENT_BUILD_KEY="${TEST_COMPONENT}_${TEST_TARGET}"
if [ "${BUILD_MATRIX[${CURRENT_BUILD_KEY}]}" == true ]
then
    go build \
        -ldflags "-X main.currentRevision=${REVISION}" \
        -tags ${TEST_TARGET} \
        -o ${TEST_COMPONENT}_${TEST_TARGET}_${CIRCLE_BUILD_NUM} \
        cmd/${TEST_COMPONENT}/${TEST_COMPONENT}.go
fi

for VERSION in "${VERSIONS[@]}"
do
    CURRENT_TEST_KEY="${TEST_COMPONENT}_${TEST_TARGET}_${VERSION}"
    if [ "${TEST_MATRIX[${CURRENT_TEST_KEY}]}" == false ]
    then
        continue
    fi

    cd ${CWD}/${VERSION}/server

    rm -f c.out output.log

    echo "Testing github.com/tapglue/multiverse/${VERSION}/server"
    go test \
        -v \
        -race \
        -coverprofile=c.out \
        -tags ${TEST_TARGET} \
        -check.v \
        -coverpkg=github.com/tapglue/multiverse/${VERSION}/core/${TEST_TARGET},github.com/tapglue/multiverse/${VERSION}/server/handlers/${TEST_TARGET},github.com/tapglue/multiverse/${VERSION}/storage/${TEST_TARGET},github.com/tapglue/multiverse/${VERSION}/validator,github.com/tapglue/multiverse/${VERSION}/server/response,github.com/tapglue/multiverse/${VERSION}/errmsg,github.com/tapglue/multiverse/${VERSION}/storage/helper \
        github.com/tapglue/multiverse/${VERSION}/server 2> output.log

    cat output.log

    gocov convert c.out | gocov annotate - > coverage_server_${VERSION}_${TEST_TARGET}.json

    # Check for race conditions as we don't have a proper exit code for them from the tool
    cat output.log | grep -v 'WARNING: DATA RACE'
done

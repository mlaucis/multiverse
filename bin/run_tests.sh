#!/bin/bash

TEST_COMPONENT=${1}
TEST_TARGET=${2}

export GOPATH=`godep path`:${GOPATH}
REVISION=`git rev-parse HEAD`

go build -ldflags "-X main.currentRevision ${REVISION}" -tags ${TEST_TARGET} -o ${TEST_COMPONENT}_${TEST_TARGET}_${CIRCLE_BUILD_NUM} cmd/${TEST_COMPONENT}/${TEST_COMPONENT}.go

if [ "${TEST_COMPONENT}" == "distributor" ]
then
    # we don't have tests for distributor yet
    exit 0
fi

CWD=`pwd`

declare -a VERSIONS=( "v02" "v03" )
for VERSION in "${VERSIONS[@]}"
do
    cd ${CWD}/${VERSION}/server

    gocov test -race -tags ${TEST_TARGET} -check.v\
        -coverpkg=github.com/tapglue/backend/${VERSION}/core/${TEST_TARGET},\
github.com/tapglue/backend/${VERSION}/server/handlers/${TEST_TARGET},\
github.com/tapglue/backend/${VERSION}/storage/${TEST_TARGET},\
github.com/tapglue/backend/${VERSION}/validator,\
github.com/tapglue/backend/${VERSION}/server/response,\
github.com/tapglue/backend/${VERSION}/errmsg,\
github.com/tapglue/backend/${VERSION}/storage/helper \
github.com/tapglue/backend/${VERSION}/server > coverage_server_${VERSION}_${TEST_TARGET}.json
done

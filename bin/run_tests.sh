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

cd v02/server
gocov test -race -tags ${TEST_TARGET} -coverpkg=github.com/tapglue/backend/v02/core/${TEST_TARGET},github.com/tapglue/backend/v02/server/handlers/${TEST_TARGET},github.com/tapglue/backend/v02/storage/${TEST_TARGET},github.com/tapglue/backend/v02/validator,github.com/tapglue/backend/v02/server/response,github.com/tapglue/backend/v02/errmsg,github.com/tapglue/backend/v02/storage/helper -check.v github.com/tapglue/backend/v02/server > coverage_server_${TEST_TARGET}.json

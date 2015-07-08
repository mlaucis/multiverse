#!/bin/bash

TEST_TARGET=${1}
export GOPATH=`godep path`:${GOPATH}
REVISION=`git rev-parse HEAD`

go build -ldflags "-X main.currentRevision ${REVISION}" -tags ${TEST_TARGET} -o intaker_${TEST_TARGET}_${CIRCLE_BUILD_NUM} cmd/intaker/intaker.go
cd v02/server
gocov test -race -tags ${TEST_TARGET} -coverpkg=github.com/tapglue/backend/v02/core/${TEST_TARGET},github.com/tapglue/backend/v02/server/handlers/${TEST_TARGET},github.com/tapglue/backend/v02/storage/${TEST_TARGET},github.com/tapglue/backend/v02/validator -check.v github.com/tapglue/backend/v02/server > coverage_${TEST_TARGET}.json

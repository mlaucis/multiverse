#!/bin/bash

cd ${GOPATH}/src/github.com/tapglue/backend/v02/server/

declare -a TEST_TARGETS=("postgres" "kinesis")

export CI=true

for TEST_TARGET in "${TEST_TARGETS[@]}"
do
    go test -race -tags ${TEST_TARGET} -check.v github.com/tapglue/backend/v02/server
done

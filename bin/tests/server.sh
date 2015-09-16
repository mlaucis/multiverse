#!/usr/bin/env bash

declare -a TEST_TARGETS=("postgres" "kinesis")

export CI=true
CWD=`pwd`

for TEST_TARGET in "${TEST_TARGETS[@]}"
do
    declare -a VERSIONS=( "v02" "v03" )
    for VERSION in "${VERSIONS[@]}"
    do
        cd ${GOPATH}/src/github.com/tapglue/multiverse/${VERSION}/server/
        go test -race -tags ${TEST_TARGET} -check.v github.com/tapglue/multiverse/${VERSION}/server
    done
done

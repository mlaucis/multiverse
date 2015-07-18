#!/bin/bash

cd ${GOPATH}/src/github.com/tapglue/backend/v02/server/

declare -a TEST_TARGETS=("postgres" "kinesis")

export CI=true

for TEST_TARGET in "${TEST_TARGETS[@]}"
do
    gocov test -race -tags ${TEST_TARGET} -coverpkg=github.com/tapglue/backend/v02/core/${TEST_TARGET},github.com/tapglue/backend/v02/server/handlers/${TEST_TARGET},github.com/tapglue/backend/v02/storage/${TEST_TARGET},github.com/tapglue/backend/v02/validator,github.com/tapglue/backend/v02/server/response,github.com/tapglue/backend/v02/errmsg,github.com/tapglue/backend/v02/storage/helper -check.v github.com/tapglue/backend/v02/server > coverage_${TEST_TARGET}.json
    gocov-html coverage_${TEST_TARGET}.json > coverage_server_${TEST_TARGET}.html

    case "$(uname -s)" in
       Darwin)
         open -a Google\ Chrome coverage_server_${TEST_TARGET}.html &
         ;;

       *)
         google-chrome coverage_server_${TEST_TARGET}.html &
         ;;
    esac

    sleep 3
    rm coverage_server_${TEST_TARGET}.json coverage_server_${TEST_TARGET}.html
done

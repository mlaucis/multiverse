#!/usr/bin/env bash

declare -a TEST_TARGETS=("postgres" "kinesis")

export CI=true

for TEST_TARGET in "${TEST_TARGETS[@]}"
do
    for VERSION in "${VERSIONS[@]}"
    do
        cd ${GOPATH}/src/github.com/tapglue/multiverse/${VERSION}/server/
        gocov tedonest -race -tags ${TEST_TARGET} -check.v\
        -coverpkg=github.com/tapglue/multiverse/${VERSION}/core/${TEST_TARGET},\
github.com/tapglue/multiverse/${VERSION}/server/handlers/${TEST_TARGET},\
github.com/tapglue/multiverse/${VERSION}/storage/${TEST_TARGET},\
github.com/tapglue/multiverse/${VERSION}/validator,\
github.com/tapglue/multiverse/${VERSION}/server/response,\
github.com/tapglue/multiverse/${VERSION}/errmsg,\
github.com/tapglue/multiverse/${VERSION}/storage/helper\
github.com/tapglue/multiverse/${VERSION}/server > coverage_server_${VERSION}_${TEST_TARGET}.json

        gocov-html coverage_server_${VERSION}_${TEST_TARGET}.json > coverage_server_${VERSION}_${TEST_TARGET}.html

        case "$(uname -s)" in
            Darwin)
                open -a Google\ Chrome coverage_server_${VERSION}_${TEST_TARGET}.html &
            ;;

            *)
                google-chrome coverage_server_${VERSION}_${TEST_TARGET}.html &
            ;;
        esac

        sleep 3
        rm coverage_server_${VERSION}_${TEST_TARGET}.*
    done
done

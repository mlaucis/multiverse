#!/bin/bash

export GOPATH=`godep path`:${GOPATH}
declare -a TEST_TARGETS=("postgres" "kinesis", "redis")

CWD=`pwd`

for TEST_TARGET in "${TEST_TARGETS[@]}"
do

    declare -a VERSIONS=( "v02" "v03" )
    for VERSION in "${VERSIONS[@]}"
    do
        if [ ${TEST_TARGET} == "redis" ] && [ ${VERSION} == "v02" ]
        then
            continue
        fi

        cd ${CWD}/${VERSION}/server

        sed -i -e 's|'$WORKSPACE'/||g' coverage_server_${VERSION}_${TEST_TARGET}.json
        cat coverage_server_${VERSION}_${TEST_TARGET}.json | gocov-html > coverage_server_${VERSION}_${TEST_TARGET}.html
        gocov-xml < coverage_server_${VERSION}_${TEST_TARGET}.json > coverage_server_${VERSION}_${TEST_TARGET}.xml
        sed -i 's/\/home\/ubuntu\/backend\///g' coverage_server_${VERSION}_${TEST_TARGET}.html
        mv coverage_server_${VERSION}_${TEST_TARGET}.html ${CIRCLE_ARTIFACTS}/
    done
done

rm -f /home/ubuntu/.go_workspace/src/github.com/tapglue/backend

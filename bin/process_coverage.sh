#!/usr/bin/env bash

set -e
set -o pipefail

export PATH=/home/ubuntu/.gimme/versions/go1.5.2.linux.amd64/bin:${PATH}

rm -f /home/ubuntu/.go_workspace/src/github.com/tapglue/multiverse
exit 0

export GOPATH=`godep path`:${GOPATH}
declare -a TEST_TARGETS=("postgres" "redis")
TEST_COMPONENT="intaker"

if [ ${CIRCLE_BRANCH} == "master" ]
then
    declare -A BUILD_MATRIX=( \
        ["intaker_postgres_v03"]=true \
        ["intaker_redis_v03"]=true \
    )
else
    declare -A BUILD_MATRIX=( \
        ["intaker_postgres_v03"]=false \
        ["intaker_redis_v03"]=false \
    )
fi

CWD=`pwd`

for TEST_TARGET in "${TEST_TARGETS[@]}"
do
    declare -a VERSIONS=( "v03" )
    for VERSION in "${VERSIONS[@]}"
    do
        CURRENT_TEST_KEY="${TEST_COMPONENT}_${TEST_TARGET}_${VERSION}"
        if [ "${BUILD_MATRIX[${CURRENT_TEST_KEY}]}" == false ]
        then
            continue
        fi

        cd ${CWD}/${VERSION}/server

        sed -i -e 's|'$WORKSPACE'/||g' coverage_server_${VERSION}_${TEST_TARGET}.json
        cat coverage_server_${VERSION}_${TEST_TARGET}.json | gocov-html > coverage_server_${VERSION}_${TEST_TARGET}.html
        gocov-xml < coverage_server_${VERSION}_${TEST_TARGET}.json > coverage_server_${VERSION}_${TEST_TARGET}.xml
        sed -i 's/\/home\/ubuntu\/multiverse\///g' coverage_server_${VERSION}_${TEST_TARGET}.html
        mv coverage_server_${VERSION}_${TEST_TARGET}.html ${CIRCLE_ARTIFACTS}/
    done
done

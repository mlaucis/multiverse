#!/usr/bin/env bash

TEST_COMPONENT=${1}
TEST_TARGET=${2}

export PATH=/home/ubuntu/.gimme/versions/go1.5.1.linux.amd64/bin:${PATH}
export GOPATH=`godep path`:${GOPATH}
REVISION=`git rev-parse HEAD`

declare -a VERSIONS=( "v02" "v03" )

if [ ${CIRCLE_BRANCH} == "master" ]
then
    declare -A TEST_MATRIX=( \
        ["intaker_postgres_v02"]=false \
        ["intaker_postgres_v03"]=true \
        ["intaker_redis_v02"]=false \
        ["intaker_redis_v03"]=true \
        ["intaker_kinesis_v02"]=false \
        ["intaker_kinesis_v03"]=true \
        ["distributor_postgres_v02"]=false \
        ["distributor_postgres_v03"]=false \
    )

    declare -A BUILD_MATRIX=( \
        ["intaker_postgres"]=true \
        ["intaker_redis"]=true \
        ["intaker_kinesis"]=true \
        ["distributor_postgres"]=true \
    )
else
    declare -A TEST_MATRIX=( \
        ["intaker_postgres_v02"]=false \
        ["intaker_postgres_v03"]=true \
        ["intaker_redis_v02"]=false \
        ["intaker_redis_v03"]=false \
        ["intaker_kinesis_v02"]=false \
        ["intaker_kinesis_v03"]=true \
        ["distributor_postgres_v02"]=false \
        ["distributor_postgres_v03"]=false \
    )

    declare -A BUILD_MATRIX=( \
        ["intaker_postgres"]=true \
        ["intaker_redis"]=true \
        ["intaker_kinesis"]=true \
        ["distributor_postgres"]=true \
    )
fi

CWD=`pwd`

CURRENT_BUILD_KEY="${TEST_COMPONENT}_${TEST_TARGET}"
if [ "${BUILD_MATRIX[${CURRENT_BUILD_KEY}]}" == true ]
then
    if [ "${TEST_COMPONENT}" == "distributor" ]
    then
        go build \
            -ldflags "-X main.currentRevision=${REVISION}" \
            -tags ${TEST_TARGET} \
            -o ${TEST_COMPONENT}_${TEST_TARGET}_${CIRCLE_BUILD_NUM} \
            cmd/${TEST_COMPONENT}/${TEST_COMPONENT}.go cmd/${TEST_COMPONENT}/executor.go
    else
        go build \
            -ldflags "-X main.currentRevision=${REVISION}" \
            -tags ${TEST_TARGET} \
            -o ${TEST_COMPONENT}_${TEST_TARGET}_${CIRCLE_BUILD_NUM} \
            cmd/${TEST_COMPONENT}/${TEST_COMPONENT}.go
    fi
fi

for VERSION in "${VERSIONS[@]}"
do
    CURRENT_TEST_KEY="${TEST_COMPONENT}_${TEST_TARGET}_${VERSION}"
    if [ "${TEST_MATRIX[${CURRENT_TEST_KEY}]}" == false ]
    then
        continue
    fi

    cd ${CWD}/${VERSION}/server

    rm c.out output.log

    echo "Testing github.com/tapglue/multiverse/${VERSION}/server"
    go test \
        -race \
        -coverprofile=c.out \
        -tags ${TEST_TARGET} \
        -check.v \
        -coverpkg=github.com/tapglue/multiverse/${VERSION}/core/${TEST_TARGET},github.com/tapglue/multiverse/${VERSION}/server/handlers/${TEST_TARGET},github.com/tapglue/multiverse/${VERSION}/storage/${TEST_TARGET},github.com/tapglue/multiverse/${VERSION}/validator,github.com/tapglue/multiverse/${VERSION}/server/response,github.com/tapglue/multiverse/${VERSION}/errmsg,github.com/tapglue/multiverse/${VERSION}/storage/helper \
        github.com/tapglue/multiverse/${VERSION}/server 2> output.log

    # Check if the exit code was good or not
    if [ $? != 0 ]
    then
        cat output.log
        exit 1
    fi

    cat output.log

    gocov convert c.out | gocov annotate - > coverage_server_${VERSION}_${TEST_TARGET}.json

    # Check for race conditions as we don't have a proper exit code for them from the tool
    cat output.log | grep 'WARNING: DATA RACE' > /dev/null

    if [ $? != 1 ]
    then
        exit 1
    fi
done

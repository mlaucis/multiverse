#!/usr/bin/env bash

TEST_TARGET=${1}

CWD=`pwd`

cp cmd/${TEST_TARGET}/config.json_dist config.json

sed -i "s/APP_ENV/test/g" config.json
sed -i "s/APP_HOST_PORT/:8082/g" config.json
sed -i "s/AWS_ENDPOINT/http:\/\/127.0.0.1:4567/g" config.json
if [ "${CIRCLECI}" = true ] ; then
    sed -i "s/REDIS_HOST/127.0.0.1:6379/g" config.json
    sed -i "s/PG_DATABASE/circle_test/g" config.json
    sed -i "s/PG_USERNAME/ubuntu/g" config.json
    sed -i "s/PG_PASSWORD/unicode/g" config.json
    sed -i "s/PG_HOSTNAME/127.0.0.1/g" config.json
    sed -i "s/PG_OPTIONS/sslmode=disable/g" config.json
    sed -i "s/PG_SLAVES//g" config.json
else
    sed -i "s/REDIS_HOST/${WERCKER_REDIS_HOST}:${WERCKER_REDIS_PORT}/g" config.json
fi
sed -i "s/REDIS_DB_ID/0/g" config.json

declare -a VERSIONS=( "v02" "v03" )
for VERSION in "${VERSIONS[@]}"
do
    cp config.json ${CWD}/${VERSION}/server/
done

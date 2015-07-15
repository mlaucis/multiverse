#!/bin/bash

export GOPATH=`godep path`:${GOPATH}
declare -a targets=("postgres" "kinesis")

cd v02/server

for target in "${targets[@]}"
do
    sed -i -e 's|'$WORKSPACE'/||g' coverage_server_${target}.json
    cat coverage_server_${target}.json | gocov-html > coverage_server_${target}.html
    gocov-xml < coverage_server_${target}.json > coverage_server_${target}.xml
done

rm -f /home/ubuntu/.go_workspace/src/github.com/tapglue/backend

for target in "${targets[@]}"
do
    sed -i 's/\/home\/ubuntu\/backend\///g' coverage_server_${target}.html
    mv coverage_server_${target}.html ${CIRCLE_ARTIFACTS}/
done

#!/bin/bash

export GOPATH=`godep path`:${GOPATH}
declare -a targets=("postgres" "kinesis")

cd v02/server

for target in "${targets[@]}"
do
    cat coverage_${target}.json | gocov-xml > coverage_${target}.xml
    sed -i 's/\/home\/ubuntu\/backend\///g' coverage_${target}.xml
done

codecov --token=${CODECOV_TOKEN}

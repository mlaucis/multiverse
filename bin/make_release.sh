#!/usr/bin/env bash

set -e
set -o pipefail

componentName=${1}
releaseTarget=${2}

CWD=`pwd`

archiveName=${componentName}_${releaseTarget}.${CIRCLE_BUILD_NUM}.tar.gz
s3Location=s3://tapglue-builds/${componentName}/${releaseTarget}
releasesFilename='releases.json'

artifactName=${componentName}_${releaseTarget}_${CIRCLE_BUILD_NUM}
tarDir=${CWD}

if [ ${componentName} == "corporate" ]
then
    if [ ${releaseTarget} == "styleguide" ]
    then
        cp ${CWD}/infrastructure/nginx/corporate/styleguide ${CWD}/style/styleguide.nginx
        tarDir=${CWD}
        artifactName=style
    elif [ ${releaseTarget} == "dashboard" ] || [ ${releaseTarget} == "website" ]
    then
        cd ${CWD}/${releaseTarget}
        npm run clean
        npm run bundle
        cp ${CWD}/infrastructure/nginx/corporate/${releaseTarget} ${CWD}/${releaseTarget}/build/${releaseTarget}.nginx
        tarDir=${CWD}/${releaseTarget}
        artifactName=build
    fi
fi

cd ${CWD}
tar -C ${tarDir} -czf ${CWD}/${archiveName} ${artifactName}
cp ${archiveName} ${CIRCLE_ARTIFACTS}/
aws s3 cp ${archiveName} ${s3Location}/
aws s3 cp ${s3Location}/${releasesFilename} ${CIRCLE_ARTIFACTS}/${releasesFilename}
sed -i -r 's/latest_build": [0-9]+/latest_build": '${CIRCLE_BUILD_NUM}'/g' ${CIRCLE_ARTIFACTS}/${releasesFilename}
aws s3 cp ${CIRCLE_ARTIFACTS}/releases.json ${s3Location}/
rm -f ${CIRCLE_ARTIFACTS}/${releasesFilename}

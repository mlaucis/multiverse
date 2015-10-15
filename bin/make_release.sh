#!/usr/bin/env bash

finalComponentName=${1}
finalReleaseTarget=${2}

CWD=`pwd`

finalArchiveName=${finalComponentName}_${finalReleaseTarget}.${CIRCLE_BUILD_NUM}.tar.gz
finalS3Location=s3://tapglue-builds/${finalComponentName}/${finalReleaseTarget}
finalReleasesFilename='releases.json'

finalArtifactName=${finalComponentName}_${finalReleaseTarget}_${CIRCLE_BUILD_NUM}
tarDir=${CWD}

if [ ${finalComponentName} == "corporate" ]
then
    if [ ${finalReleaseTarget} == "styleguide" ]
    then
        cp ${CWD}/infrastructure/nginx/corporate/styleguide ${CWD}/style/styleguide.nginx
        tarDir=${CWD}
        finalArtifactName=style
    elif [ ${finalReleaseTarget} == "dashboard" ] || [ ${finalReleaseTarget} == "website" ]
    then
        cd ${CWD}/${finalReleaseTarget}
        npm run clean
        npm run bundle
        cp ${CWD}/infrastructure/nginx/corporate/${finalReleaseTarget} ${CWD}/${finalReleaseTarget}/build/${finalReleaseTarget}.nginx
        tarDir=${CWD}/${finalReleaseTarget}
        finalArtifactName=build
    fi
fi

cd ${CWD}
tar -C ${tarDir} -czf ${CWD}/${finalArchiveName} ${finalArtifactName}
cp ${finalArchiveName} ${CIRCLE_ARTIFACTS}/
aws s3 cp ${finalArchiveName} ${finalS3Location}/
aws s3 cp ${finalS3Location}/${finalReleasesFilename} ${CIRCLE_ARTIFACTS}/${finalReleasesFilename}
sed -i -r 's/latest_build": [0-9]+/latest_build": '${CIRCLE_BUILD_NUM}'/g' ${CIRCLE_ARTIFACTS}/${finalReleasesFilename}
aws s3 cp ${CIRCLE_ARTIFACTS}/releases.json ${finalS3Location}/
rm -f ${CIRCLE_ARTIFACTS}/${finalReleasesFilename}

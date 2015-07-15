#!/bin/bash

finalComponentName=${1}
finalReleaseTarget=${2}

finalExecName=${finalComponentName}_${finalReleaseTarget}_${CIRCLE_BUILD_NUM}
finalArchiveName=${finalComponentName}_${finalReleaseTarget}.${CIRCLE_BUILD_NUM}.tar.gz
finalS3Location=s3://tapglue-builds/${finalComponentName}/${finalReleaseTarget}
finalReleasesFilename=releases.json

tar -czf ${finalArchiveName} ${finalExecName}
cp ${finalArchiveName} ${CIRCLE_ARTIFACTS}/
aws s3 cp ${finalArchiveName} ${finalS3Location}/
aws s3 cp ${finalS3Location}/${finalReleasesFilename} ${CIRCLE_ARTIFACTS}/${finalReleasesFilename}
sed -i -r 's/latest_build": [0-9]+/latest_build": '${CIRCLE_BUILD_NUM}'/g' ${CIRCLE_ARTIFACTS}/${finalReleasesFilename}
aws s3 cp ${CIRCLE_ARTIFACTS}/releases.json ${finalS3Location}/
rm -f ${CIRCLE_ARTIFACTS}/${finalReleasesFilename}

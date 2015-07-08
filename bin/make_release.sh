#!/bin/bash

finalReleaseTarget=${1}
finalExecName=intaker_${finalReleaseTarget}_$CIRCLE_BUILD_NUM
finalArchiveName=intaker_${finalReleaseTarget}.$CIRCLE_BUILD_NUM.tar.gz
finalS3Location=s3://tapglue-builds/intaker/${finalReleaseTarget}
finalReleasesFilename=releases.json

tar -czf ${finalArchiveName} ${finalExecName}
cp ${finalArchiveName} ${CIRCLE_ARTIFACTS}/
aws s3 cp ${finalArchiveName} ${finalS3Location}/
aws s3 cp ${finalS3Location}/${finalReleasesFilename} ${CIRCLE_ARTIFACTS}/${finalReleasesFilename}
sed -i -r 's/latest_build": [0-9]+/latest_build": '${CIRCLE_BUILD_NUM}'/g' ${CIRCLE_ARTIFACTS}/${finalReleasesFilename}
aws s3 cp ${CIRCLE_ARTIFACTS}/releases.json ${finalS3Location}/
rm -f ${CIRCLE_ARTIFACTS}/${finalReleasesFilename}

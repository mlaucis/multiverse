#!/usr/bin/env bash

exec 1> >(logger -t $(basename ${0})) 2>&1

INSTANCE_ID=`wget -qO- http://instance-data/latest/meta-data/instance-id`
REGION=`wget -qO- http://instance-data/latest/meta-data/placement/availability-zone | sed 's/.$//'`
INSTANCE_TAGS="{"`aws ec2 describe-tags --region $REGION --filter "Name=resource-id,Values=$INSTANCE_ID" --output=text | sed -r 's/TAGS\t(.*)\t.*\t.*\t(.*)/"\1": "\2",/g'`"\"blank_tag\":\"blank_tag\"}"
INTAKER_DEPLOY_TARGET=`echo ${INSTANCE_TAGS} | python -c 'import sys,json;data=json.loads(sys.stdin.read()); print data["intaker_target"]'`
INSTALLER_CHANNEL=`echo ${INSTANCE_TAGS} | python -c 'import sys,json;data=json.loads(sys.stdin.read()); print data["installer_channel"]'`

logger -t intaker_installer got INSTANCE_ID: ${INSTANCE_ID}
logger -t intaker_installer got REGION: ${REGION}
logger -t intaker_installer got INSTANCE_TAGS: ${INSTANCE_TAGS}
logger -t tapglue_installer got INSTALLER_CAHNNEL: ${INSTALLER_CHANNEL}
logger -t intaker_installer got INTAKER_DEPLOY_TARGET: ${INTAKER_DEPLOY_TARGET}

if [ -z "${INTAKER_DEPLOY_TARGET}" ]; then
    logger -t intaker_installer installer target not found
    exit 1
fi

mkdir -p ~/releases/intaker/${INTAKER_DEPLOY_TARGET}
cd ~/releases/intaker/${INTAKER_DEPLOY_TARGET}

aws s3 cp s3://tapglue-builds/intaker/${INTAKER_DEPLOY_TARGET}/releases.json ./

releaseVersion=`cat ./releases.json | python -c 'import sys,json;data=json.loads(sys.stdin.read()); print data["current_release"]'`
execName=intaker_${INTAKER_DEPLOY_TARGET}_${releaseVersion}

aws s3 cp s3://tapglue-builds/intaker/${INTAKER_DEPLOY_TARGET}/intaker_${INTAKER_DEPLOY_TARGET}.${releaseVersion}.tar.gz ./
aws s3 cp s3://tapglue-builds/intaker/${INTAKER_DEPLOY_TARGET}/config.json ./

tar -zxvf intaker_${INTAKER_DEPLOY_TARGET}.${releaseVersion}.tar.gz

echo '#!/bin/bash

exec 1> >(logger -t $(basename ${0})) 2>&1

logger -t intaker runner received run command for intaker ${1} version ${2}

./intaker_${1}_${2}
' > run.sh

chmod +x ./run.sh

logger -t intaker_installer deployed intaker_${INTAKER_DEPLOY_TARGET}_${releaseVersion}

rm -f intaker_${INTAKER_DEPLOY_TARGET}.${releaseVersion}.tar.gz

kill -9 `ps aux | grep intaker_${INTAKER_DEPLOY_TARGET}_ | awk -F" " '{print $2}'`

nohup ./run.sh ${INTAKER_DEPLOY_TARGET} ${releaseVersion} &

#!/bin/bash

exec 1> >(logger -t $(basename ${0})) 2>&1

INSTANCE_ID=`wget -qO- http://instance-data/latest/meta-data/instance-id`
REGION=`wget -qO- http://instance-data/latest/meta-data/placement/availability-zone | sed 's/.$//'`
INSTANCE_TAGS="{"`aws ec2 describe-tags --region $REGION --filter "Name=resource-id,Values=$INSTANCE_ID" --output=text | sed -r 's/TAGS\t(.*)\t.*\t.*\t(.*)/"\1": "\2",/g'`"\"blank_tag\":\"blank_tag\"}"
INSTALLER_TARGET=`echo ${INSTANCE_TAGS} | python -c 'import sys,json;data=json.loads(sys.stdin.read()); print data["tapglue_installer"]'`

logger -t tapglue_installer got INSTANCE_ID: ${INSTANCE_ID}
logger -t tapglue_installer got REGION: ${REGION}
logger -t tapglue_installer got INSTANCE_TAGS: ${INSTANCE_TAGS}
logger -t tapglue_installer got INSTALLER_TARGET: ${INSTALLER_TARGET}

if [ -z "${INSTALLER_TARGET}" ]; then
    logger -t tapglue_installer installer target not found
    exit 1
fi

INSTALLER_SCRIPT="install_${INSTALLER_TARGET}.sh"

mkdir -p ~/releases/${INSTALLER_TARGET}/
cd ~/releases/${INSTALLER_TARGET}/

aws s3 cp s3://tapglue-builds/installer/${INSTALLER_SCRIPT} ./
chmod +x ./${INSTALLER_SCRIPT}

nohup ./${INSTALLER_SCRIPT} &

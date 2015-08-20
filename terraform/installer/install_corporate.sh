#!/usr/bin/env bash

exec 1> >(logger -t $(basename ${0})) 2>&1

INSTANCE_ID=`wget -qO- http://instance-data/latest/meta-data/instance-id`
REGION=`wget -qO- http://instance-data/latest/meta-data/placement/availability-zone | sed 's/.$//'`
INSTANCE_TAGS="{"`aws ec2 describe-tags --region $REGION --filter "Name=resource-id,Values=$INSTANCE_ID" --output=text | sed -r 's/TAGS\t(.*)\t.*\t.*\t(.*)/"\1": "\2",/g'`"\"blank_tag\":\"blank_tag\"}"
INTAKER_DEPLOY_TARGET=`echo ${INSTANCE_TAGS} | python -c 'import sys,json;data=json.loads(sys.stdin.read()); print data["intaker_target"]'`
INSTALLER_CHANNEL=`echo ${INSTANCE_TAGS} | python -c 'import sys,json;data=json.loads(sys.stdin.read()); print data["installer_channel"]'`

logger -t corporate_installer got INSTANCE_ID: ${INSTANCE_ID}
logger -t corporate_installer got REGION: ${REGION}
logger -t corporate_installer got INSTANCE_TAGS: ${INSTANCE_TAGS}
logger -t corporate_installer got INSTALLER_CAHNNEL: ${INSTALLER_CHANNEL}
logger -t corporate_installer got INTAKER_DEPLOY_TARGET: ${INTAKER_DEPLOY_TARGET}


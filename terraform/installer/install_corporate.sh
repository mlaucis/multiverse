#!/usr/bin/env bash

exec 1> >(logger -t $(basename ${0})) 2>&1

# We want to stop nginx as early as possible so that the healthcheck doesn't kick in and mark the instance as available
sudo service nginx stop

INSTANCE_ID=`wget -qO- http://instance-data/latest/meta-data/instance-id`
REGION=`wget -qO- http://instance-data/latest/meta-data/placement/availability-zone | sed 's/.$//'`
INSTANCE_TAGS="{"`aws ec2 describe-tags --region ${REGION} --filter "Name=resource-id,Values=$INSTANCE_ID" --output=text | sed -r 's/TAGS\t(.*)\t.*\t.*\t(.*)/"\1": "\2",/g'`"\"blank_tag\":\"blank_tag\"}"
INSTALLER_CHANNEL=`echo ${INSTANCE_TAGS} | python -c 'import sys,json;data=json.loads(sys.stdin.read()); print data["installer_channel"]'`

INSTALLER_COMPONENT='corporate'

logger -t ${INSTALLER_COMPONENT}_installer got INSTANCE_ID: ${INSTANCE_ID}
logger -t ${INSTALLER_COMPONENT}_installer got REGION: ${REGION}
logger -t ${INSTALLER_COMPONENT}_installer got INSTANCE_TAGS: ${INSTANCE_TAGS}
logger -t ${INSTALLER_COMPONENT}_installer got INSTALLER_CAHNNEL: ${INSTALLER_CHANNEL}

mkdir -p ~/releases/${INSTALLER_COMPONENT}/

declare -a INSTALLER_TARGETS=( "styleguide" )
for INSTALLER_TARGET in "${INSTALLER_TARGETS[@]}"
do
    cd ~/releases/${INSTALLER_COMPONENT}
    mkdir -p ${INSTALLER_TARGET}
    cd ~/releases/${INSTALLER_COMPONENT}/${INSTALLER_TARGET}

    # Wipe existing installation, if any (how? why?)
    rm -rf style

    aws s3 cp s3://tapglue-builds/${INSTALLER_COMPONENT}/${INSTALLER_TARGET}/releases.json ./
    releaseVersion=`cat ./releases.json | python -c 'import sys,json;data=json.loads(sys.stdin.read()); print data["current_release"]'`

    aws s3 cp s3://tapglue-builds/${INSTALLER_COMPONENT}/${INSTALLER_TARGET}/${INSTALLER_COMPONENT}_${INSTALLER_TARGET}.${releaseVersion}.tar.gz ./

    tar -zxvf ${INSTALLER_COMPONENT}_${INSTALLER_TARGET}.${releaseVersion}.tar.gz

    mv ./style/styleguide.nginx /etc/nginx/sites-available/styleguide

    ln -nfs /etc/nginx/sites-available/styleguide /etc/nginx/sites-enabled/styleguide

    logger -t ${INSTALLER_COMPONENT}_installer deployed ${INSTALLER_COMPONENT}_${INSTALLER_TARGET}_${releaseVersion}

    rm -f ${INSTALLER_COMPONENT}_${INSTALLER_TARGET}.${releaseVersion}.tar.gz

done

# Once everything is done we can start nginx again
sudo service nginx start

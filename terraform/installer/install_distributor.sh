#!/usr/bin/env bash

exec 1> >(logger -t $(basename ${0})) 2>&1

INSTANCE_ID=`wget -qO- http://instance-data/latest/meta-data/instance-id`
REGION=`wget -qO- http://instance-data/latest/meta-data/placement/availability-zone | sed 's/.$//'`
INSTANCE_TAGS="{"`aws ec2 describe-tags --region ${REGION} --filter "Name=resource-id,Values=$INSTANCE_ID" --output=text | sed -r 's/TAGS\t(.*)\t.*\t.*\t(.*)/"\1": "\2",/g'`"\"blank_tag\":\"blank_tag\"}"
DISTRIBUTOR_DEPLOY_TARGET=`echo ${INSTANCE_TAGS} | python -c 'import sys,json;data=json.loads(sys.stdin.read()); print data["distributor_target"]'`
INSTALLER_CHANNEL=`echo ${INSTANCE_TAGS} | python -c 'import sys,json;data=json.loads(sys.stdin.read()); print data["installer_channel"]'`

logger -t distributor_installer got INSTANCE_ID: ${INSTANCE_ID}
logger -t distributor_installer got REGION: ${REGION}
logger -t distributor_installer got INSTANCE_TAGS: ${INSTANCE_TAGS}
logger -t distributor_installer got INSTALLER_CAHNNEL: ${INSTALLER_CHANNEL}
logger -t distributor_installer got DISTRIBUTOR_DEPLOY_TARGET: ${DISTRIBUTOR_DEPLOY_TARGET}

if [ -z "${DISTRIBUTOR_DEPLOY_TARGET}" ]; then
    logger -t distributor_installer installer target not found
    exit 1
fi

mkdir -p ~/releases/distributor/${DISTRIBUTOR_DEPLOY_TARGET}
cd ~/releases/distributor/${DISTRIBUTOR_DEPLOY_TARGET}

aws s3 cp s3://tapglue-builds/distributor/${DISTRIBUTOR_DEPLOY_TARGET}/releases.json ./

releaseVersion=`cat ./releases.json | python -c 'import sys,json;data=json.loads(sys.stdin.read()); print data["current_release"]'`
execName=distributor_${DISTRIBUTOR_DEPLOY_TARGET}_${releaseVersion}

aws s3 cp s3://tapglue-builds/distributor/${DISTRIBUTOR_DEPLOY_TARGET}/distributor_${DISTRIBUTOR_DEPLOY_TARGET}.${releaseVersion}.tar.gz ./
aws s3 cp s3://tapglue-builds/distributor/${DISTRIBUTOR_DEPLOY_TARGET}/config.json ./

tar -zxvf distributor_${DISTRIBUTOR_DEPLOY_TARGET}.${releaseVersion}.tar.gz

echo '#!/usr/bin/env bash

exec 1> >(logger -t $(basename ${0})) 2>&1

logger -t distributor runner received run command for distributor ${1} version ${2}

./distributor_${1}_${2} -target '${DISTRIBUTOR_DEPLOY_TARGET}'
' > run.sh

chmod +x ./run.sh

logger -t distributor_installer deployed distributor_${DISTRIBUTOR_DEPLOY_TARGET}_${releaseVersion}

rm -f distributor_${DISTRIBUTOR_DEPLOY_TARGET}.${releaseVersion}.tar.gz

kill -9 `ps aux | grep distributor_${DISTRIBUTOR_DEPLOY_TARGET}_ | awk -F" " '{print $2}'`

nohup ./run.sh ${DISTRIBUTOR_DEPLOY_TARGET} ${releaseVersion} &

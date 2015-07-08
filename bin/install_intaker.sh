#!/bin/bash

# run crontab -e and add the following line
# @reboot $SCRIPT/PATH/install_intaker.sh $TARGET 2>&1 | /usr/bin/logger -t install_intaker
#
# $TARGET can be: kinesis or postgres
exec 1> >(logger -t $(basename ${0})) 2>&1

INSTANCE_ID=`wget -qO- http://instance-data/latest/meta-data/instance-id`
REGION=`wget -qO- http://instance-data/latest/meta-data/placement/availability-zone | sed 's/.$//'`
INSTANCE_TAGS=`aws ec2 describe-tags --region $REGION --filter "Name=resource-id,Values=$INSTANCE_ID" --output=text | sed -r 's/TAGS\t(.*)\t.*\t.*\t(.*)/\1="\2"/'`
INTAKER_DEPLOY_TARGET=`echo ${INSTANCE_TAGS} | grep 'intaker_target' | sed -r 's/(.*)="(.*)"/\2/'`

if [ -z "${INTAKER_DEPLOY_TARGET}" ]; then
    logger -t intaker_installer installer target not found
    exit 1
fi

mkdir -p ~/releases/intaker/${INTAKER_DEPLOY_TARGET}
cd ~/releases/intaker/${INTAKER_DEPLOY_TARGET}

aws s3 cp s3://tapglue-builds/intaker/${INTAKER_DEPLOY_TARGET}/releases.json ./

releaseVersion=`cat ./releases.json | python -mjson.tool | grep -i current | cut -d' ' -f 6 | sed 's/,//g'`
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

kill -9 `pgrep intaker_${INTAKER_DEPLOY_TARGET}_`

nohup `./run.sh ${INTAKER_DEPLOY_TARGET} ${releaseVersion}` &

logger -t intaker_installer started intaker_${INTAKER_DEPLOY_TARGET}_${releaseVersion}

rm -f intaker_${INTAKER_DEPLOY_TARGET}.${releaseVersion}.tar.gz

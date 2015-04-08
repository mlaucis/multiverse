#!/bin/bash
cd $GOPATH/src/github.com/tapglue/backend
cp config.json_dist config.json

sed -i "s/APP_ENV/test/g" config.json
sed -i "s/APP_HOST_PORT/:8082/g" config.json
if [ "$CIRCLECI" = true ] ; then
    sed -i "s/REDIS_HOST/127.0.0.1:6379/g" config.json
else
    sed -i "s/REDIS_HOST/${WERCKER_REDIS_HOST}:${WERCKER_REDIS_PORT}/g" config.json
fi
sed -i "s/REDIS_DB_ID/0/g" config.json

cp config.json v01/server/
cp config.json v02/server/

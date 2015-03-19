#!/bin/bash
cd $GOPATH/src/github.com/tapglue/backend
cp config.json_dist config.json

sed -i "s/APP_ENV/test/g" config.json
sed -i "s/APP_HOST_PORT/:8082/g" config.json
sed -i "s/NEWRELIC_KEY/${NEW_RELIC_LICENSE_KEY}/g" config.json
sed -i "s/NEWRELIC_NAME/test - tapglue/g" config.json
sed -i "s/REDIS_HOST/${WERCKER_REDIS_HOST}:${WERCKER_REDIS_PORT}/g" config.json
sed -i "s/REDIS_DB_ID/0/g" config.json

cp config.json v1/core/
cp config.json v1/server/

#!/bin/bash
cd $GOPATH/src/github.com/tapglue/backend
cp config.json_dist config.json

sed -i "s/APP_ENV/test/g" config.json
sed -i "s/APP_HOST_PORT/:8082/g" config.json
sed -i "s/NEWRELIC_KEY/${NEW_RELIC_LICENSE_KEY}/g" config.json
sed -i "s/NEWRELIC_NAME/test - tapglue/g" config.json
sed -i "s/DB_USER/${WERCKER_MYSQL_USERNAME}/g" config.json
sed -i "s/DB_PASS/${WERCKER_MYSQL_PASSWORD}/g" config.json
sed -i "s/DB_DB/${WERCKER_MYSQL_DATABASE}/g" config.json
sed -i "s/DB_MAX_IDLE/10/g" config.json
sed -i "s/DB_MAX_OPEN/300/g" config.json
sed -i "s/DB_MASTER_DEBUG/true/g" config.json
sed -i "s/DB_MASTER_HOST/${WERCKER_MYSQL_HOST}/g" config.json
sed -i "s/DB_MASTER_PORT/3306/g" config.json
sed -i "s/DB_SLAVE1_DEBUG/true/g" config.json
sed -i "s/DB_SLAVE1_HOST/${WERCKER_MYSQL_HOST}/g" config.json
sed -i "s/DB_SLAVE1_PORT/3306/g" config.json

cp config.json aerospike/
cp config.json db/
cp config.json server/

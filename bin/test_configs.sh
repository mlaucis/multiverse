#!/bin/bash
cd $GOPATH/src/github.com/tapglue/backend
cp config.json_dist config.json

sed -i "s/APP_ENV/test/g" config.json
sed -i "s/APP_HOST_PORT/:8082/g" config.json
sed -i "s/DB_USER/${WERCKER_MYSQL_USERNAME}/g" config.json
sed -i "s/DB_PASS/${WERCKER_MYSQL_PASSWORD}/g" config.json
sed -i "s/DB_DB/${WERCKER_MYSQL_DATABASE}/g" config.json
sed -i "s/DB_MAX_IDLE/10/g" config.json
sed -i "s/DB_MAX_OPEN/300/g" config.json
sed -i "s/DB_MASTER_DEBUG/true/g" config.json
sed -i "s/DB_MASTER_HOST/${WERCKER_MYSQL_HOST}/g" config.json
sed -i "s/DB_MASTER_PORT/${WERCKER_MYSQL_PORT}/g" config.json
sed -i "s/DB_SLAVE1_DEBUG/true/g" config.json
sed -i "s/DB_SLAVE1_HOST/${WERCKER_MYSQL_HOST}/g" config.json
sed -i "s/DB_SLAVE1_PORT/${WERCKER_MYSQL_PORT}/g" config.json

cp config.json db/
cp config.json server/

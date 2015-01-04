#!/bin/bash
export WERCKER_SOURCE_DIR=${GOPATH}/src/github.com/tapglue/backend
sed -i "s/CURRENT_TIMESTAMP/'2015-01-01 01:23:45'/g" ${WERCKER_SOURCE_DIR}/resources/sql/tapglue.sql
mysql -u ${WERCKER_MYSQL_USERNAME} -p${WERCKER_MYSQL_PASSWORD} -h ${WERCKER_MYSQL_HOST} -P ${WERCKER_MYSQL_PORT} ${WERCKER_MYSQL_DATABASE} < ${WERCKER_SOURCE_DIR}/resources/sql/tapglue.sql

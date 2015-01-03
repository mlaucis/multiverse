#!/bin/bash
mysql -u ${WERCKER_MYSQL_USERNAME} -p${WERCKER_MYSQL_PASSWORD} -h ${WERCKER_MYSQL_HOST} -P ${WERCKER_MYSQL_PORT} ${WERCKER_MYSQL_DATABASE} < ${GOPATH}/src/github.com/tapglue/backend/resources/sql/tapglue.sql

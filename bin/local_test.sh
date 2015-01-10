#!/bin/bash

DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )
mysql -u $1 -p$2 -h $3 -P $4 $5 -e 'SET foreign_key_checks = 0;DROP TABLE `accounts`, `account_users`, `applications`, `events`, `sessions`, `users`, `user_connections`;SET foreign_key_checks = 1;'
mysql -u $1 -p$2 -h $3 -P $4 $5 < ${DIR}/../resources/sql/tapglue.sql
go test -check.v

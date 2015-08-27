#!/usr/bin/env bash
mkdir -p ~/releases/

aws s3 cp s3://tapglue-builds/intaker/postgres/releases.json ~/releases/

releaseVersion=`cat ~/releases/releases.json | python -mjson.tool | grep -i current | cut -d' ' -f 6 | sed 's/,//g'`
execName=intaker_postgres_$releaseVersion

mkdir -p ~/releases/intaker/

aws s3 cp s3://tapglue-builds/intaker/postgres/intaker_postgres.$releaseVersion.tar.gz ~/releases/intaker/
aws s3 cp s3://tapglue-builds/intaker/postgres/config.json ~/releases/intaker/

cd ~/releases/intaker/
tar -zxvf intaker_postgres.$releaseVersion.tar.gz
nohup ./$execName >> ./log 2>&1 &

rm -f intaker_postgres.$releaseVersion.tar.gz

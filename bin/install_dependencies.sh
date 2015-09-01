#!/usr/bin/env bash

export PATH=/home/ubuntu/.gimme/versions/go1.5.linux.amd64/bin:${PATH}
CWD=`pwd`

echo "Installing Kinesalite"
docker run -d -t -p 4567:4567 dlsniper/kinesalite:1.8.0

echo "Installing gimme"
sudo curl -sL -o /usr/local/bin/gimme https://raw.githubusercontent.com/travis-ci/gimme/master/gimme
sudo chmod +x /usr/local/bin/gimme

echo "Installing go 1.5"
eval "$(GIMME_GO_VERSION=1.5 gimme)"

echo "Installing go dependencies"
go get github.com/tools/godep github.com/axw/gocov/gocov github.com/matm/gocov-html gopkg.in/check.v1

echo "Installing dashboard dependencies"
cd ${CWD}/dashboard
npm install

echo "Installing website dependencies"
cd ${CWD}/website
npm install


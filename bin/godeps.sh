#!/usr/bin/env bash

set -e

export PATH=/home/ubuntu/.gimme/versions/go1.5.2.linux.amd64/bin:${PATH}

go get -v -u github.com/tools/godep
cp -R $GOPATH/src/github.com/Tapglue $GOPATH/_src
godep restore
mv $GOPATH/_src $GOPATH/src/github.com/tapglue

#!/bin/bash
go get -v -u github.com/tools/godep
cp -R $GOPATH/src/github.com/Tapglue $GOPATH/_src
godep restore
mv $GOPATH/_src $GOPATH/src/github.com/tapglue

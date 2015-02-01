#!/bin/bash

cd ${GOPATH}/src/github.com/tapglue/backend/server
gocov test -race > coverage.json
gocov-html coverage.json > coverage.html
google-chrome coverage.html &
sleep 1
rm coverage.json coverage.html

#!/bin/bash

cd ${GOPATH}/src/github.com/tapglue/backend/core
gocov test -race > coverage.json
gocov-html coverage.json > coverage.html

case "$(uname -s)" in
   Darwin)
     open -a Google\ Chrome coverage.html &
     ;;

   *)
     google-chrome coverage.html &
     ;;
esac

sleep 3
rm coverage.json coverage.html

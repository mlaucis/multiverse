#!/bin/bash

cd ${GOPATH}/src/github.com/tapglue/backend/v01/server/
gocov test -race -coverpkg=github.com/tapglue/backend/v01/server,github.com/tapglue/backend/v01/core > coverage.json
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

cd ${GOPATH}/src/github.com/tapglue/backend/v02/server/
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

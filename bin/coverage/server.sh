#!/bin/bash

cd ${GOPATH}/src/github.com/tapglue/backend/v01/server/
gocov test -race -coverpkg=github.com/tapglue/backend/v01/server,github.com/tapglue/backend/v01/core,github.com/tapglue/backend/v01/context,github.com/tapglue/backend/v01/storage,github.com/tapglue/backend/v01/storage/redis,github.com/tapglue/backend/v01/validator,github.com/tapglue/backend/v01/validator/keys,github.com/tapglue/backend/v01/validator/tokens,github.com/tapglue/backend/v01/fixtures,github.com/tapglue/backend/v01/entity > coverage.json
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
gocov test -race -coverpkg=github.com/tapglue/backend/v02/server,github.com/tapglue/backend/v02/core,github.com/tapglue/backend/v02/context,github.com/tapglue/backend/v02/storage,github.com/tapglue/backend/v02/storage/redis,github.com/tapglue/backend/v02/validator,github.com/tapglue/backend/v02/validator/keys,github.com/tapglue/backend/v02/validator/tokens,github.com/tapglue/backend/v02/fixtures,github.com/tapglue/backend/v02/entity > coverage.json
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

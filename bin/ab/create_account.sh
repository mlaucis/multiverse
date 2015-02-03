#!/bin/bash

cd ${GOPATH}/src/github.com/tapglue/backend

echo "Building the backend app"
go build -o backend backend.go

echo "Launching the backend app"
./backend &
sleep 1

echo "Starting the siege"
ab -n 200000 -c 100 -p resources/ab_test/new_account_payload  -T 'application/json' localhost:8082/0.1/accounts &
sleep 3

echo "Starting pprof"
go tool pprof backend http://localhost:8082/debug/pprof/profile

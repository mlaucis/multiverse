#!/usr/bin/env bash

cd ${GOPATH}/src/github.com/tapglue/backend

echo "Building the backend app"
go build -o backend backend.go

echo "Launching the backend app"
2>/dev/null 1>&2 ./backend &
sleep 1

echo "Starting the siege"
ab -n 500000 -c 100 localhost:8082/0.1/account/100 &
sleep 1

echo "Starting pprof"
go tool pprof backend http://localhost:8082/debug/pprof/profile

kill `pidof backend`

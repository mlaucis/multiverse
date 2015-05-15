# Tapglue backend 

This repository contains the implementation of tapglues backend.

[![Circle CI](https://circleci.com/gh/tapglue/backend.svg?style=svg&circle-token=22a2b029440b825d23a4f0118274af084da917b2)](https://circleci.com/gh/tapglue/backend)
[![codecov.io](https://codecov.io/github/tapglue/backend/coverage.svg?token=OHlqgNOv66&branch=master)](https://codecov.io/github/tapglue/backend?branch=master)

## Documentation

See [Documentation](https://github.com/tapglue/backend/wiki) for entities, api design and more.

## System Requirements

- go (latest)
- redis 2.8 or newer

## Installation

Following steps are need to download and install this project.

### Getting started

Download the git repository to get started.

```shell
$ mkdir -p $GOPATH/src/github.com/tapglue/backend
$ git clone https://github.com/tapglue/backend.git
$ cd backend
```

### Dependencies

All dependecies should be fecthed correctly by running:

```shell
$ go get github.com/tapglue/backend
```

or, if you cloned it locally in your GOPATH

```shell
$ cd $GOPATH/src/github.com/tapglue/backend
$ go get ./...
```

### Server configuration

Configure the server including ports and database settings in the [config.json](config.json).

```json
{
  "env": "dev",
  "use_artwork": false,
  "listenHost": ":8082",
  "newrelic": {
    "key": "",
    "name": "dev - tapglue",
    "enabled": false
  },
  "redis": {
    "hosts": [
      "127.0.0.1:6379"
    ],
    "password": "",
    "db": 0,
    "pool_size": 30
  }
}
```

### Start server

```shell
$ go run -race backend.go
```

## Tests

```shell
$ cd core
$ go test -check.v
$ cd ../server
$ go test -check.v
```

## Coverage

```shell
$ bin/coverage/*.sh
```

## Benchmarks

```shell
$ cd core
$ go test -bench=. -benchmem
$ cd ../server
$ go test -bench=. -benchmem
```

## Profilling

```shell
$ bin/ab/*.sh
```

## Code commit

Before doing a commit, please run the following in the ```$GOPATH/src/github.com/tapglue/backend```  
```shell
goimports -w ./.. && golint ./... && go vet ./...
```

## Deployment

TBD


## Security test

Tool https://github.com/sqlmapproject/sqlmap

```bash
python sqlmap.py -u "http://127.0.0.1:8083/0.2/user/db9617bf-275e-521a-88c3-b6ef69d3af05*/events" \
-z "ignore-401,flu,bat" --banner -f \
--headers="Authorization: Basic ZTdhYWZjNDgxMWU4N2UyOTA3NjliNTdmOGNjYWI4NTA6U0RZcmJrUnVLR2w5ZUY5alZIazhNeXQ2Vm5jPQ=="
```

Bad keywords:

- google wfuzz: https://wfuzz.googlecode.com/svn/trunk/wordlist/Injections/SQL.txt

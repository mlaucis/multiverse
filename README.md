# Tapglue multiverse

This repository contains the implementation of tapglues multiverse.

[![Circle CI](https://circleci.com/gh/tapglue/multiverse.svg?style=svg&circle-token=22a2b029440b825d23a4f0118274af084da917b2)](https://circleci.com/gh/tapglue/multiverse)

## Documentation

See [Documentation](https://github.com/tapglue/multiverse/wiki) for entities, api design and more.

## System Requirements

- go (latest)
- postgres
- postgis
- redis 2.8 or newer

## Install Go

Go is the main language that is being used to develop our backend application.

### Manual Installation

Follow the instruction on the [Go website](https://golang.org/doc/install) to install it manually.

### Homebrew Installation

`brew install go`

### Configure GOPATH

Make sure that your `GOPATH` is configured correctly. You can type `go env` to evaluate your setup.

## Install Redis

Redis is used as a cache to store information such as user sessions.

### Manual Installation

Follow the instruction on the [Redis website](http://redis.io/download) to install it manually.

### Homebrew Installation

`brew install redis`

## Install Postgres

Postgres is the main database that is used to store all data related to:

- orgs
- members
- apps
- users
- connections
- events
- objects

### Manual Installation

Follow the instruction on the [Postgres website](http://www.postgresql.org/download/) to install it manually.

### Homebrew Installation

`brew install postgres`

## Install Postgis

PostGIS is a spatial database extender for Postgres. This can be used to create geo-based-feeds.

### Manual Installation

Follow the instruction on the [Postgis website](http://postgis.net/install/) to install it manually.

### Homebrew Installation

`brew install postgis`

## Getting started

Download the git repository to get started.

```shell
$ mkdir -p $GOPATH/src/github.com/tapglue/multiverse
$ git clone git@github.com:tapglue/multiverse.git
$ cd multiverse
```

## Dependencies

You can install dependencies with `godep` or manually.

### Godep Installation

Go to the root directory and run `godep restore`

All dependencies should be installed into your `GOPATH`.

### Manual Installation

All dependencies should be fetched correctly by running:

```shell
$ go get github.com/tapglue/multiverse
```

or, if you cloned it locally in your GOPATH

```shell
$ cd $GOPATH/src/github.com/tapglue/multiverse
$ go get ./...
```

## Database

We need to start and create the databases before we start our backend application.

### Launch Redis


```shell
redis-server /usr/local/etc/redis.conf
```


### Launch Postgres

```shell
postgres -D /usr/local/var/postgres/ -d 2 -t pl
```

### Create Databases for Tests & Development

```shell
createdb tapglue_test
```

```shell
createdb tapglue_dev
```

### Create Schemas and test data


```shell
psql -E -d tapglue_test -f v04/resources/db/pgsql.sql
```

```shell
psql -E -d tapglue_dev -f v04/resources/db/pgsql.sql
```

## Server configuration

Configure the server including ports and database settings in the [config.json](config.json).

### Test Configuration

Navigate to `v04/server` and create a `config.json`

```json
{
  "env": "dev",
  "listenHost": ":8083",
  "skip_security": false,
  "json_logs": true,
  "use_syslog": false,
  "use_ssl": false,
  "use_low_sec": false,
  "redis": {
    "hosts": [
      "localhost:6379"
    ],
    "password": "",
    "db": 0,
    "pool_size": 30
  },
  "postgres": {
    "database": "tapglue_test",
    "master": {
      "username": "whoami",
      "password": "",
      "host": "127.0.0.1",
      "options": "sslmode=disable&connect_timeout=5"
    },
    "slaves": [
      {
        "username": "whoami",
        "password": "",
        "host": "127.0.0.1",
        "options": "sslmode=disable&connect_timeout=5"
      }
    ]
  }
}
```

The username must be the one you are logged in with `whoami`

### Development Configuration

Navigate to `cmd/intaker` and create a `config.json`

```json
{
  "env": "dev",
  "listenHost": ":8083",
  "skip_security": false,
  "json_logs": true,
  "use_syslog": false,
  "use_ssl": false,
  "use_low_sec": false,
  "redis": {
    "hosts": [
      "localhost:6379"
    ],
    "password": "",
    "db": 0,
    "pool_size": 30
  },
  "postgres": {
    "database": "tapglue_dev",
    "master": {
      "username": "whoami",
      "password": "",
      "host": "127.0.0.1",
      "options": "sslmode=disable&connect_timeout=5"
    },
    "slaves": [
      {
        "username": "whoami",
        "password": "",
        "host": "127.0.0.1",
        "options": "sslmode=disable&connect_timeout=5"
      }
    ]
  }
}

```

## Start server


```shell
go run -tags postgres cmd/intaker/intaker.go
```

## Running tests

Change to `v04/server` and run

```shell
CI=true go test -tags postgres -check.v ./...
```

Other tests from root directory:

```shell
go test -v ./controller/...
```

```shell
go test -v -tags integration ./service/...
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

## Profiling

```shell
$ bin/ab/*.sh
```

## Code commit

Before doing a commit, please run the following in the ```$GOPATH/src/github.com/tapglue/multiverse```
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

### SSH Tunneling to the database

Replace:
- `BASTION_IP` with the IP of the Bastion host
- `PRIVATE_IP` with the IP of the private host
- `DB_IP` with the IP of the database

The command:
```bash
ssh -A user@BASTION_IP -L 54321:localhost:54320 ssh -L 54320:PRIVATE_IP:5432 DB_IP cat -
```

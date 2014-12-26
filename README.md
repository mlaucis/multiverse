Gluee backend [![wercker status](https://app.wercker.com/status/79fb688d3dd5889a31d18cf6fee31a24/s/ "wercker status")](https://app.wercker.com/project/bykey/79fb688d3dd5889a31d18cf6fee31a24)
=====================

This repository contains the implementation of gluees backend.

## Build status

[![wercker status](https://app.wercker.com/status/79fb688d3dd5889a31d18cf6fee31a24/m "wercker status")](https://app.wercker.com/project/bykey/79fb688d3dd5889a31d18cf6fee31a24)

## Documentation

See [Documentation](https://github.com/Gluee/backend/wiki) for entities, api design and more.

## System Requirements

4 CPU Cores and 1GB RAM would be the baseline.

## Installation

Following steps are need to download and install this project.

### Getting started

Download the git repository to get started.

```shell
$ git clone https://github.com/Gluee/backend.git
$ cd backend
```

### Dependencies

MySQL driver

```shell
$ go get github.com/go-sql-driver/mysql
```

Postgres driver

```shell
$ go get github.com/lib/pq
```

Postgres driver

```shell
$ go get github.com/jmoiron/sqlx
```

Extensions to database/sql

```shell
$ go get github.com/gorilla/mux
```

Registry for global request variables

```shell
$ go get github.com/gorilla/context
```

### Configure server

Configure the server including ports and database settings in the [config file](config.json).

```json
{
	"env": "dev",
	"listenHost": ":8082",
	"db": {
		"username": "root",
		"password": "x",
		"database": "gluee",
		"max_idle": 10,
		"max_open": 300,
		"master": {
			"debug": true,
			"host": "127.0.0.1",
			"port": 3306
		},
		"slaves": [
			{
				"debug": true,
				"host": "127.0.0.1",
				"port": 3306
			}
		]
	}
}
```

### Start server

```shell
$ go run backend.go
```

## Tests

```shell
$ go test
```

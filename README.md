# Tapglue backend [![wercker status](https://app.wercker.com/status/37a8675b2ae12075851f297ce6a36ead/s "wercker status")](https://app.wercker.com/project/bykey/37a8675b2ae12075851f297ce6a36ead)

This repository contains the implementation of tapglues backend.

## Build status

[![wercker status](https://app.wercker.com/status/37a8675b2ae12075851f297ce6a36ead/m "wercker status")](https://app.wercker.com/project/bykey/37a8675b2ae12075851f297ce6a36ead)

## Documentation

See [Documentation](https://github.com/tapglue/backend/wiki) for entities, api design and more.

## System Requirements

4 CPU Cores and 1GB RAM would be the baseline.

## Installation

Following steps are need to download and install this project.

### Getting started

Download the git repository to get started.

```shell
$ git clone https://github.com/tapglue/backend.git
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

Finally get, compile and install everything

```shell
$ go get
$ go install
```

### Server configuration

Configure the server including ports and database settings in the [config.json](config.json).

```json
{
	"env": "dev",
	"listenHost": ":8082",
	"db": {
		"username": "tapglue",
		"password": "x",
		"database": "tapglue",
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

### Database configuration

```sql
CREATE SCHEMA `tapglue` ;
```

#### Create database

#### Create tables

Create a database called `tapglue` and execute the SQL [tapglue.sql](https://github.com/tapglue/backend/blob/master/resources/sql/tapglue.sql) to create all tables and settings.

#### Create users

````sql
CREATE USER 'tapglue'@'localhost' IDENTIFIED BY 'x';
GRANT ALL PRIVILEGES ON * . * TO 'tapglue'@'localhost';
```

### Start server

```shell
$ go run backend.go
```

## Tests

```shell
$ go test
```

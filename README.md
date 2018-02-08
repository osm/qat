## About

qat executes SQL queries, suitable for use in shell scripts.

## Supported SQL servers

* Postgres
* MSSQL

## Installation

```sh
go get github.com/osm/qat
```

## Data sources

### Postgres

```sh
qat -source "host=localhost port=5432 user=foo password=bar dbname=baz sslmode=disable" -query "SELECT 1"
```

### MSSQL

```sh
qat -source "server=localhost;port=1433;user id=foo;password=bar" -query "SELECT 1"
```

## Usage

```sh
# Accepts query from stdin
echo "SELECT * FROM foo" | qat -source "host=localhost port=5432 user=foo password=bar dbname=baz sslmode=disable"

# Query passed as a parameter
qat -source "host=localhost port=5432 user=foo password=bar dbname=baz sslmode=disable" -query "SELECT * FROM foo"

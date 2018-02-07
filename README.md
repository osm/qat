## About

qat executes SQL queries, suitable for use in shell scripts.

## Supported SQL servers

* Postgres

## Installation

```sh
go get github.com/osm/qat
```

## Usage

```sh
# Accepts query from stdin
echo "SELECT * FROM foo" | qat -user foo -password bar -port 5432 -server localhost -name foo -query -

# Query passed as a parameter
qat -user foo -password bar -port 5432 -server localhost -name foo -query "SELECT * FROM foo"

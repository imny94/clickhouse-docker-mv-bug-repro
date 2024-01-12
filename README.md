# Clickhouse docker image bug reproduction

## Issue
When using the official clickhouse docker image, the first time that a table is queried, it returns no results. Subsequent queries return the correct results.

## Local Environment
- go version: go1.21.5 darwin/arm64
- docker Engine: 20.10.21
- docker Compose: v2.13.0

## How to reproduce
This issue only happens on the first time that the script is run, so you'll need to create a fresh environment and do the following:

### If running for the *first* time
1. Run `docker-compose up -d`
    - This should create a local instance of clickhouse for you, with the relevant tables set-up (1 raw table, 1 aggregate table, and 1 materialized view)
1. Run `go run main.go`

### Subsequent runs
1. Run `docker-compose down` to stop the clickhouse instance
1. Remove the attached clickhouse volume by deleting the `./out` folder, i.e. `rm -rf ./out`
1. Run `docker-compose up -d`
    - This should create a local instance of clickhouse for you, with the relevant tables set-up (1 raw table, 1 aggregate table, and 1 materialized view)
1. Run `go run main.go`

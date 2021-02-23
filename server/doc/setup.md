# Postgres setup

	createdb -U postgres anigram

# Migrations

Migrations are done using [migrate](https://github.com/golang-migrate/migrate).

## Installation

See the [docs](https://github.com/golang-migrate/migrate/tree/master/cli) for the migrate cli tool.

## Creating a new migration

From the project's root directory:

    migrate create -dir ./migrations -ext sql initial_tables

## Migrating

From the project's root directory:

    migrate -database postgresql://anigram@localhost:5432/anigram?sslmode=disable -path ./migrations up

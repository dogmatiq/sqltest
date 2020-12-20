# SQL Testing Utilities

[![Build Status](https://github.com/dogmatiq/sqltest/workflows/CI/badge.svg)](https://github.com/dogmatiq/sqltest/actions?workflow=CI)
[![Code Coverage](https://img.shields.io/codecov/c/github/dogmatiq/sqltest/main.svg)](https://codecov.io/github/dogmatiq/sqltest)
[![Latest Version](https://img.shields.io/github/tag/dogmatiq/sqltest.svg?label=semver)](https://semver.org)
[![Documentation](https://img.shields.io/badge/go.dev-reference-007d9c)](https://pkg.go.dev/github.com/dogmatiq/sqltest)
[![Go Report Card](https://goreportcard.com/badge/github.com/dogmatiq/sqltest)](https://goreportcard.com/report/github.com/dogmatiq/sqltest)

This is a Go library containing utilities that help when writing tests that use
real SQL database servers.

It is only intended for use as a test dependencies for projects within the
Dogmatiq organization.

The primary feature is the `NewDatabase()` function which creates a temporary
database that can be discarded at the end of each test.

## Database Products

The database products in the table below are currently supported. Some products
are supported via multiple different Go SQL drivers.

Each of these databases can be started using the provided [Docker stack](docker-stack.yml),
and/or started within each project's Github Actions workflow. This project's
[CI workflow](.github/workflows.ci.yml) serves as an example of how to start
each product under Github Actions.

Product     | Supported Drivers             | Notes
------------|-------------------------------|-----------------------------------
MySQL       | [mysql]                       |
MariaDB     | [mysql]                       |
PostgreSQL  | [pgx] (preferred), [postgres] |
SQLite      | [sqlite3]                     | Embedded database, requires CGO

[mysql]: https://github.com/go-sql-driver/mysql
[pgx]: https://github.com/jackc/pgx
[postgres]: https://github.com/lib/pq
[sqlite3]: https://github.com/mattn/go-sqlite3

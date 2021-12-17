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

| Product    | Supported Drivers             | Notes                           |
| ---------- | ----------------------------- | ------------------------------- |
| MySQL      | [mysql]                       | &mdash;                         |
| MariaDB    | [mysql]                       | &mdash;                         |
| PostgreSQL | [pgx] (preferred), [postgres] | &mdash;                         |
| SQLite     | [sqlite3]                     | Embedded database, requires CGO |

[mysql]: https://github.com/go-sql-driver/mysql
[pgx]: https://github.com/jackc/pgx
[postgres]: https://github.com/lib/pq
[sqlite3]: https://github.com/mattn/go-sqlite3

## Docker Stack

The [`docker-stack.yml`](docker-stack.yml) file in this repository starts
services for each of the supported database products. These services are typically required to run the tests for any project that depends on `dogmatiq/sqltest`.

The stack can be deployed using the following command:

```console
curl https://raw.githubusercontent.com/dogmatiq/sqltest/main/docker-stack.yml | docker stack deploy dogmatiq-sqltest --compose-file -
```

## GitHub Actions Configuration

Projects that depend on `dogmatiq/sqltest` should use the `go+sql` workflow. This workflow starts services for each of the supported
database products, and runs the tests both with and without CGO enabled.

The workflow is chosen by changing the repository definition in the [Terraform configuration]. To see an example, check the definition of `dogmatiq/sqltest` itself, which uses this workflow.

[terraform configuration]: https://github.com/dogmatiq/repos

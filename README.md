<div align="center">

# SQL Testing Utilities

A Go library containing utilities that help when writing tests that use real SQL
database servers.

[![Documentation](https://img.shields.io/badge/go.dev-documentation-007d9c?&style=for-the-badge)](https://pkg.go.dev/github.com/dogmatiq/sqltest)
[![Latest Version](https://img.shields.io/github/tag/dogmatiq/sqltest.svg?&style=for-the-badge&label=semver)](https://github.com/dogmatiq/sqltest/releases)
[![Build Status](https://img.shields.io/github/actions/workflow/status/dogmatiq/sqltest/ci.yml?style=for-the-badge&branch=main)](https://github.com/dogmatiq/sqltest/actions/workflows/ci.yml)
[![Code Coverage](https://img.shields.io/codecov/c/github/dogmatiq/sqltest/main.svg?style=for-the-badge)](https://codecov.io/github/dogmatiq/sqltest)

</div>

This library is only intended for use as a test dependencies for projects within
the Dogmatiq organization.

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

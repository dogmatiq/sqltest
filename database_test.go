package sqltest_test

import (
	"context"
	"time"

	"github.com/dogmatiq/dapper"
	. "github.com/dogmatiq/sqltest"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("type Database", func() {
	Describe("func NewDatabase()", func() {
		DescribeTable(
			"it returns DSN that can be used to open a database",
			func(d Driver, p Product) {
				ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
				defer cancel()

				dsn, err := NewDatabase(ctx, d, p)
				Expect(err).ShouldNot(HaveOccurred())
				defer dsn.Close()

				dapper.Print(dsn)

				db, err := dsn.Open()
				Expect(err).ShouldNot(HaveOccurred())

				err = db.PingContext(ctx)
				Expect(err).ShouldNot(HaveOccurred())

				err = dsn.Close()
				Expect(err).ShouldNot(HaveOccurred())

				// Database should be closed now, ping should fail.
				err = db.PingContext(ctx)
				Expect(err).Should(HaveOccurred())
			},
			Entry("MySQL using the 'mysql' driver", MySQLDriver, MySQL),
			Entry("MariaDB using the 'mysql' driver", MySQLDriver, MariaDB),
			Entry("PostgreSQL using the 'pgx' driver", PGXDriver, PostgreSQL),
			Entry("PostgreSQL using the 'postgres' driver", PostgresDriver, PostgreSQL),
			Entry("SQLite using the 'sqlite3' driver", SQLite3Driver, SQLite),
		)
	})

	Describe("func Open()", func() {
		DescribeTable(
			"returns an error",
			func(db *Database) {
				_, err := db.Open()
				Expect(err).To(MatchError("attempted to open a database pool for a closed database"))
			},
			Entry("nil pointer", (*Database)(nil)),
			Entry("zero value", &Database{}),
		)
	})

	Describe("func Close()", func() {
		DescribeTable(
			"does not return an error",
			func(db *Database) {
				err := db.Close()
				Expect(err).ShouldNot(HaveOccurred())
			},
			Entry("nil pointer", (*Database)(nil)),
			Entry("zero value", &Database{}),
		)
	})
})

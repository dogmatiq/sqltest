package sqltest_test

import (
	"context"
	"fmt"
	"time"

	. "github.com/dogmatiq/sqltest"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("func NewTemporaryDatabase()", func() {
	entries := []TableEntry{
		entry(MySQLDriver, MySQL),
		entry(MySQLDriver, MariaDB),
		entry(PGXDriver, PostgreSQL),
		entry(PostgresDriver, PostgreSQL),
	}

	DescribeTable(
		"it returns DSN that can be used to open a database",
		func(d Driver, p Product) {
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()

			dsn, err := NewTemporaryDatabase(ctx, d, p)
			Expect(err).ShouldNot(HaveOccurred())
			defer dsn.Close()

			db, err := dsn.Open()
			Expect(err).ShouldNot(HaveOccurred())

			err = db.PingContext(ctx)
			Expect(err).ShouldNot(HaveOccurred())

			err = dsn.Close()
			Expect(err).ShouldNot(HaveOccurred())
		},
		entries...,
	)
})

func entry(d Driver, p Product) TableEntry {
	return Entry(
		fmt.Sprintf("%s via '%s' driver", p.Name(), d.Name()),
		d,
		p,
	)
}

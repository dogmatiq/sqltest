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

var _ = Describe("func Open()", func() {
	DescribeTable(
		"opens a connection to a database",
		func(d Driver, p Product) {
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()

			db, close, err := Open(ctx, d, p)
			Expect(err).ShouldNot(HaveOccurred())
			defer close()

			err = db.PingContext(ctx)
			Expect(err).ShouldNot(HaveOccurred())

			err = close()
			Expect(err).ShouldNot(HaveOccurred())
		},
		entry(MySQLDriver, MySQL),
	)
})

func entry(d Driver, p Product) TableEntry {
	return Entry(
		fmt.Sprintf("%s via '%s' driver", p.Name(), d.Name()),
		d,
		p,
	)
}

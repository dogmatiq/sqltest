package sqltest_test

import (
	. "github.com/dogmatiq/sqltest"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("func CompatiblePairs()", func() {
	It("returns all pairs by default", func() {
		pairs := CompatiblePairs()

		Expect(pairs).To(ConsistOf(
			Pair{
				Product: MySQL,
				Driver:  MySQLDriver,
			},
			Pair{
				Product: MariaDB,
				Driver:  MySQLDriver,
			},
			Pair{
				Product: PostgreSQL,
				Driver:  PGXDriver,
			},
			Pair{
				Product: PostgreSQL,
				Driver:  PostgresDriver,
			},
			Pair{
				Product: SQLite,
				Driver:  SQLite3Driver,
			},
		))
	})

	It("limits results to the provided products", func() {
		pairs := CompatiblePairs(MariaDB, PostgreSQL)

		Expect(pairs).To(ConsistOf(
			Pair{
				Product: MariaDB,
				Driver:  MySQLDriver,
			},
			Pair{
				Product: PostgreSQL,
				Driver:  PGXDriver,
			},
			Pair{
				Product: PostgreSQL,
				Driver:  PostgresDriver,
			},
		))
	})
})

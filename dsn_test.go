package sqltest_test

import (
	. "github.com/dogmatiq/sqltest"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("type DSN", func() {
	Describe("func Close()", func() {
		It("can be called on a nil *DSN", func() {
			var dsn *DSN

			err := dsn.Close()
			Expect(err).ShouldNot(HaveOccurred())
		})

		It("can be called on a zero-value DSN", func() {
			dsn := &DSN{}

			err := dsn.Close()
			Expect(err).ShouldNot(HaveOccurred())
		})
	})
})

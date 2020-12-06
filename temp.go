package sqltest

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strings"
	"time"

	"go.uber.org/multierr"
)

// NewTemporaryDatabase creates a temporary test database for a specific pair of
// database product and SQL driver and returns its DSN.
//
// The DSN is an io.Closer that must be closed when the database is no longer
// needed. It is safe to close the DSN even if this function returns an error.
func NewTemporaryDatabase(
	ctx context.Context,
	d Driver,
	p Product,
) (_ *DSN, err error) {
	defaultTemplateDSN, ok := p.TemplateDSN(d)
	if !ok {
		return nil, fmt.Errorf(
			"%s can not be accessed using the '%s' driver",
			p.Name(),
			d.Name(),
		)
	}

	templateDSN := os.Getenv(
		strings.ToUpper(
			fmt.Sprintf("DOGMATIQ_TEST_DSN_%s_%s", p.Name(), d.Name()),
		),
	)
	if templateDSN == "" {
		templateDSN = defaultTemplateDSN
	}

	prefix, err := d.DatabaseNameFromDSN(templateDSN)
	if err != nil {
		return nil, fmt.Errorf(
			"unable to extract database name from %s DSN (%s) using the '%s' driver: %w",
			p.Name(),
			templateDSN,
			d.Name(),
			err,
		)
	}

	tempDB := generateDatabaseName(prefix)
	tempDSN, err := d.DSNForTesting(templateDSN, tempDB)
	if err != nil {
		return nil, fmt.Errorf(
			"unable to generate a %s DSN for the test database using the '%s' driver: %w",
			p.Name(),
			d.Name(),
			err,
		)
	}
	defer ifErr(&err, func() {
		tempDSN.Close()
	})

	schemaDSN, err := d.DSNForSchemaManipulation(templateDSN)
	if err != nil {
		return nil, fmt.Errorf(
			"unable to generate a %s DSN for manipulating the schema using the '%s' driver: %w",
			p.Name(),
			d.Name(),
			err,
		)
	}
	defer ifErr(&err, func() {
		schemaDSN.Close()
	})

	db, err := sql.Open(d.Name(), schemaDSN.DataSource)
	if err != nil {
		return nil, fmt.Errorf(
			"unable to open %s database (%s) using the '%s' driver: %w",
			p.Name(),
			schemaDSN.DataSource,
			d.Name(),
			err,
		)
	}
	defer ifErr(&err, func() {
		db.Close()
	})

	if err := p.CreateDatabase(ctx, db, tempDB); err != nil {
		return nil, fmt.Errorf(
			"unable to create temporary %s database using the '%s' driver: %w",
			p.Name(),
			d.Name(),
			err,
		)
	}

	return &DSN{
		Driver:     tempDSN.Driver,
		DataSource: tempDSN.DataSource,
		Closer: func(*DSN) error {
			// Allow 3 seconds to drop the temporary database.
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()

			err := p.DropDatabase(ctx, db, tempDB)
			err = multierr.Append(err, db.Close())
			err = multierr.Append(err, schemaDSN.Close())
			err = multierr.Append(err, tempDSN.Close())

			return err
		},
	}, nil
}

// generateDatabaseName returns a name for a temporary test database based on
// the current time and the process PID.
func generateDatabaseName(prefix string) string {
	if prefix == "" {
		prefix = "dogmatiq"
	}

	now := time.Now()

	return fmt.Sprintf(
		"%s_%s_PID%d_%d",
		prefix,
		now.Format("20060102_150405"),
		os.Getpid(),
		now.UnixNano(),
	)
}

// ifErr calls fn if *err is a non-nil error.
func ifErr(err *error, fn func()) {
	if *err != nil {
		fn()
	}
}

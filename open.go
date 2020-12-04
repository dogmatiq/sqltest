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

// Open creates a temporary test database for the specific pair of database
// product and SQL driver and returns a connection pool.
//
// The close function must be called when the database is no longer needed. It
// is always a valid function.
//
// If err is nil, close() drops the temporary database and closes the
// connection(s); otherwise it is a no-op. This allows close() to be called
// in AfterEach() blocks without first checking if it exists.
func Open(
	ctx context.Context,
	d Driver,
	p Product,
) (
	db *sql.DB,
	close func() error,
	err error,
) {
	dsn, closeDSN, err := DSN(ctx, d, p)
	if err != nil {
		return nil, nil, err
	}

	db, err = sql.Open(d.Name(), dsn)
	if err != nil {
		closeDSN()
		return nil, nil, err
	}

	return db, func() error {
		dbErr := db.Close()
		dsnErr := closeDSN()
		return multierr.Append(dbErr, dsnErr)
	}, nil
}

// DSN creates a temporary test database for a specific pair of database product
// and SQL driver and returns its DSN.
//
// If err is nil, close() drops the temporary database; otherwise it is a no-op.
// This allows close() to be called in AfterEach() blocks without first checking
// if it exists.
func DSN(
	ctx context.Context,
	d Driver,
	p Product,
) (
	dsn string,
	close func() error,
	err error,
) {
	defaultTemplateDSN, ok := p.TemplateDSN(d)
	if !ok {
		return "", noop, fmt.Errorf(
			"%s can not be accessed using the the '%s' driver",
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
		return "", noop, fmt.Errorf(
			"unable to extract database name from %s DSN (%s) using the '%s' driver: %w",
			p.Name(),
			templateDSN,
			d.Name(),
			err,
		)
	}

	tempDB := generateDatabaseName(prefix)
	dsn, err = d.DSNForTesting(templateDSN, tempDB)
	if err != nil {
		return "", noop, fmt.Errorf(
			"unable to generate a %s DSN for the test database using the '%s' driver: %w",
			p.Name(),
			d.Name(),
			err,
		)
	}

	schemaDSN, err := d.DSNForSchemaManipulation(templateDSN)
	if err != nil {
		return "", noop, fmt.Errorf(
			"unable to generate a %s DSN for manipulating the schema using the '%s' driver: %w",
			p.Name(),
			d.Name(),
			err,
		)
	}

	db, err := sql.Open(d.Name(), schemaDSN)
	if err != nil {
		return "", noop, fmt.Errorf(
			"unable to connect to %s using the '%s' driver: %w",
			p.Name(),
			d.Name(),
			err,
		)
	}
	defer func() {
		// Close the database if we returned an error to the user.
		if err != nil {
			db.Close()
		}
	}()

	if err := p.CreateDatabase(ctx, db, tempDB); err != nil {
		return "", noop, fmt.Errorf(
			"unable to create temporary %s database using the '%s' driver: %w",
			p.Name(),
			d.Name(),
			err,
		)
	}

	return dsn, func() error {
		// Allow 3 seconds to drop the temporary database.
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		return p.DropDatabase(ctx, db, tempDB)
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

func noop() error { return nil }

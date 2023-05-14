package postgres

import (
	"context"
	"database/sql/driver"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Query struct {
	Name     string
	QueryRaw string
}

type Closer interface {
	Close() error
}

type QueryExecer interface {

	// Exec acquires a connection from the Pool and executes the given SQL.
	// SQL can be either a prepared statement name or an SQL string.
	// Arguments should be referenced positionally from the SQL string as $1, $2, etc.
	// The acquired connection is returned to the pool when the Exec function returns.
	Exec(ctx context.Context, q Query, args ...interface{}) (pgconn.CommandTag, error)

	// Query acquires a connection and executes a query that returns pgx.Rows.
	// Arguments should be referenced positionally from the SQL string as $1, $2, etc.
	// See pgx.Rows documentation to close the returned Rows and return the acquired connection to the Pool.
	//
	// If there is an error, the returned pgx.Rows will be returned in an error state.
	// If preferred, ignore the error returned from Query and handle errors using the returned pgx.Rows.
	//
	// For extra control over how the query is executed, the types QuerySimpleProtocol, QueryResultFormats, and
	// QueryResultFormatsByOID may be used as the first args to control exactly how the query is executed. This is rarely
	// needed. See the documentation for those types for details.
	Query(ctx context.Context, q Query, args ...interface{}) (pgx.Rows, error)

	// QueryRow acquires a connection and executes a query that is expected
	// to return at most one row (pgx.Row). Errors are deferred until pgx.Row's
	// Scan method is called. If the query selects no rows, pgx.Row's Scan will
	// return ErrNoRows. Otherwise, pgx.Row's Scan scans the first selected row
	// and discards the rest. The acquired connection is returned to the Pool when
	// pgx.Row's Scan method is called.
	//
	// Arguments should be referenced positionally from the SQL string as $1, $2, etc.
	//
	// For extra control over how the query is executed, the types QuerySimpleProtocol, QueryResultFormats, and
	// QueryResultFormatsByOID may be used as the first args to control exactly how the query is executed. This is rarely
	// needed. See the documentation for those types for details.
	QueryRow(ctx context.Context, q Query, args ...interface{}) pgx.Row
}

type NamedExecer interface {
	Get(ctx context.Context, dest interface{}, q Query, args ...interface{}) error
	Select(ctx context.Context, dest interface{}, q Query, args ...interface{}) error
}

type Postgres interface {
	driver.Pinger
	Closer

	QueryExecer
	NamedExecer
}

type postgres struct {
	pgxPool *pgxpool.Pool
}

func (p *postgres) Ping(ctx context.Context) error {
	return p.pgxPool.Ping(ctx)
}

func (p *postgres) Close() error {
	p.pgxPool.Close()
	return nil
}

func (p *postgres) Exec(ctx context.Context, q Query, args ...interface{}) (pgconn.CommandTag, error) {
	return p.pgxPool.Exec(ctx, q.QueryRaw, args...)
}

func (p *postgres) Query(ctx context.Context, q Query, args ...interface{}) (pgx.Rows, error) {
	return p.pgxPool.Query(ctx, q.QueryRaw, args...)
}

func (p *postgres) QueryRow(ctx context.Context, q Query, args ...interface{}) pgx.Row {
	return p.pgxPool.QueryRow(ctx, q.QueryRaw, args...)
}

// Get ScanOne is a package-level helper function that uses the DefaultAPI object.
// See API.ScanOne for details.
func (p *postgres) Get(ctx context.Context, dest interface{}, q Query, args ...interface{}) error {
	rows, err := p.Query(ctx, q, args...)
	if err != nil {
		return err
	}

	return pgxscan.ScanOne(dest, rows)
}

func (p *postgres) Select(ctx context.Context, dest interface{}, q Query, args ...interface{}) error {
	rows, err := p.Query(ctx, q, args...)
	if err != nil {
		return err
	}

	return pgxscan.ScanAll(dest, rows)
}

package database

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// Client interface providing access to the database.
type Client interface {
	DB() DB
	Close() error
}

// DB interface combining all database capabilities.
type DB interface {
	SQLExecer
	Pinger
	Transactor
	Close()
}

// TxManager interface for handling transaction scopes.
type TxManager interface {
	ReadCommitted(ctx context.Context, f Handler) error
}

// Handler is a function executed within a transaction.
type Handler func(ctx context.Context) error

// Transactor interface for starting transactions.
type Transactor interface {
	BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error)
}

// SQLExecer interface combining named and basic query execution.
type SQLExecer interface {
	NamedExecer
	QueryExecer
}

// NamedExecer interface for executing queries with named parameters.
type NamedExecer interface {
	ScanOneContext(ctx context.Context, dest interface{}, q Query, args ...interface{}) error
	ScanAllContext(ctx context.Context, dest interface{}, q Query, args ...interface{}) error
}

// QueryExecer interface for executing basic queries.
type QueryExecer interface {
	ExecContext(ctx context.Context, q Query, args ...interface{}) (pgconn.CommandTag, error)
	QueryContext(ctx context.Context, q Query, args ...interface{}) (pgx.Rows, error)
	QueryRowContext(ctx context.Context, q Query, args ...interface{}) pgx.Row
}

// Pinger interface for checking database connectivity.
type Pinger interface {
	Ping(ctx context.Context) error
}

// Query struct representing a database query.
type Query struct {
	Name     string
	QueryRaw string
}

type key string

const TxKey key = "tx"
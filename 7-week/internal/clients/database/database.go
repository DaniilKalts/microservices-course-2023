package database

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

var (
	ErrNotFound        = errors.New("no rows found")
	ErrUniqueViolation = errors.New("unique constraint violation")
)

type Query struct {
	Name     string
	QueryRaw string
}

type Client interface {
	DB() DB
	Close() error
}

type DB interface {
	SQLExecutor
	Pinger
	Transactor
	Close()
}

type SQLExecutor interface {
	NamedExecutor
	QueryExecutor
}

type NamedExecutor interface {
	ScanOneContext(ctx context.Context, dest interface{}, q Query, args ...interface{}) error
	ScanAllContext(ctx context.Context, dest interface{}, q Query, args ...interface{}) error
}

type QueryExecutor interface {
	ExecContext(ctx context.Context, q Query, args ...interface{}) (pgconn.CommandTag, error)
	QueryContext(ctx context.Context, q Query, args ...interface{}) (pgx.Rows, error)
	QueryRowContext(ctx context.Context, q Query, args ...interface{}) pgx.Row
}

type Pinger interface {
	Ping(ctx context.Context) error
}

type Transactor interface {
	BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error)
}

type TxManager interface {
	ReadCommitted(ctx context.Context, f TxFunc) error
}

type TxFunc func(ctx context.Context) error

type key string

const TxKey key = "tx"

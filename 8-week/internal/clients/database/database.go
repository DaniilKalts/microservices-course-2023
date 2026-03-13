package database

import (
	"context"
	"errors"
)

// --- Errors ---

var (
	ErrNotFound        = errors.New("no rows found")
	ErrUniqueViolation = errors.New("unique constraint violation")
	ErrTimeout         = errors.New("query timeout exceeded")
)

// --- Types ---

type Query struct {
	Name     string
	QueryRaw string
}

type ExecResult interface {
	RowsAffected() int64
}

// --- Core interfaces ---

// Client is the top-level entry point to the database layer.
type Client interface {
	DB() DB
	Close() error
}

// DB combines query execution, health checks, and transactions.
type DB interface {
	SQLExecutor
	Pinger
	Transactor
	Close()
}

// --- Query execution ---

type SQLExecutor interface {
	NamedExecutor
	QueryExecutor
}

type NamedExecutor interface {
	ScanOneContext(ctx context.Context, dest interface{}, q Query, args ...interface{}) error
	ScanAllContext(ctx context.Context, dest interface{}, q Query, args ...interface{}) error
}

type QueryExecutor interface {
	ExecContext(ctx context.Context, q Query, args ...interface{}) (ExecResult, error)
}

type Pinger interface {
	Ping(ctx context.Context) error
}

// --- Transactions ---

type IsoLevel string

const (
	IsoLevelReadCommitted  IsoLevel = "read committed"
	IsoLevelRepeatableRead IsoLevel = "repeatable read"
	IsoLevelSerializable   IsoLevel = "serializable"
)

type TxOptions struct {
	IsoLevel IsoLevel
}

type Tx interface {
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}

type Transactor interface {
	BeginTx(ctx context.Context, txOptions TxOptions) (Tx, error)
}

type TxManager interface {
	ReadCommitted(ctx context.Context, f TxFunc) error
}

type TxFunc func(ctx context.Context) error

type key string

const TxKey key = "tx"

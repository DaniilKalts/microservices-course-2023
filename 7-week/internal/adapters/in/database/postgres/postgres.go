package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	"github.com/DaniilKalts/microservices-course-2023/7-week/internal/clients/database"
	"github.com/DaniilKalts/microservices-course-2023/7-week/internal/clients/database/prettier"
)

type pg struct {
	pool   *pgxpool.Pool
	logger *zap.Logger
}

func NewDB(pool *pgxpool.Pool, logger *zap.Logger) (database.DB, error) {
	if pool == nil {
		return nil, errors.New("postgres pool is nil")
	}
	if logger == nil {
		return nil, errors.New("postgres logger is nil")
	}

	return &pg{
		pool:   pool,
		logger: logger,
	}, nil
}

func (p *pg) ScanOneContext(ctx context.Context, dest interface{}, q database.Query, args ...interface{}) error {
	row, err := p.QueryContext(ctx, q, args...)
	if err != nil {
		return err
	}

	return pgxscan.ScanOne(dest, row)
}

func (p *pg) ScanAllContext(ctx context.Context, dest interface{}, q database.Query, args ...interface{}) error {
	rows, err := p.QueryContext(ctx, q, args...)
	if err != nil {
		return err
	}

	return pgxscan.ScanAll(dest, rows)
}

func (p *pg) ExecContext(ctx context.Context, q database.Query, args ...interface{}) (pgconn.CommandTag, error) {
	startedAt := time.Now()

	var (
		tag pgconn.CommandTag
		err error
	)

	tx, ok := ctx.Value(database.TxKey).(pgx.Tx)
	if ok {
		tag, err = tx.Exec(ctx, q.QueryRaw, args...)
	} else {
		tag, err = p.pool.Exec(ctx, q.QueryRaw, args...)
	}

	p.logQuery("exec", q, args, time.Since(startedAt), err)

	return tag, err
}

func (p *pg) QueryContext(ctx context.Context, q database.Query, args ...interface{}) (pgx.Rows, error) {
	startedAt := time.Now()

	var (
		rows pgx.Rows
		err  error
	)

	tx, ok := ctx.Value(database.TxKey).(pgx.Tx)
	if ok {
		rows, err = tx.Query(ctx, q.QueryRaw, args...)
	} else {
		rows, err = p.pool.Query(ctx, q.QueryRaw, args...)
	}

	p.logQuery("query", q, args, time.Since(startedAt), err)

	return rows, err
}

func (p *pg) QueryRowContext(ctx context.Context, q database.Query, args ...interface{}) pgx.Row {
	startedAt := time.Now()

	var row pgx.Row

	tx, ok := ctx.Value(database.TxKey).(pgx.Tx)
	if ok {
		row = tx.QueryRow(ctx, q.QueryRaw, args...)
	} else {
		row = p.pool.QueryRow(ctx, q.QueryRaw, args...)
	}

	p.logQuery("query_row", q, args, time.Since(startedAt), nil)

	return row
}

func (p *pg) Ping(ctx context.Context) error {
	return p.pool.Ping(ctx)
}

func (p *pg) BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error) {
	return p.pool.BeginTx(ctx, txOptions)
}

func (p *pg) Close() {
	p.pool.Close()
}

func (p *pg) logQuery(operation string, q database.Query, args []interface{}, duration time.Duration, err error) {
	fields := []zap.Field{
		zap.String("operation", operation),
		zap.String("query_name", q.Name),
		zap.Int("args_count", len(args)),
		zap.Float64("duration_ms", float64(duration)/float64(time.Millisecond)),
	}

	if p.logger.Core().Enabled(zap.DebugLevel) {
		fields = append(fields, zap.String("query", prettier.Pretty(q.QueryRaw, prettier.PlaceholderDollar, args...)))
	}

	if err != nil {
		fields = append(fields, zap.Error(err))
		p.logger.Error("database operation failed", fields...)
		return
	}

	p.logger.Debug("database operation completed", fields...)
}

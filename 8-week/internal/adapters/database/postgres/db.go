package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"go.uber.org/zap"

	"github.com/DaniilKalts/microservices-course-2023/8-week/internal/clients/database"
	"github.com/DaniilKalts/microservices-course-2023/8-week/internal/clients/database/prettier"
)

type db struct {
	pool         *pgxpool.Pool
	queryTimeout time.Duration
	logger       *zap.Logger
}

func NewDB(pool *pgxpool.Pool, queryTimeout time.Duration, logger *zap.Logger) (database.DB, error) {
	if pool == nil {
		return nil, errors.New("postgres pool is nil")
	}
	if logger == nil {
		return nil, errors.New("postgres logger is nil")
	}

	return &db{
		pool:         pool,
		queryTimeout: queryTimeout,
		logger:       logger,
	}, nil
}

func (p *db) ScanOneContext(ctx context.Context, dest any, q database.Query, args ...any) error {
	ctx, cancel := context.WithTimeout(ctx, p.queryTimeout)
	defer cancel()

	row, err := p.queryContext(ctx, q, args...)
	if err != nil {
		return err
	}

	if err = pgxscan.ScanOne(dest, row); err != nil {
		if pgxscan.NotFound(err) {
			return database.ErrNotFound
		}
		return translateErr(err)
	}

	return nil
}

func (p *db) ScanAllContext(ctx context.Context, dest any, q database.Query, args ...any) error {
	ctx, cancel := context.WithTimeout(ctx, p.queryTimeout)
	defer cancel()

	rows, err := p.queryContext(ctx, q, args...)
	if err != nil {
		return err
	}

	return pgxscan.ScanAll(dest, rows)
}

func (p *db) ExecContext(ctx context.Context, q database.Query, args ...any) (database.ExecResult, error) {
	ctx, cancel := context.WithTimeout(ctx, p.queryTimeout)
	defer cancel()

	span, ctx := p.startDBSpan(ctx, "db.exec", q)
	defer span.Finish()

	startedAt := time.Now()

	var (
		tag execResult
		err error
	)

	tx, ok := ctx.Value(database.TxKey).(pgx.Tx)
	if ok {
		tag.ct, err = tx.Exec(ctx, q.QueryRaw, args...)
	} else {
		tag.ct, err = p.pool.Exec(ctx, q.QueryRaw, args...)
	}

	p.logQuery("exec", q, args, time.Since(startedAt), err)
	if err != nil {
		ext.Error.Set(span, true)
		span.LogKV("event", "error", "message", err.Error())
		return tag, translateErr(err)
	}

	return tag, nil
}

// queryContext is an internal method used by ScanOneContext and ScanAllContext.
func (p *db) queryContext(ctx context.Context, q database.Query, args ...any) (pgx.Rows, error) {
	span, ctx := p.startDBSpan(ctx, "db.query", q)
	defer span.Finish()

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
	if err != nil {
		ext.Error.Set(span, true)
		span.LogKV("event", "error", "message", err.Error())
		return rows, translateErr(err)
	}

	return rows, nil
}

func (p *db) Ping(ctx context.Context) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "db.ping")
	defer span.Finish()

	span.SetTag("component", "database")
	span.SetTag("db.system", "postgresql")

	err := p.pool.Ping(ctx)
	if err != nil {
		p.logger.Error("failed to ping database", zap.Error(err))
		ext.Error.Set(span, true)
		span.LogKV("event", "error", "message", err.Error())
	}

	return err
}

func (p *db) BeginTx(ctx context.Context, txOptions database.TxOptions) (database.Tx, error) {
	ctx, cancel := context.WithTimeout(ctx, p.queryTimeout)
	defer cancel()

	span, ctx := opentracing.StartSpanFromContext(ctx, "db.begin_tx")
	defer span.Finish()

	span.SetTag("component", "database")
	span.SetTag("db.system", "postgresql")

	pgxOpts := pgx.TxOptions{IsoLevel: pgx.TxIsoLevel(string(txOptions.IsoLevel))}

	tx, err := p.pool.BeginTx(ctx, pgxOpts)
	if err != nil {
		p.logger.Error("failed to begin transaction", zap.String("isolation_level", string(txOptions.IsoLevel)), zap.Error(err))
		ext.Error.Set(span, true)
		span.LogKV("event", "error", "message", err.Error())
		return nil, err
	}

	p.logger.Debug("transaction started", zap.String("isolation_level", string(txOptions.IsoLevel)))

	return tx, nil
}

func (p *db) Close() {
	p.pool.Close()
}

func (p *db) logQuery(operation string, q database.Query, args []any, duration time.Duration, err error) {
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

func (p *db) startDBSpan(ctx context.Context, operationName string, q database.Query) (opentracing.Span, context.Context) {
	span, spanCtx := opentracing.StartSpanFromContext(ctx, operationName)
	span.SetTag("component", "database")
	span.SetTag("db.system", "postgresql")
	span.SetTag("db.query_name", q.Name)

	return span, spanCtx
}

// execResult wraps pgconn.CommandTag to satisfy database.ExecResult.
type execResult struct {
	ct interface {
		RowsAffected() int64
	}
}

func (r execResult) RowsAffected() int64 {
	if r.ct == nil {
		return 0
	}
	return r.ct.RowsAffected()
}

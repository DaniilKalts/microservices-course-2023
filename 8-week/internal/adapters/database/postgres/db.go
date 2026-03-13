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

func (d *db) ScanOneContext(ctx context.Context, dest any, q database.Query, args ...any) error {
	ctx, cancel := context.WithTimeout(ctx, d.queryTimeout)
	defer cancel()

	row, err := d.queryContext(ctx, q, args...)
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

func (d *db) ScanAllContext(ctx context.Context, dest any, q database.Query, args ...any) error {
	ctx, cancel := context.WithTimeout(ctx, d.queryTimeout)
	defer cancel()

	rows, err := d.queryContext(ctx, q, args...)
	if err != nil {
		return err
	}

	return pgxscan.ScanAll(dest, rows)
}

func (d *db) ExecContext(ctx context.Context, q database.Query, args ...any) (database.ExecResult, error) {
	ctx, cancel := context.WithTimeout(ctx, d.queryTimeout)
	defer cancel()

	span, ctx := d.startDBSpan(ctx, "db.exec", q)
	defer span.Finish()

	startedAt := time.Now()

	var (
		tag execResult
		err error
	)

	tx, hasTx := ctx.Value(database.TxKey).(pgx.Tx)
	if hasTx {
		tag.ct, err = tx.Exec(ctx, q.QueryRaw, args...)
	} else {
		tag.ct, err = d.pool.Exec(ctx, q.QueryRaw, args...)
	}

	d.logQuery("exec", q, args, time.Since(startedAt), err)
	if err != nil {
		ext.Error.Set(span, true)
		span.LogKV("event", "error", "message", err.Error())
		return tag, translateErr(err)
	}

	return tag, nil
}

func (d *db) queryContext(ctx context.Context, q database.Query, args ...any) (pgx.Rows, error) {
	span, ctx := d.startDBSpan(ctx, "db.query", q)
	defer span.Finish()

	startedAt := time.Now()

	var (
		rows pgx.Rows
		err  error
	)

	tx, hasTx := ctx.Value(database.TxKey).(pgx.Tx)
	if hasTx {
		rows, err = tx.Query(ctx, q.QueryRaw, args...)
	} else {
		rows, err = d.pool.Query(ctx, q.QueryRaw, args...)
	}

	d.logQuery("query", q, args, time.Since(startedAt), err)
	if err != nil {
		ext.Error.Set(span, true)
		span.LogKV("event", "error", "message", err.Error())
		return rows, translateErr(err)
	}

	return rows, nil
}

func (d *db) Ping(ctx context.Context) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "db.ping")
	defer span.Finish()

	span.SetTag("component", "database")
	span.SetTag("db.system", "postgresql")

	err := d.pool.Ping(ctx)
	if err != nil {
		d.logger.Error("failed to ping database", zap.Error(err))
		ext.Error.Set(span, true)
		span.LogKV("event", "error", "message", err.Error())
	}

	return err
}

func (d *db) BeginTx(ctx context.Context, txOptions database.TxOptions) (database.Tx, error) {
	ctx, cancel := context.WithTimeout(ctx, d.queryTimeout)
	defer cancel()

	span, ctx := opentracing.StartSpanFromContext(ctx, "db.begin_tx")
	defer span.Finish()

	span.SetTag("component", "database")
	span.SetTag("db.system", "postgresql")

	pgxOpts := pgx.TxOptions{IsoLevel: pgx.TxIsoLevel(string(txOptions.IsoLevel))}

	tx, err := d.pool.BeginTx(ctx, pgxOpts)
	if err != nil {
		d.logger.Error("failed to begin transaction", zap.String("isolation_level", string(txOptions.IsoLevel)), zap.Error(err))
		ext.Error.Set(span, true)
		span.LogKV("event", "error", "message", err.Error())
		return nil, err
	}

	d.logger.Debug("transaction started", zap.String("isolation_level", string(txOptions.IsoLevel)))

	return tx, nil
}

func (d *db) Close() {
	d.pool.Close()
}

func (d *db) logQuery(operation string, q database.Query, args []any, duration time.Duration, err error) {
	fields := []zap.Field{
		zap.String("operation", operation),
		zap.String("query_name", q.Name),
		zap.Int("args_count", len(args)),
		zap.Float64("duration_ms", duration.Seconds()*1000),
	}

	if d.logger.Core().Enabled(zap.DebugLevel) {
		fields = append(fields, zap.String("query", prettier.Pretty(q.QueryRaw, prettier.PlaceholderDollar, args...)))
	}

	if err != nil {
		fields = append(fields, zap.Error(err))
		d.logger.Error("database operation failed", fields...)
		return
	}

	d.logger.Debug("database operation completed", fields...)
}

func (d *db) startDBSpan(ctx context.Context, operationName string, q database.Query) (opentracing.Span, context.Context) {
	span, spanCtx := opentracing.StartSpanFromContext(ctx, operationName)
	span.SetTag("component", "database")
	span.SetTag("db.system", "postgresql")
	span.SetTag("db.query_name", q.Name)

	return span, spanCtx
}

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

package transaction

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/DaniilKalts/microservices-course-2023/8-week/internal/clients/database"
)

type manager struct {
	db     database.Transactor
	logger *zap.Logger
}

func NewTransactionManager(db database.Transactor, logger *zap.Logger) database.TxManager {
	return &manager{
		db:     db,
		logger: logger,
	}
}

func (m *manager) transaction(ctx context.Context, opts pgx.TxOptions, fn database.TxFunc) (err error) {
	if _, ok := ctx.Value(database.TxKey).(pgx.Tx); ok {
		return fn(ctx)
	}

	tx, err := m.db.BeginTx(ctx, opts)
	if err != nil {
		return errors.Wrap(err, "failed to begin transaction")
	}

	ctx = context.WithValue(ctx, database.TxKey, tx)

	defer func() {
		if r := recover(); r != nil {
			err = errors.Errorf("panic recovered: %v", r)
			m.logger.Error("panic recovered in transaction", zap.Any("panic", r))
		}
		if err != nil {
			if errRollback := tx.Rollback(ctx); errRollback != nil {
				m.logger.Error("failed to rollback transaction", zap.Error(errRollback), zap.NamedError("original_error", err))
				err = errors.Wrap(err, errRollback.Error())
			}
			return
		}
		if err = tx.Commit(ctx); err != nil {
			m.logger.Error("failed to commit transaction", zap.Error(err))
			err = errors.Wrap(err, "failed to commit transaction")
		}
	}()

	if err = fn(ctx); err != nil {
		err = errors.Wrap(err, "failed to execute transaction")
	}

	return err
}

func (m *manager) ReadCommitted(ctx context.Context, f database.TxFunc) error {
	txOpts := pgx.TxOptions{IsoLevel: pgx.ReadCommitted}
	return m.transaction(ctx, txOpts, f)
}

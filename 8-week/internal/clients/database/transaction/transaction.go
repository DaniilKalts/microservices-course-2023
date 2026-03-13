package transaction

import (
	"context"
	"fmt"

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

func (m *manager) transaction(ctx context.Context, opts database.TxOptions, fn database.TxFunc) (err error) {
	if _, ok := ctx.Value(database.TxKey).(database.Tx); ok {
		return fn(ctx)
	}

	tx, err := m.db.BeginTx(ctx, opts)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	ctx = context.WithValue(ctx, database.TxKey, tx)

	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic recovered: %v", r)
			m.logger.Error("panic recovered in transaction", zap.Any("panic", r))
		}
		if err != nil {
			if errRollback := tx.Rollback(ctx); errRollback != nil {
				m.logger.Error("failed to rollback transaction", zap.Error(errRollback), zap.NamedError("original_error", err))
				err = fmt.Errorf("%w: %v", err, errRollback)
			}
			return
		}
		if err = tx.Commit(ctx); err != nil {
			m.logger.Error("failed to commit transaction", zap.Error(err))
			err = fmt.Errorf("failed to commit transaction: %w", err)
		}
	}()

	if err = fn(ctx); err != nil {
		err = fmt.Errorf("failed to execute transaction: %w", err)
	}

	return err
}

func (m *manager) ReadCommitted(ctx context.Context, f database.TxFunc) error {
	txOpts := database.TxOptions{IsoLevel: database.IsoLevelReadCommitted}
	return m.transaction(ctx, txOpts, f)
}

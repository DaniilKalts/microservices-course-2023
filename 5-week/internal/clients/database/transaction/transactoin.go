package transaction

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"

	"github.com/DaniilKalts/microservices-course-2023/5-week/internal/clients/database"
)

type manager struct {
	db database.Transactor
}

func NewTransactionManager(db database.Transactor) database.TxManager {
	return &manager{
		db: db,
	}
}

func (m *manager) transaction(ctx context.Context, opts pgx.TxOptions, fn database.Handler) (err error) {
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
		}
		if err != nil {
			if errRollback := tx.Rollback(ctx); errRollback != nil {
				err = errors.Wrap(err, errRollback.Error())
			}
			return
		}
		if err = tx.Commit(ctx); err != nil {
			err = errors.Wrap(err, "failed to commit transaction")
		}
	}()

	if err = fn(ctx); err != nil {
		err = errors.New("failed to execute transaction")
	}

	return nil
}

func (m *manager) ReadCommitted(ctx context.Context, f database.Handler) error {
	txOpts := pgx.TxOptions{IsoLevel: pgx.ReadCommitted}
	return m.transaction(ctx, txOpts, f)
}

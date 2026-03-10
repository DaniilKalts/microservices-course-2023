package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	"github.com/DaniilKalts/microservices-course-2023/8-week/internal/clients/database"
)

type pgClient struct {
	db database.DB
}

func New(ctx context.Context, dsn string, logger *zap.Logger) (database.Client, error) {
	if logger == nil {
		return nil, errors.New("postgres logger is nil")
	}

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to db: %w", err)
	}

	db, err := NewDB(pool, logger)
	if err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to init postgres db wrapper: %w", err)
	}

	return &pgClient{
		db: db,
	}, nil
}

func (c *pgClient) DB() database.DB {
	return c.db
}

func (c *pgClient) Close() error {
	if c.db != nil {
		c.db.Close()
	}
	return nil
}

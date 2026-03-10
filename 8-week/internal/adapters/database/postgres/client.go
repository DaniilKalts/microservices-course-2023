package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	"github.com/DaniilKalts/microservices-course-2023/8-week/internal/clients/database"
	"github.com/DaniilKalts/microservices-course-2023/8-week/internal/config"
)

type client struct {
	db database.DB
}

func New(ctx context.Context, cfg config.PostgresConfig, logger *zap.Logger) (database.Client, error) {
	if logger == nil {
		return nil, errors.New("postgres logger is nil")
	}

	pool, err := pgxpool.New(ctx, cfg.DSN())
	if err != nil {
		return nil, fmt.Errorf("failed to connect to db: %w", err)
	}

	db, err := NewDB(pool, cfg.QueryTimeout, logger)
	if err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to init postgres db wrapper: %w", err)
	}

	return &client{
		db: db,
	}, nil
}

func (c *client) DB() database.DB {
	return c.db
}

func (c *client) Close() error {
	if c.db != nil {
		c.db.Close()
	}
	return nil
}

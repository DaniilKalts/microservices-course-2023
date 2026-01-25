package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/DaniilKalts/microservices-course-2023/3-week/internal/clients/database"
)

type pgClient struct {
	db database.DB
}

func New(ctx context.Context, dsn string) (database.Client, error) {
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to db: %v", err)
	}

	return &pgClient{
		db: &pg{pool: pool},
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
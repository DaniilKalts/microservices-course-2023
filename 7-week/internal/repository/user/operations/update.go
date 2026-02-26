package operations

import (
	"context"
	"time"

	sq "github.com/Masterminds/squirrel"

	"github.com/DaniilKalts/microservices-course-2023/7-week/internal/clients/database"
)

type UpdateInput struct {
	ID    string
	Name  *string
	Email *string
}

func Update(ctx context.Context, dbc database.Client, input UpdateInput) error {
	builderUpdate := sq.Update("users").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"id": input.ID})

	if input.Name != nil {
		builderUpdate = builderUpdate.Set("name", *input.Name)
	}
	if input.Email != nil {
		builderUpdate = builderUpdate.Set("email", *input.Email)
	}

	builderUpdate = builderUpdate.Set("updated_at", time.Now())

	query, args, err := builderUpdate.ToSql()
	if err != nil {
		return err
	}

	q := database.Query{Name: "user.Update", QueryRaw: query}
	if _, err = dbc.DB().ExecContext(ctx, q, args...); err != nil {
		return err
	}

	return nil
}

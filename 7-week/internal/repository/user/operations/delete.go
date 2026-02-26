package operations

import (
	"context"

	sq "github.com/Masterminds/squirrel"

	"github.com/DaniilKalts/microservices-course-2023/7-week/internal/clients/database"
)

type DeleteInput struct {
	ID string
}

func Delete(ctx context.Context, dbc database.Client, input DeleteInput) error {
	builderDelete := sq.Delete("users").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"id": input.ID})

	query, args, err := builderDelete.ToSql()
	if err != nil {
		return err
	}

	q := database.Query{Name: "user.Delete", QueryRaw: query}
	if _, err = dbc.DB().ExecContext(ctx, q, args...); err != nil {
		return err
	}

	return nil
}

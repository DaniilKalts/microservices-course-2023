package user

import (
	"context"

	sq "github.com/Masterminds/squirrel"

	"github.com/DaniilKalts/microservices-course-2023/5-week/internal/clients/database"
)

func (r *Repository) Delete(ctx context.Context, id string) error {
	builderDelete := sq.Delete("users").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"id": id})

	query, args, err := builderDelete.ToSql()
	if err != nil {
		return err
	}

	q := database.Query{Name: "user.Delete", QueryRaw: query}
	if _, err = r.dbc.DB().ExecContext(ctx, q, args...); err != nil {
		return err
	}

	return nil
}

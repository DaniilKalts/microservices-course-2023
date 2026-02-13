package user

import (
	"context"
	"time"

	sq "github.com/Masterminds/squirrel"

	"github.com/DaniilKalts/microservices-course-2023/4-week/internal/clients/database"
	domainUser "github.com/DaniilKalts/microservices-course-2023/4-week/internal/domain/user"
)

func (r *Repository) Update(ctx context.Context, id string, patch *domainUser.UpdatePatch) error {
	builderUpdate := sq.Update("users").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"id": id})

	if patch.Name != nil {
		builderUpdate = builderUpdate.Set("name", *patch.Name)
	}
	if patch.Email != nil {
		builderUpdate = builderUpdate.Set("email", *patch.Email)
	}

	builderUpdate = builderUpdate.Set("updated_at", time.Now())

	query, args, err := builderUpdate.ToSql()
	if err != nil {
		return err
	}

	q := database.Query{Name: "user.Update", QueryRaw: query}
	if _, err = r.dbc.DB().ExecContext(ctx, q, args...); err != nil {
		return err
	}

	return nil
}

package user

import (
	"context"

	sq "github.com/Masterminds/squirrel"

	"github.com/DaniilKalts/microservices-course-2023/4-week/internal/clients/database"
	domainUser "github.com/DaniilKalts/microservices-course-2023/4-week/internal/domain/user"
)

func (r *Repository) Get(ctx context.Context, id string) (*domainUser.Entity, error) {
	builderSelect := sq.Select("id", "name", "email", "role", "created_at", "updated_at").
		From("users").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"id": id}).
		Limit(1)

	query, args, err := builderSelect.ToSql()
	if err != nil {
		return nil, err
	}

	var user dbUser
	q := database.Query{Name: "user.Get", QueryRaw: query}
	if err = r.dbc.DB().ScanOneContext(ctx, &user, q, args...); err != nil {
		return nil, err
	}

	return toDomainFromDBUser(&user), nil
}

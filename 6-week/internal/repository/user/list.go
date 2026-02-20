package user

import (
	"context"

	sq "github.com/Masterminds/squirrel"

	"github.com/DaniilKalts/microservices-course-2023/6-week/internal/clients/database"
	domainUser "github.com/DaniilKalts/microservices-course-2023/6-week/internal/domain/user"
)

func (r *Repository) List(ctx context.Context) ([]domainUser.User, error) {
	builderSelect := sq.Select("id", "name", "email", "role", "created_at", "updated_at").
		From("users").
		PlaceholderFormat(sq.Dollar).
		OrderBy("created_at DESC")

	query, args, err := builderSelect.ToSql()
	if err != nil {
		return nil, err
	}

	var users []dbUser
	q := database.Query{Name: "user.List", QueryRaw: query}
	if err = r.dbc.DB().ScanAllContext(ctx, &users, q, args...); err != nil {
		return nil, err
	}

	entities := make([]domainUser.User, 0, len(users))
	for i := range users {
		entities = append(entities, *toDomainFromDBUser(&users[i]))
	}

	return entities, nil
}

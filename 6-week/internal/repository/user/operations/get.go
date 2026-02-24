package operations

import (
	"context"

	sq "github.com/Masterminds/squirrel"

	"github.com/DaniilKalts/microservices-course-2023/6-week/internal/clients/database"
	domainUser "github.com/DaniilKalts/microservices-course-2023/6-week/internal/domain/user"
	"github.com/DaniilKalts/microservices-course-2023/6-week/internal/repository/user/mapper"
	"github.com/DaniilKalts/microservices-course-2023/6-week/internal/repository/user/model"
)

type GetInput struct {
	ID string
}

func Get(ctx context.Context, dbc database.Client, input GetInput) (*domainUser.User, error) {
	builderSelect := sq.Select("id", "name", "email", "role", "created_at", "updated_at").
		From("users").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"id": input.ID}).
		Limit(1)

	query, args, err := builderSelect.ToSql()
	if err != nil {
		return nil, err
	}

	var user model.DBUser
	q := database.Query{Name: "user.Get", QueryRaw: query}
	if err = dbc.DB().ScanOneContext(ctx, &user, q, args...); err != nil {
		return nil, err
	}

	return mapper.ToDomainFromDBUser(&user), nil
}

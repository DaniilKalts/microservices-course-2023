package operations

import (
	"context"

	sq "github.com/Masterminds/squirrel"

	"github.com/DaniilKalts/microservices-course-2023/7-week/internal/clients/database"
	domainUser "github.com/DaniilKalts/microservices-course-2023/7-week/internal/domain/user"
	"github.com/DaniilKalts/microservices-course-2023/7-week/internal/repository/user/mapper"
	"github.com/DaniilKalts/microservices-course-2023/7-week/internal/repository/user/model"
)

func List(ctx context.Context, dbc database.Client) ([]domainUser.User, error) {
	builderSelect := sq.Select("id", "name", "email", "role", "created_at", "updated_at").
		From("users").
		PlaceholderFormat(sq.Dollar).
		OrderBy("created_at DESC")

	query, args, err := builderSelect.ToSql()
	if err != nil {
		return nil, err
	}

	var users []model.DBUser
	q := database.Query{Name: "user.List", QueryRaw: query}
	if err = dbc.DB().ScanAllContext(ctx, &users, q, args...); err != nil {
		return nil, err
	}

	return mapper.ToDomainFromDBUsers(users), nil
}

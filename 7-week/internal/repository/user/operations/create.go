package operations

import (
	"context"

	sq "github.com/Masterminds/squirrel"

	"github.com/DaniilKalts/microservices-course-2023/7-week/internal/clients/database"
	domainUser "github.com/DaniilKalts/microservices-course-2023/7-week/internal/domain/user"
	"github.com/DaniilKalts/microservices-course-2023/7-week/internal/repository/user/mapper"
)

type CreateInput struct {
	User         *domainUser.User
	PasswordHash string
}

func Create(ctx context.Context, dbc database.Client, input CreateInput) (string, error) {
	dbUser := mapper.ToDBUserFromDomain(input.User, input.PasswordHash)

	builderCreate := sq.Insert("users").
		PlaceholderFormat(sq.Dollar).
		Columns("id", "name", "email", "password_hash", "role").
		Values(dbUser.ID, dbUser.Name, dbUser.Email, dbUser.PasswordHash, dbUser.Role).
		Suffix("RETURNING id")

	query, args, err := builderCreate.ToSql()
	if err != nil {
		return "", err
	}

	var userID string
	if err = dbc.DB().ScanOneContext(ctx, &userID, database.Query{Name: "user.Create", QueryRaw: query}, args...); err != nil {
		return "", err
	}

	return userID, nil
}

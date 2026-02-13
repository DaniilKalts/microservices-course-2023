package user

import (
	"context"

	sq "github.com/Masterminds/squirrel"

	"github.com/DaniilKalts/microservices-course-2023/4-week/internal/clients/database"
	domainUser "github.com/DaniilKalts/microservices-course-2023/4-week/internal/domain/user"
)

func (r *Repository) Create(ctx context.Context, user *domainUser.Entity, passwordHash string) (string, error) {
	dbUser := toDBUserFromDomain(user)
	dbUser.PasswordHash = passwordHash

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
	if err = r.dbc.DB().ScanOneContext(ctx, &userID, database.Query{Name: "user.Create", QueryRaw: query}, args...); err != nil {
		return "", err
	}

	return userID, nil
}

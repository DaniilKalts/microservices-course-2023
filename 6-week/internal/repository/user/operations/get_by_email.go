package operations

import (
	"context"

	sq "github.com/Masterminds/squirrel"

	"github.com/DaniilKalts/microservices-course-2023/6-week/internal/clients/database"
	domainUser "github.com/DaniilKalts/microservices-course-2023/6-week/internal/domain/user"
	"github.com/DaniilKalts/microservices-course-2023/6-week/internal/repository/user/model"
)

type GetByEmailInput struct {
	Email string
}

type Credentials struct {
	ID           string
	PasswordHash string
	Role         domainUser.Role
}

func GetByEmail(
	ctx context.Context,
	dbc database.Client,
	input GetByEmailInput,
) (*Credentials, error) {
	builderSelect := sq.Select("id", "password_hash", "role").
		From("users").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"email": input.Email}).
		Limit(1)

	query, args, err := builderSelect.ToSql()
	if err != nil {
		return nil, err
	}

	var user model.DBUser
	q := database.Query{Name: "user.GetByEmail", QueryRaw: query}
	if err = dbc.DB().ScanOneContext(ctx, &user, q, args...); err != nil {
		return nil, err
	}

	return &Credentials{
		ID:           user.ID,
		PasswordHash: user.PasswordHash,
		Role:         domainUser.Role(user.Role),
	}, nil
}

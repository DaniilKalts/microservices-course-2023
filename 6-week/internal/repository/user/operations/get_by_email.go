package operations

import (
	"context"

	sq "github.com/Masterminds/squirrel"

	"github.com/DaniilKalts/microservices-course-2023/6-week/internal/clients/database"
	domainAuth "github.com/DaniilKalts/microservices-course-2023/6-week/internal/domain/auth"
	"github.com/DaniilKalts/microservices-course-2023/6-week/internal/repository/user/mapper"
	"github.com/DaniilKalts/microservices-course-2023/6-week/internal/repository/user/model"
)

type GetCredentialsByEmailInput struct {
	Email string
}

func GetCredentialsByEmail(
	ctx context.Context,
	dbc database.Client,
	input GetCredentialsByEmailInput,
) (*domainAuth.Credentials, error) {
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
	q := database.Query{Name: "user.GetCredentialsByEmail", QueryRaw: query}
	if err = dbc.DB().ScanOneContext(ctx, &user, q, args...); err != nil {
		return nil, err
	}

	return mapper.ToCredentialsFromDBUser(&user), nil
}

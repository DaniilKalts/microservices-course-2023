package user

import (
	"context"
	"errors"
	"time"

	sq "github.com/Masterminds/squirrel"

	"github.com/DaniilKalts/microservices-course-2023/8-week/internal/clients/database"
	domainUser "github.com/DaniilKalts/microservices-course-2023/8-week/internal/domain/user"
)

var sb = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

type Repository interface {
	Create(ctx context.Context, user *domainUser.User, passwordHash string) (string, error)
	List(ctx context.Context) ([]domainUser.User, error)
	GetByID(ctx context.Context, id string) (*domainUser.User, error)
	GetCredentialsByEmail(ctx context.Context, email string) (*domainUser.Credentials, error)
	Update(ctx context.Context, input domainUser.UpdateInput) error
	Delete(ctx context.Context, id string) error
}

type repository struct {
	client database.Client
}

func NewRepository(client database.Client) Repository {
	return &repository{
		client: client,
	}
}

func (repo *repository) Create(ctx context.Context, user *domainUser.User, passwordHash string) (string, error) {
	u := toDBUser(user, passwordHash)

	query, args, err := sb.Insert("users").
		Columns("id", "name", "email", "password_hash", "role").
		Values(u.ID, u.Name, u.Email, u.PasswordHash, u.Role).
		Suffix("RETURNING id").
		ToSql()
	if err != nil {
		return "", err
	}

	var userID string
	if err = repo.client.DB().ScanOneContext(ctx, &userID, database.Query{Name: "user.Create", QueryRaw: query}, args...); err != nil {
		if errors.Is(err, database.ErrUniqueViolation) {
			return "", domainUser.ErrEmailAlreadyExists
		}
		return "", err
	}

	return userID, nil
}

func (repo *repository) List(ctx context.Context) ([]domainUser.User, error) {
	query, args, err := sb.Select(userColumns...).
		From("users").
		OrderBy("created_at DESC").
		ToSql()
	if err != nil {
		return nil, err
	}

	var users []dbUser
	if err = repo.client.DB().ScanAllContext(ctx, &users, database.Query{Name: "user.List", QueryRaw: query}, args...); err != nil {
		return nil, err
	}

	return toDomainUsers(users), nil
}

func (repo *repository) GetByID(ctx context.Context, id string) (*domainUser.User, error) {
	query, args, err := sb.Select(userColumns...).
		From("users").
		Where(sq.Eq{"id": id}).
		Limit(1).
		ToSql()
	if err != nil {
		return nil, err
	}

	var user dbUser
	if err = repo.client.DB().ScanOneContext(ctx, &user, database.Query{Name: "user.GetByID", QueryRaw: query}, args...); err != nil {
		if errors.Is(err, database.ErrNotFound) {
			return nil, domainUser.ErrNotFound
		}
		return nil, err
	}

	return toDomainUser(&user), nil
}

func (repo *repository) GetCredentialsByEmail(ctx context.Context, email string) (*domainUser.Credentials, error) {
	query, args, err := sb.Select("id", "password_hash", "role").
		From("users").
		Where(sq.Eq{"email": email}).
		Limit(1).
		ToSql()
	if err != nil {
		return nil, err
	}

	var creds dbCredentials
	if err = repo.client.DB().ScanOneContext(ctx, &creds, database.Query{Name: "user.GetCredentialsByEmail", QueryRaw: query}, args...); err != nil {
		if errors.Is(err, database.ErrNotFound) {
			return nil, domainUser.ErrNotFound
		}
		return nil, err
	}

	return toDomainCredentials(&creds), nil
}

func (repo *repository) Update(ctx context.Context, input domainUser.UpdateInput) error {
	builder := sb.Update("users").
		Where(sq.Eq{"id": input.ID})

	hasFields := false

	if input.Name != nil {
		builder = builder.Set("name", *input.Name)
		hasFields = true
	}
	if input.Email != nil {
		builder = builder.Set("email", *input.Email)
		hasFields = true
	}
	if input.PasswordHash != nil {
		builder = builder.Set("password_hash", *input.PasswordHash)
		hasFields = true
	}

	if !hasFields {
		return domainUser.ErrNoFieldsToUpdate
	}

	builder = builder.Set("updated_at", time.Now())

	query, args, err := builder.ToSql()
	if err != nil {
		return err
	}

	result, err := repo.client.DB().ExecContext(ctx, database.Query{Name: "user.Update", QueryRaw: query}, args...)
	if err != nil {
		if errors.Is(err, database.ErrUniqueViolation) {
			return domainUser.ErrEmailAlreadyExists
		}
		return err
	}

	if result.RowsAffected() == 0 {
		return domainUser.ErrNotFound
	}

	return nil
}

func (repo *repository) Delete(ctx context.Context, id string) error {
	query, args, err := sb.Delete("users").
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return err
	}

	result, err := repo.client.DB().ExecContext(ctx, database.Query{Name: "user.Delete", QueryRaw: query}, args...)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return domainUser.ErrNotFound
	}

	return nil
}

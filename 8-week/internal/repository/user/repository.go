package user

import (
	"context"
	"errors"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"

	"github.com/DaniilKalts/microservices-course-2023/8-week/internal/clients/database"
	domainUser "github.com/DaniilKalts/microservices-course-2023/8-week/internal/domain/user"
	"github.com/DaniilKalts/microservices-course-2023/8-week/pkg/tracing"
)

var sb = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

type UpdateInput struct {
	ID           string
	Name         *string
	Email        *string
	PasswordHash *string
}

type Repository interface {
	Create(ctx context.Context, user *domainUser.User, passwordHash string) (string, error)
	List(ctx context.Context) ([]domainUser.User, error)
	GetByID(ctx context.Context, id string) (*domainUser.User, error)
	GetCredentialsByEmail(ctx context.Context, email string) (*domainUser.Credentials, error)
	Update(ctx context.Context, input UpdateInput) error
	Delete(ctx context.Context, id string) error
}

type repository struct {
	client database.Client
	logger *zap.Logger
}

func NewRepository(client database.Client, logger *zap.Logger) Repository {
	return &repository{
		client: client,
		logger: logger,
	}
}

func (repo *repository) Create(ctx context.Context, user *domainUser.User, passwordHash string) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "repository.user.Create")
	defer span.Finish()

	u := toDBUser(user, passwordHash)

	query, args, err := sb.Insert("users").
		Columns("id", "name", "email", "password_hash", "role").
		Values(u.ID, u.Name, u.Email, u.PasswordHash, u.Role).
		Suffix("RETURNING id").
		ToSql()
	if err != nil {
		tracing.LogError(repo.logger, span, "failed to build query", err)
		return "", err
	}

	var userID string
	if err = repo.client.DB().ScanOneContext(ctx, &userID, database.Query{Name: "user.Create", QueryRaw: query}, args...); err != nil {
		if errors.Is(err, database.ErrUniqueViolation) {
			tracing.LogWarn(repo.logger, span, "email already exists", err, zap.String("email", user.Email))
			return "", domainUser.ErrEmailAlreadyExists
		}
		tracing.LogError(repo.logger, span, "failed to create user", err)
		return "", err
	}

	return userID, nil
}

func (repo *repository) List(ctx context.Context) ([]domainUser.User, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "repository.user.List")
	defer span.Finish()

	query, args, err := sb.Select(userColumns...).
		From("users").
		OrderBy("created_at DESC").
		ToSql()
	if err != nil {
		tracing.LogError(repo.logger, span, "failed to build query", err)
		return nil, err
	}

	var users []dbUser
	if err = repo.client.DB().ScanAllContext(ctx, &users, database.Query{Name: "user.List", QueryRaw: query}, args...); err != nil {
		tracing.LogError(repo.logger, span, "failed to list users", err)
		return nil, err
	}

	return toDomainUsers(users), nil
}

func (repo *repository) GetByID(ctx context.Context, id string) (*domainUser.User, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "repository.user.GetByID")
	defer span.Finish()

	query, args, err := sb.Select(userColumns...).
		From("users").
		Where(sq.Eq{"id": id}).
		Limit(1).
		ToSql()
	if err != nil {
		tracing.LogError(repo.logger, span, "failed to build query", err)
		return nil, err
	}

	var user dbUser
	if err = repo.client.DB().ScanOneContext(ctx, &user, database.Query{Name: "user.GetByID", QueryRaw: query}, args...); err != nil {
		if errors.Is(err, database.ErrNotFound) {
			tracing.LogWarn(repo.logger, span, "user not found", err, zap.String("user_id", id))
			return nil, domainUser.ErrNotFound
		}
		tracing.LogError(repo.logger, span, "failed to get user", err, zap.String("user_id", id))
		return nil, err
	}

	return toDomainUser(&user), nil
}

func (repo *repository) GetCredentialsByEmail(ctx context.Context, email string) (*domainUser.Credentials, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "repository.user.GetCredentialsByEmail")
	defer span.Finish()

	query, args, err := sb.Select("id", "password_hash", "role").
		From("users").
		Where(sq.Eq{"email": email}).
		Limit(1).
		ToSql()
	if err != nil {
		tracing.LogError(repo.logger, span, "failed to build query", err)
		return nil, err
	}

	var creds dbCredentials
	if err = repo.client.DB().ScanOneContext(ctx, &creds, database.Query{Name: "user.GetCredentialsByEmail", QueryRaw: query}, args...); err != nil {
		if errors.Is(err, database.ErrNotFound) {
			tracing.LogWarn(repo.logger, span, "credentials not found", err, zap.String("email", email))
			return nil, domainUser.ErrNotFound
		}
		tracing.LogError(repo.logger, span, "failed to get credentials", err, zap.String("email", email))
		return nil, err
	}

	return toCredentials(&creds), nil
}

func (repo *repository) Update(ctx context.Context, input UpdateInput) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "repository.user.Update")
	defer span.Finish()

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
		tracing.LogWarn(repo.logger, span, "no fields to update", domainUser.ErrNoFieldsToUpdate, zap.String("user_id", input.ID))
		return domainUser.ErrNoFieldsToUpdate
	}

	builder = builder.Set("updated_at", time.Now())

	query, args, err := builder.ToSql()
	if err != nil {
		tracing.LogError(repo.logger, span, "failed to build query", err, zap.String("user_id", input.ID))
		return err
	}

	result, err := repo.client.DB().ExecContext(ctx, database.Query{Name: "user.Update", QueryRaw: query}, args...)
	if err != nil {
		if errors.Is(err, database.ErrUniqueViolation) {
			tracing.LogWarn(repo.logger, span, "email already exists", err, zap.String("user_id", input.ID))
			return domainUser.ErrEmailAlreadyExists
		}
		tracing.LogError(repo.logger, span, "failed to update user", err, zap.String("user_id", input.ID))
		return err
	}

	if result.RowsAffected() == 0 {
		tracing.LogWarn(repo.logger, span, "user not found", domainUser.ErrNotFound, zap.String("user_id", input.ID))
		return domainUser.ErrNotFound
	}

	return nil
}

func (repo *repository) Delete(ctx context.Context, id string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "repository.user.Delete")
	defer span.Finish()

	query, args, err := sb.Delete("users").
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		tracing.LogError(repo.logger, span, "failed to build query", err, zap.String("user_id", id))
		return err
	}

	result, err := repo.client.DB().ExecContext(ctx, database.Query{Name: "user.Delete", QueryRaw: query}, args...)
	if err != nil {
		tracing.LogError(repo.logger, span, "failed to delete user", err, zap.String("user_id", id))
		return err
	}

	if result.RowsAffected() == 0 {
		tracing.LogWarn(repo.logger, span, "user not found", domainUser.ErrNotFound, zap.String("user_id", id))
		return domainUser.ErrNotFound
	}

	return nil
}

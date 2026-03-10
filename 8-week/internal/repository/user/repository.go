package user

import (
	"context"
	"errors"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"go.uber.org/zap"

	"github.com/DaniilKalts/microservices-course-2023/8-week/internal/clients/database"
	domainUser "github.com/DaniilKalts/microservices-course-2023/8-week/internal/domain/user"
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
		repo.logger.Error("failed to build query", zap.Error(err))
		ext.Error.Set(span, true)
		span.LogKV("event", "error", "message", err.Error())
		return "", err
	}

	var userID string
	if err = repo.client.DB().ScanOneContext(ctx, &userID, database.Query{Name: "user.Create", QueryRaw: query}, args...); err != nil {
		if errors.Is(err, database.ErrUniqueViolation) {
			repo.logger.Warn("email already exists", zap.String("email", user.Email), zap.Error(err))
			span.LogKV("event", "warning", "message", err.Error())
			return "", domainUser.ErrEmailAlreadyExists
		}
		repo.logger.Error("failed to create user", zap.Error(err))
		ext.Error.Set(span, true)
		span.LogKV("event", "error", "message", err.Error())
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
		repo.logger.Error("failed to build query", zap.Error(err))
		ext.Error.Set(span, true)
		span.LogKV("event", "error", "message", err.Error())
		return nil, err
	}

	var users []dbUser
	if err = repo.client.DB().ScanAllContext(ctx, &users, database.Query{Name: "user.List", QueryRaw: query}, args...); err != nil {
		repo.logger.Error("failed to list users", zap.Error(err))
		ext.Error.Set(span, true)
		span.LogKV("event", "error", "message", err.Error())
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
		repo.logger.Error("failed to build query", zap.Error(err))
		ext.Error.Set(span, true)
		span.LogKV("event", "error", "message", err.Error())
		return nil, err
	}

	var user dbUser
	if err = repo.client.DB().ScanOneContext(ctx, &user, database.Query{Name: "user.GetByID", QueryRaw: query}, args...); err != nil {
		if errors.Is(err, database.ErrNotFound) {
			repo.logger.Warn("user not found", zap.String("user_id", id), zap.Error(err))
			span.LogKV("event", "warning", "message", err.Error())
			return nil, domainUser.ErrNotFound
		}
		repo.logger.Error("failed to get user", zap.String("user_id", id), zap.Error(err))
		ext.Error.Set(span, true)
		span.LogKV("event", "error", "message", err.Error())
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
		repo.logger.Error("failed to build query", zap.Error(err))
		ext.Error.Set(span, true)
		span.LogKV("event", "error", "message", err.Error())
		return nil, err
	}

	var creds dbCredentials
	if err = repo.client.DB().ScanOneContext(ctx, &creds, database.Query{Name: "user.GetCredentialsByEmail", QueryRaw: query}, args...); err != nil {
		if errors.Is(err, database.ErrNotFound) {
			repo.logger.Warn("credentials not found", zap.String("email", email), zap.Error(err))
			span.LogKV("event", "warning", "message", err.Error())
			return nil, domainUser.ErrNotFound
		}
		repo.logger.Error("failed to get credentials", zap.String("email", email), zap.Error(err))
		ext.Error.Set(span, true)
		span.LogKV("event", "error", "message", err.Error())
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
		repo.logger.Warn("no fields to update", zap.String("user_id", input.ID), zap.Error(domainUser.ErrNoFieldsToUpdate))
		span.LogKV("event", "warning", "message", domainUser.ErrNoFieldsToUpdate.Error())
		return domainUser.ErrNoFieldsToUpdate
	}

	builder = builder.Set("updated_at", time.Now())

	query, args, err := builder.ToSql()
	if err != nil {
		repo.logger.Error("failed to build query", zap.String("user_id", input.ID), zap.Error(err))
		ext.Error.Set(span, true)
		span.LogKV("event", "error", "message", err.Error())
		return err
	}

	result, err := repo.client.DB().ExecContext(ctx, database.Query{Name: "user.Update", QueryRaw: query}, args...)
	if err != nil {
		if errors.Is(err, database.ErrUniqueViolation) {
			repo.logger.Warn("email already exists", zap.String("user_id", input.ID), zap.Error(err))
			span.LogKV("event", "warning", "message", err.Error())
			return domainUser.ErrEmailAlreadyExists
		}
		repo.logger.Error("failed to update user", zap.String("user_id", input.ID), zap.Error(err))
		ext.Error.Set(span, true)
		span.LogKV("event", "error", "message", err.Error())
		return err
	}

	if result.RowsAffected() == 0 {
		repo.logger.Warn("user not found", zap.String("user_id", input.ID), zap.Error(domainUser.ErrNotFound))
		span.LogKV("event", "warning", "message", domainUser.ErrNotFound.Error())
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
		repo.logger.Error("failed to build query", zap.String("user_id", id), zap.Error(err))
		ext.Error.Set(span, true)
		span.LogKV("event", "error", "message", err.Error())
		return err
	}

	result, err := repo.client.DB().ExecContext(ctx, database.Query{Name: "user.Delete", QueryRaw: query}, args...)
	if err != nil {
		repo.logger.Error("failed to delete user", zap.String("user_id", id), zap.Error(err))
		ext.Error.Set(span, true)
		span.LogKV("event", "error", "message", err.Error())
		return err
	}

	if result.RowsAffected() == 0 {
		repo.logger.Warn("user not found", zap.String("user_id", id), zap.Error(domainUser.ErrNotFound))
		span.LogKV("event", "warning", "message", domainUser.ErrNotFound.Error())
		return domainUser.ErrNotFound
	}

	return nil
}

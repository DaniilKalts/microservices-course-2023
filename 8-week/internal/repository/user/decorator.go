package user

import (
	"context"
	"errors"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"go.uber.org/zap"

	domainUser "github.com/DaniilKalts/microservices-course-2023/8-week/internal/domain/user"
)

// --- Tracing ---

type tracingRepository struct {
	wrapped Repository
}

func WithTracing(repo Repository) Repository {
	return &tracingRepository{wrapped: repo}
}

func (r *tracingRepository) Create(ctx context.Context, user *domainUser.User, passwordHash string) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "repository.user.Create")
	defer span.Finish()

	id, err := r.wrapped.Create(ctx, user, passwordHash)
	if err != nil {
		ext.Error.Set(span, true)
		span.LogKV("event", "error", "message", err.Error())
	}

	return id, err
}

func (r *tracingRepository) List(ctx context.Context) ([]domainUser.User, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "repository.user.List")
	defer span.Finish()

	users, err := r.wrapped.List(ctx)
	if err != nil {
		ext.Error.Set(span, true)
		span.LogKV("event", "error", "message", err.Error())
	}

	return users, err
}

func (r *tracingRepository) GetByID(ctx context.Context, id string) (*domainUser.User, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "repository.user.GetByID")
	defer span.Finish()

	user, err := r.wrapped.GetByID(ctx, id)
	if err != nil {
		ext.Error.Set(span, true)
		span.LogKV("event", "error", "message", err.Error())
	}

	return user, err
}

func (r *tracingRepository) GetCredentialsByEmail(ctx context.Context, email string) (*domainUser.Credentials, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "repository.user.GetCredentialsByEmail")
	defer span.Finish()

	creds, err := r.wrapped.GetCredentialsByEmail(ctx, email)
	if err != nil {
		ext.Error.Set(span, true)
		span.LogKV("event", "error", "message", err.Error())
	}

	return creds, err
}

func (r *tracingRepository) Update(ctx context.Context, input domainUser.UpdateInput) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "repository.user.Update")
	defer span.Finish()

	err := r.wrapped.Update(ctx, input)
	if err != nil {
		ext.Error.Set(span, true)
		span.LogKV("event", "error", "message", err.Error())
	}

	return err
}

func (r *tracingRepository) Delete(ctx context.Context, id string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "repository.user.Delete")
	defer span.Finish()

	err := r.wrapped.Delete(ctx, id)
	if err != nil {
		ext.Error.Set(span, true)
		span.LogKV("event", "error", "message", err.Error())
	}

	return err
}

// --- Logging ---

type loggingRepository struct {
	wrapped Repository
	logger  *zap.Logger
}

func WithLogging(repo Repository, logger *zap.Logger) Repository {
	return &loggingRepository{wrapped: repo, logger: logger}
}

func (r *loggingRepository) Create(ctx context.Context, user *domainUser.User, passwordHash string) (string, error) {
	id, err := r.wrapped.Create(ctx, user, passwordHash)
	if err != nil {
		if errors.Is(err, domainUser.ErrEmailAlreadyExists) {
			r.logger.Warn("email already exists", zap.String("email", user.Email), zap.Error(err))
		} else {
			r.logger.Error("failed to create user", zap.Error(err))
		}
	}

	return id, err
}

func (r *loggingRepository) List(ctx context.Context) ([]domainUser.User, error) {
	users, err := r.wrapped.List(ctx)
	if err != nil {
		r.logger.Error("failed to list users", zap.Error(err))
	}

	return users, err
}

func (r *loggingRepository) GetByID(ctx context.Context, id string) (*domainUser.User, error) {
	user, err := r.wrapped.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, domainUser.ErrNotFound) {
			r.logger.Warn("user not found", zap.String("user_id", id), zap.Error(err))
		} else {
			r.logger.Error("failed to get user", zap.String("user_id", id), zap.Error(err))
		}
	}

	return user, err
}

func (r *loggingRepository) GetCredentialsByEmail(ctx context.Context, email string) (*domainUser.Credentials, error) {
	creds, err := r.wrapped.GetCredentialsByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, domainUser.ErrNotFound) {
			r.logger.Warn("credentials not found", zap.String("email", email), zap.Error(err))
		} else {
			r.logger.Error("failed to get credentials", zap.String("email", email), zap.Error(err))
		}
	}

	return creds, err
}

func (r *loggingRepository) Update(ctx context.Context, input domainUser.UpdateInput) error {
	err := r.wrapped.Update(ctx, input)
	if err != nil {
		switch {
		case errors.Is(err, domainUser.ErrNoFieldsToUpdate):
			r.logger.Warn("no fields to update", zap.String("user_id", input.ID), zap.Error(err))
		case errors.Is(err, domainUser.ErrEmailAlreadyExists):
			r.logger.Warn("email already exists", zap.String("user_id", input.ID), zap.Error(err))
		case errors.Is(err, domainUser.ErrNotFound):
			r.logger.Warn("user not found", zap.String("user_id", input.ID), zap.Error(err))
		default:
			r.logger.Error("failed to update user", zap.String("user_id", input.ID), zap.Error(err))
		}
	}

	return err
}

func (r *loggingRepository) Delete(ctx context.Context, id string) error {
	err := r.wrapped.Delete(ctx, id)
	if err != nil {
		if errors.Is(err, domainUser.ErrNotFound) {
			r.logger.Warn("user not found", zap.String("user_id", id), zap.Error(err))
		} else {
			r.logger.Error("failed to delete user", zap.String("user_id", id), zap.Error(err))
		}
	}

	return err
}

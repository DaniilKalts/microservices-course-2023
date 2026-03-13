package user

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"

	domainUser "github.com/DaniilKalts/microservices-course-2023/8-week/internal/domain/user"
)

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

package user

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"

	domainUser "github.com/DaniilKalts/microservices-course-2023/8-week/internal/domain/user"
)

type tracingService struct {
	wrapped Service
}

func WithTracing(svc Service) Service {
	return &tracingService{wrapped: svc}
}

func (s *tracingService) Create(ctx context.Context, input domainUser.CreateInput) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "service.user.Create")
	defer span.Finish()

	id, err := s.wrapped.Create(ctx, input)
	if err != nil {
		ext.Error.Set(span, true)
		span.LogKV("event", "error", "message", err.Error())
	}

	return id, err
}

func (s *tracingService) List(ctx context.Context) ([]domainUser.User, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "service.user.List")
	defer span.Finish()

	users, err := s.wrapped.List(ctx)
	if err != nil {
		ext.Error.Set(span, true)
		span.LogKV("event", "error", "message", err.Error())
	}

	return users, err
}

func (s *tracingService) Get(ctx context.Context, id string) (*domainUser.User, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "service.user.Get")
	defer span.Finish()

	user, err := s.wrapped.Get(ctx, id)
	if err != nil {
		ext.Error.Set(span, true)
		span.LogKV("event", "error", "message", err.Error())
	}

	return user, err
}

func (s *tracingService) GetCredentialsByEmail(ctx context.Context, email string) (*domainUser.Credentials, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "service.user.GetCredentialsByEmail")
	defer span.Finish()

	creds, err := s.wrapped.GetCredentialsByEmail(ctx, email)
	if err != nil {
		ext.Error.Set(span, true)
		span.LogKV("event", "error", "message", err.Error())
	}

	return creds, err
}

func (s *tracingService) Update(ctx context.Context, input domainUser.UpdateInput) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "service.user.Update")
	defer span.Finish()

	err := s.wrapped.Update(ctx, input)
	if err != nil {
		ext.Error.Set(span, true)
		span.LogKV("event", "error", "message", err.Error())
	}

	return err
}

func (s *tracingService) Delete(ctx context.Context, id string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "service.user.Delete")
	defer span.Finish()

	err := s.wrapped.Delete(ctx, id)
	if err != nil {
		ext.Error.Set(span, true)
		span.LogKV("event", "error", "message", err.Error())
	}

	return err
}

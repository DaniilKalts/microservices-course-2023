package auth

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

func (s *tracingService) Register(ctx context.Context, input RegisterInput) (domainUser.User, TokenPair, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "service.auth.Register")
	defer span.Finish()

	user, tokens, err := s.wrapped.Register(ctx, input)
	if err != nil {
		ext.Error.Set(span, true)
		span.LogKV("event", "error", "message", err.Error())
	}

	return user, tokens, err
}

func (s *tracingService) Login(ctx context.Context, input LoginInput) (TokenPair, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "service.auth.Login")
	defer span.Finish()

	tokens, err := s.wrapped.Login(ctx, input)
	if err != nil {
		ext.Error.Set(span, true)
		span.LogKV("event", "error", "message", err.Error())
	}

	return tokens, err
}

func (s *tracingService) Logout(ctx context.Context, input LogoutInput) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "service.auth.Logout")
	defer span.Finish()

	err := s.wrapped.Logout(ctx, input)
	if err != nil {
		ext.Error.Set(span, true)
		span.LogKV("event", "error", "message", err.Error())
	}

	return err
}

func (s *tracingService) Refresh(ctx context.Context, input RefreshInput) (TokenPair, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "service.auth.Refresh")
	defer span.Finish()

	tokens, err := s.wrapped.Refresh(ctx, input)
	if err != nil {
		ext.Error.Set(span, true)
		span.LogKV("event", "error", "message", err.Error())
	}

	return tokens, err
}

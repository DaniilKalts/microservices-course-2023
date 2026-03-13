package auth

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"go.uber.org/zap"

	domainUser "github.com/DaniilKalts/microservices-course-2023/8-week/internal/domain/user"
)

// --- Tracing ---

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

// --- Logging ---

type loggingService struct {
	wrapped Service
	logger  *zap.Logger
}

func WithLogging(svc Service, logger *zap.Logger) Service {
	return &loggingService{wrapped: svc, logger: logger}
}

func (s *loggingService) Register(ctx context.Context, input RegisterInput) (domainUser.User, TokenPair, error) {
	user, tokens, err := s.wrapped.Register(ctx, input)
	if err != nil {
		s.logger.Error("failed to register user", zap.String("email", input.Email), zap.Error(err))
		return user, tokens, err
	}

	s.logger.Info("user registered", zap.String("user_id", user.ID), zap.String("email", input.Email))

	return user, tokens, nil
}

func (s *loggingService) Login(ctx context.Context, input LoginInput) (TokenPair, error) {
	tokens, err := s.wrapped.Login(ctx, input)
	if err != nil {
		s.logger.Warn("login failed", zap.String("email", input.Email), zap.Error(err))
		return tokens, err
	}

	s.logger.Info("user logged in", zap.String("email", input.Email))

	return tokens, nil
}

func (s *loggingService) Logout(ctx context.Context, input LogoutInput) error {
	err := s.wrapped.Logout(ctx, input)
	if err != nil {
		s.logger.Warn("logout failed", zap.Error(err))
		return err
	}

	s.logger.Info("user logged out")

	return nil
}

func (s *loggingService) Refresh(ctx context.Context, input RefreshInput) (TokenPair, error) {
	tokens, err := s.wrapped.Refresh(ctx, input)
	if err != nil {
		s.logger.Warn("token refresh failed", zap.Error(err))
		return tokens, err
	}

	s.logger.Info("tokens refreshed")

	return tokens, nil
}

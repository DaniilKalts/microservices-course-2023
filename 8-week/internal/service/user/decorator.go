package user

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

// --- Logging ---

type loggingService struct {
	wrapped Service
	logger  *zap.Logger
}

func WithLogging(svc Service, logger *zap.Logger) Service {
	return &loggingService{wrapped: svc, logger: logger}
}

func (s *loggingService) Create(ctx context.Context, input domainUser.CreateInput) (string, error) {
	id, err := s.wrapped.Create(ctx, input)
	if err != nil {
		s.logger.Error("failed to create user", zap.Error(err))
		return id, err
	}

	s.logger.Info("user created", zap.String("user_id", id))

	return id, nil
}

func (s *loggingService) List(ctx context.Context) ([]domainUser.User, error) {
	return s.wrapped.List(ctx)
}

func (s *loggingService) Get(ctx context.Context, id string) (*domainUser.User, error) {
	return s.wrapped.Get(ctx, id)
}

func (s *loggingService) GetCredentialsByEmail(ctx context.Context, email string) (*domainUser.Credentials, error) {
	return s.wrapped.GetCredentialsByEmail(ctx, email)
}

func (s *loggingService) Update(ctx context.Context, input domainUser.UpdateInput) error {
	err := s.wrapped.Update(ctx, input)
	if err != nil {
		s.logger.Error("failed to update user", zap.String("user_id", input.ID), zap.Error(err))
		return err
	}

	s.logger.Info("user updated", zap.String("user_id", input.ID))

	return nil
}

func (s *loggingService) Delete(ctx context.Context, id string) error {
	err := s.wrapped.Delete(ctx, id)
	if err != nil {
		s.logger.Error("failed to delete user", zap.String("user_id", id), zap.Error(err))
		return err
	}

	s.logger.Info("user deleted", zap.String("user_id", id))

	return nil
}

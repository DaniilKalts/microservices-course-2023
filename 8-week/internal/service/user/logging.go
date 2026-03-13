package user

import (
	"context"

	"go.uber.org/zap"

	domainUser "github.com/DaniilKalts/microservices-course-2023/8-week/internal/domain/user"
)

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

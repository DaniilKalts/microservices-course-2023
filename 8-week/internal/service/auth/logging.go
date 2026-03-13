package auth

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

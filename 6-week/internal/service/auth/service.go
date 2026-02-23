package auth

import (
	"context"

	domainAuth "github.com/DaniilKalts/microservices-course-2023/6-week/internal/domain/auth"
	"github.com/DaniilKalts/microservices-course-2023/6-week/internal/service/auth/operations"
	userService "github.com/DaniilKalts/microservices-course-2023/6-week/internal/service/user"
	"github.com/DaniilKalts/microservices-course-2023/6-week/pkg/jwt"
)

type Service interface {
	Register(ctx context.Context, name, email, password, passwordConfirm string) (string, domainAuth.TokenPair, error)
	Login(ctx context.Context, email, password string) (domainAuth.TokenPair, error)
	Logout(ctx context.Context, refreshToken string) error
	Refresh(ctx context.Context, refreshToken string) (domainAuth.TokenPair, error)
}

type service struct {
	userService userService.Service
	jwtManager  jwt.Manager
}

func NewService(userSvc userService.Service, jwtManager jwt.Manager) Service {
	return &service{
		userService: userSvc,
		jwtManager:  jwtManager,
	}
}

func (s *service) Register(ctx context.Context, name, email, password, passwordConfirm string) (string, domainAuth.TokenPair, error) {
	return operations.Register(ctx, s.userService, s.jwtManager, name, email, password, passwordConfirm)
}

func (s *service) Login(ctx context.Context, email, password string) (domainAuth.TokenPair, error) {
	return operations.Login(ctx, email, password)
}

func (s *service) Logout(ctx context.Context, refreshToken string) error {
	return operations.Logout(ctx, s.jwtManager, refreshToken)
}

func (s *service) Refresh(ctx context.Context, refreshToken string) (domainAuth.TokenPair, error) {
	return operations.Refresh(ctx, s.jwtManager, refreshToken)
}

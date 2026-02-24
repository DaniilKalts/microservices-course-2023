package auth

import (
	"context"

	domainAuth "github.com/DaniilKalts/microservices-course-2023/6-week/internal/domain/auth"
	"github.com/DaniilKalts/microservices-course-2023/6-week/internal/service/auth/operations"
	userService "github.com/DaniilKalts/microservices-course-2023/6-week/internal/service/user"
	"github.com/DaniilKalts/microservices-course-2023/6-week/pkg/jwt"
)

type Service interface {
	Register(ctx context.Context, input operations.RegisterInput) (string, domainAuth.TokenPair, error)
	Login(ctx context.Context, input operations.LoginInput) (domainAuth.TokenPair, error)
	Logout(ctx context.Context, input operations.LogoutInput) error
	Refresh(ctx context.Context, input operations.RefreshInput) (domainAuth.TokenPair, error)
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

func (s *service) Register(ctx context.Context, input operations.RegisterInput) (string, domainAuth.TokenPair, error) {
	return operations.Register(ctx, s.userService, s.jwtManager, input)
}

func (s *service) Login(ctx context.Context, input operations.LoginInput) (domainAuth.TokenPair, error) {
	return operations.Login(ctx, input)
}

func (s *service) Logout(ctx context.Context, input operations.LogoutInput) error {
	return operations.Logout(ctx, s.jwtManager, input)
}

func (s *service) Refresh(ctx context.Context, input operations.RefreshInput) (domainAuth.TokenPair, error) {
	return operations.Refresh(ctx, s.jwtManager, input)
}

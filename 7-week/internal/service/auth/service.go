package auth

import (
	"context"

	domainAuth "github.com/DaniilKalts/microservices-course-2023/7-week/internal/domain/auth"
	domainUser "github.com/DaniilKalts/microservices-course-2023/7-week/internal/domain/user"
	"github.com/DaniilKalts/microservices-course-2023/7-week/internal/repository"
	"github.com/DaniilKalts/microservices-course-2023/7-week/internal/service/auth/operations"
	userService "github.com/DaniilKalts/microservices-course-2023/7-week/internal/service/user"
	"github.com/DaniilKalts/microservices-course-2023/7-week/pkg/jwt"
)

type Service interface {
	Register(ctx context.Context, input operations.RegisterInput) (domainUser.User, domainAuth.TokenPair, error)
	Login(ctx context.Context, input operations.LoginInput) (domainAuth.TokenPair, error)
	Logout(ctx context.Context, input operations.LogoutInput) error
	Refresh(ctx context.Context, input operations.RefreshInput) (domainAuth.TokenPair, error)
}

type service struct {
	userService userService.Service
	userRepo    repository.UserRepository
	jwtManager  jwt.Manager
}

func NewService(userSvc userService.Service, userRepo repository.UserRepository, jwtManager jwt.Manager) Service {
	return &service{
		userService: userSvc,
		userRepo:    userRepo,
		jwtManager:  jwtManager,
	}
}

func (s *service) Register(ctx context.Context, input operations.RegisterInput) (domainUser.User, domainAuth.TokenPair, error) {
	return operations.Register(ctx, s.userService, s.jwtManager, input)
}

func (s *service) Login(ctx context.Context, input operations.LoginInput) (domainAuth.TokenPair, error) {
	return operations.Login(ctx, s.userRepo, s.jwtManager, input)
}

func (s *service) Logout(ctx context.Context, input operations.LogoutInput) error {
	return operations.Logout(ctx, s.jwtManager, input)
}

func (s *service) Refresh(ctx context.Context, input operations.RefreshInput) (domainAuth.TokenPair, error) {
	return operations.Refresh(ctx, s.jwtManager, input)
}

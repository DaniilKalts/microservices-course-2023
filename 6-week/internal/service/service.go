package service

import (
	"context"

	domainAuth "github.com/DaniilKalts/microservices-course-2023/6-week/internal/domain/auth"
	domainUser "github.com/DaniilKalts/microservices-course-2023/6-week/internal/domain/user"
)

type UserService interface {
	Create(ctx context.Context, user *domainUser.User, password, passwordConfirm string) (string, error)
	List(ctx context.Context) ([]domainUser.User, error)
	Get(ctx context.Context, id string) (*domainUser.User, error)
	Update(ctx context.Context, id string, patch *domainUser.UpdatePatch) error
	Delete(ctx context.Context, id string) error
}

type AuthService interface {
	Register(ctx context.Context, name, email, password string) (string, domainAuth.TokenPair, error)
	Login(ctx context.Context, email, password string) (domainAuth.TokenPair, error)
	Logout(ctx context.Context, refreshToken string) error
	Refresh(ctx context.Context, refreshToken string) (domainAuth.TokenPair, error)
}

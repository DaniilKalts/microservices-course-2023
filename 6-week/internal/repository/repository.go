package repository

import (
	"context"

	domainUser "github.com/DaniilKalts/microservices-course-2023/6-week/internal/domain/user"
)

//go:generate minimock -i UserRepository -o ./mocks -s _mock.go
type UserRepository interface {
	Create(ctx context.Context, user *domainUser.User, passwordHash string) (string, error)
	List(ctx context.Context) ([]domainUser.User, error)
	Get(ctx context.Context, id string) (*domainUser.User, error)
	Update(ctx context.Context, id string, patch *domainUser.UpdatePatch) error
	Delete(ctx context.Context, id string) error
}

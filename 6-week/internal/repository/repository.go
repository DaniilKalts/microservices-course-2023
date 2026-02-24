package repository

import (
	"context"

	domainUser "github.com/DaniilKalts/microservices-course-2023/6-week/internal/domain/user"
)

type UserUpdateInput struct {
	ID    string
	Name  *string
	Email *string
}

//go:generate minimock -i UserRepository -o ./mocks -s _mock.go
type UserRepository interface {
	Create(ctx context.Context, user *domainUser.User, passwordHash string) (string, error)
	List(ctx context.Context) ([]domainUser.User, error)
	Get(ctx context.Context, id string) (*domainUser.User, error)
	Update(ctx context.Context, input UserUpdateInput) error
	Delete(ctx context.Context, id string) error
}

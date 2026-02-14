package repository

import (
	"context"

	domainUser "github.com/DaniilKalts/microservices-course-2023/5-week/internal/domain/user"
)

//go:generate minimock -i UserRepository -o ./mocks -s _mock.go

type UserRepository interface {
	Create(ctx context.Context, user *domainUser.Entity, passwordHash string) (string, error)
	Get(ctx context.Context, id string) (*domainUser.Entity, error)
	Update(ctx context.Context, id string, patch *domainUser.UpdatePatch) error
	Delete(ctx context.Context, id string) error
}

package service

import (
	"context"

	domainUser "github.com/DaniilKalts/microservices-course-2023/4-week/internal/domain/user"
)

type UserService interface {
	Create(ctx context.Context, user *domainUser.Entity, password, passwordConfirm string) (string, error)
	Get(ctx context.Context, id string) (*domainUser.Entity, error)
	Update(ctx context.Context, id string, patch *domainUser.UpdatePatch) error
	Delete(ctx context.Context, id string) error
}

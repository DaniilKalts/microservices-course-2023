package user

import (
	"context"

	domainUser "github.com/DaniilKalts/microservices-course-2023/6-week/internal/domain/user"
	"github.com/DaniilKalts/microservices-course-2023/6-week/internal/repository"
	"github.com/DaniilKalts/microservices-course-2023/6-week/internal/service/user/operations"
)

type Service interface {
	Create(ctx context.Context, user *domainUser.User, password, passwordConfirm string) (string, error)
	List(ctx context.Context) ([]domainUser.User, error)
	Get(ctx context.Context, id string) (*domainUser.User, error)
	Update(ctx context.Context, id string, patch *domainUser.UpdatePatch) error
	Delete(ctx context.Context, id string) error
}

type service struct {
	repo repository.UserRepository
}

func NewService(repo repository.UserRepository) Service {
	return &service{repo: repo}
}

func (s *service) Create(ctx context.Context, user *domainUser.User, password, passwordConfirm string) (string, error) {
	return operations.Create(ctx, s.repo, user, password, passwordConfirm)
}

func (s *service) List(ctx context.Context) ([]domainUser.User, error) {
	return operations.List(ctx, s.repo)
}

func (s *service) Get(ctx context.Context, id string) (*domainUser.User, error) {
	return operations.Get(ctx, s.repo, id)
}

func (s *service) Update(ctx context.Context, id string, patch *domainUser.UpdatePatch) error {
	return operations.Update(ctx, s.repo, id, patch)
}

func (s *service) Delete(ctx context.Context, id string) error {
	return operations.Delete(ctx, s.repo, id)
}

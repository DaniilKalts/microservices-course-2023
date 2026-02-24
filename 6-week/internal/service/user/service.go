package user

import (
	"context"

	domainUser "github.com/DaniilKalts/microservices-course-2023/6-week/internal/domain/user"
	"github.com/DaniilKalts/microservices-course-2023/6-week/internal/repository"
	"github.com/DaniilKalts/microservices-course-2023/6-week/internal/service/user/operations"
)

type Service interface {
	Create(ctx context.Context, input operations.CreateInput) (string, error)
	List(ctx context.Context) ([]domainUser.User, error)
	Get(ctx context.Context, input operations.GetInput) (*domainUser.User, error)
	Update(ctx context.Context, input operations.UpdateInput) error
	Delete(ctx context.Context, input operations.DeleteInput) error
}

type service struct {
	repo repository.UserRepository
}

func NewService(repo repository.UserRepository) Service {
	return &service{repo: repo}
}

func (s *service) Create(ctx context.Context, input operations.CreateInput) (string, error) {
	return operations.Create(ctx, s.repo, input)
}

func (s *service) List(ctx context.Context) ([]domainUser.User, error) {
	return operations.List(ctx, s.repo)
}

func (s *service) Get(ctx context.Context, input operations.GetInput) (*domainUser.User, error) {
	return operations.Get(ctx, s.repo, input)
}

func (s *service) Update(ctx context.Context, input operations.UpdateInput) error {
	return operations.Update(ctx, s.repo, input)
}

func (s *service) Delete(ctx context.Context, input operations.DeleteInput) error {
	return operations.Delete(ctx, s.repo, input)
}

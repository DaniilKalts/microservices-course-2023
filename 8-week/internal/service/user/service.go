package user

import (
	"context"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	domainUser "github.com/DaniilKalts/microservices-course-2023/8-week/internal/domain/user"
)

type Repository interface {
	Create(ctx context.Context, user *domainUser.User, passwordHash string) (string, error)
	List(ctx context.Context) ([]domainUser.User, error)
	GetByID(ctx context.Context, id string) (*domainUser.User, error)
	GetCredentialsByEmail(ctx context.Context, email string) (*domainUser.Credentials, error)
	Update(ctx context.Context, input domainUser.UpdateInput) error
	Delete(ctx context.Context, id string) error
}

//go:generate minimock -i github.com/DaniilKalts/microservices-course-2023/8-week/internal/service/user.Service -o mock.go -n UserServiceMock -p user

type Service interface {
	Create(ctx context.Context, input domainUser.CreateInput) (string, error)
	List(ctx context.Context) ([]domainUser.User, error)
	Get(ctx context.Context, id string) (*domainUser.User, error)
	GetCredentialsByEmail(ctx context.Context, email string) (*domainUser.Credentials, error)
	Update(ctx context.Context, input domainUser.UpdateInput) error
	Delete(ctx context.Context, id string) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{
		repo: repo,
	}
}

func (s *service) Create(ctx context.Context, input domainUser.CreateInput) (string, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return "", err
	}
	input.User.ID = id.String()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	userID, err := s.repo.Create(ctx, input.User, string(hashedPassword))
	if err != nil {
		return "", err
	}

	return userID, nil
}

func (s *service) List(ctx context.Context) ([]domainUser.User, error) {
	users, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (s *service) Get(ctx context.Context, id string) (*domainUser.User, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *service) GetCredentialsByEmail(ctx context.Context, email string) (*domainUser.Credentials, error) {
	creds, err := s.repo.GetCredentialsByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	return creds, nil
}

func (s *service) Update(ctx context.Context, input domainUser.UpdateInput) error {
	if input.Password != nil {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*input.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}

		passwordHash := string(hashedPassword)
		input.PasswordHash = &passwordHash
	}

	if err := s.repo.Update(ctx, input); err != nil {
		return err
	}

	return nil
}

func (s *service) Delete(ctx context.Context, id string) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}

	return nil
}

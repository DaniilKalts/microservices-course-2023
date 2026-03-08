package user

import (
	"context"
	"errors"

	"fmt"

	"github.com/google/uuid"
	"github.com/opentracing/opentracing-go"
	"golang.org/x/crypto/bcrypt"

	domainUser "github.com/DaniilKalts/microservices-course-2023/7-week/internal/domain/user"
	userRepository "github.com/DaniilKalts/microservices-course-2023/7-week/internal/repository/user"
)

var (
	ErrPasswordMismatch   = errors.New("passwords don't match")
	ErrNotFound           = errors.New("user not found")
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrNoFieldsToUpdate   = errors.New("no fields to update")
)

type CreateInput struct {
	User            *domainUser.User
	Password        string
	PasswordConfirm string
}

type UpdateInput struct {
	ID       string
	Name     *string
	Email    *string
	Password *string
}

type Service interface {
	Create(ctx context.Context, input CreateInput) (string, error)
	List(ctx context.Context) ([]domainUser.User, error)
	Get(ctx context.Context, id string) (*domainUser.User, error)
	GetCredentialsByEmail(ctx context.Context, email string) (*domainUser.Credentials, error)
	Update(ctx context.Context, input UpdateInput) error
	Delete(ctx context.Context, id string) error
}

type service struct {
	repo userRepository.Repository
}

func NewService(repo userRepository.Repository) Service {
	return &service{
		repo: repo,
	}
}

func (s *service) Create(ctx context.Context, input CreateInput) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "service.user.Create")
	defer span.Finish()

	if input.Password != input.PasswordConfirm {
		return "", ErrPasswordMismatch
	}

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
		return "", mapRepoError(err)
	}

	return userID, nil
}

func (s *service) List(ctx context.Context) ([]domainUser.User, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "service.user.List")
	defer span.Finish()

	users, err := s.repo.List(ctx)
	if err != nil {
		return nil, mapRepoError(err)
	}

	return users, nil
}

func (s *service) Get(ctx context.Context, id string) (*domainUser.User, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "service.user.Get")
	defer span.Finish()

	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, mapRepoError(err)
	}

	return user, nil
}

func (s *service) GetCredentialsByEmail(ctx context.Context, email string) (*domainUser.Credentials, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "service.user.GetCredentialsByEmail")
	defer span.Finish()

	creds, err := s.repo.GetCredentialsByEmail(ctx, email)
	if err != nil {
		return nil, mapRepoError(err)
	}

	return creds, nil
}

func (s *service) Update(ctx context.Context, input UpdateInput) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "service.user.Update")
	defer span.Finish()

	repoInput := userRepository.UpdateInput{
		ID:    input.ID,
		Name:  input.Name,
		Email: input.Email,
	}

	if input.Password != nil {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*input.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}

		passwordHash := string(hashedPassword)
		repoInput.PasswordHash = &passwordHash
	}

	if err := s.repo.Update(ctx, repoInput); err != nil {
		return mapRepoError(err)
	}

	return nil
}

func (s *service) Delete(ctx context.Context, id string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "service.user.Delete")
	defer span.Finish()

	if err := s.repo.Delete(ctx, id); err != nil {
		return mapRepoError(err)
	}

	return nil
}

func mapRepoError(err error) error {
	switch {
	case errors.Is(err, userRepository.ErrNotFound):
		return fmt.Errorf("%w: %v", ErrNotFound, err)
	case errors.Is(err, userRepository.ErrEmailAlreadyExists):
		return fmt.Errorf("%w: %v", ErrEmailAlreadyExists, err)
	case errors.Is(err, userRepository.ErrNoFieldsToUpdate):
		return fmt.Errorf("%w: %v", ErrNoFieldsToUpdate, err)
	default:
		return err
	}
}

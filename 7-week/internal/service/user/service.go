package user

import (
	"context"

	"github.com/google/uuid"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"

	domainUser "github.com/DaniilKalts/microservices-course-2023/7-week/internal/domain/user"
	userRepository "github.com/DaniilKalts/microservices-course-2023/7-week/internal/repository/user"
	"github.com/DaniilKalts/microservices-course-2023/7-week/pkg/tracing"
)

type Service interface {
	Create(ctx context.Context, input CreateInput) (string, error)
	List(ctx context.Context) ([]domainUser.User, error)
	Get(ctx context.Context, id string) (*domainUser.User, error)
	GetCredentialsByEmail(ctx context.Context, email string) (*domainUser.Credentials, error)
	Update(ctx context.Context, input UpdateInput) error
	Delete(ctx context.Context, id string) error
}

type service struct {
	repo   userRepository.Repository
	logger *zap.Logger
}

func NewService(repo userRepository.Repository, logger *zap.Logger) Service {
	return &service{
		repo:   repo,
		logger: logger,
	}
}

func (s *service) Create(ctx context.Context, input CreateInput) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "service.user.Create")
	defer span.Finish()

	id, err := uuid.NewV7()
	if err != nil {
		tracing.LogError(s.logger, span, "failed to generate uuid", err)
		return "", err
	}
	input.User.ID = id.String()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		tracing.LogError(s.logger, span, "failed to hash password", err)
		return "", err
	}

	userID, err := s.repo.Create(ctx, input.User, string(hashedPassword))
	if err != nil {
		return "", err
	}

	s.logger.Info("user created", zap.String("user_id", userID))

	return userID, nil
}

func (s *service) List(ctx context.Context) ([]domainUser.User, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "service.user.List")
	defer span.Finish()

	users, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (s *service) Get(ctx context.Context, id string) (*domainUser.User, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "service.user.Get")
	defer span.Finish()

	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *service) GetCredentialsByEmail(ctx context.Context, email string) (*domainUser.Credentials, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "service.user.GetCredentialsByEmail")
	defer span.Finish()

	creds, err := s.repo.GetCredentialsByEmail(ctx, email)
	if err != nil {
		return nil, err
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
			tracing.LogError(s.logger, span, "failed to hash password", err, zap.String("user_id", input.ID))
			return err
		}

		passwordHash := string(hashedPassword)
		repoInput.PasswordHash = &passwordHash
	}

	if err := s.repo.Update(ctx, repoInput); err != nil {
		return err
	}

	s.logger.Info("user updated", zap.String("user_id", input.ID))

	return nil
}

func (s *service) Delete(ctx context.Context, id string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "service.user.Delete")
	defer span.Finish()

	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}

	s.logger.Info("user deleted", zap.String("user_id", id))

	return nil
}

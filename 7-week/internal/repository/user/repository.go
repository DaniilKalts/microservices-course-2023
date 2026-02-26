package user

import (
	"context"

	"github.com/DaniilKalts/microservices-course-2023/7-week/internal/clients/database"
	domainAuth "github.com/DaniilKalts/microservices-course-2023/7-week/internal/domain/auth"
	domainUser "github.com/DaniilKalts/microservices-course-2023/7-week/internal/domain/user"
	"github.com/DaniilKalts/microservices-course-2023/7-week/internal/repository/user/operations"
)

type Repository interface {
	Create(ctx context.Context, user *domainUser.User, passwordHash string) (string, error)
	List(ctx context.Context) ([]domainUser.User, error)
	GetByID(ctx context.Context, id string) (*domainUser.User, error)
	GetCredentialsByEmail(ctx context.Context, email string) (*domainAuth.Credentials, error)
	Update(ctx context.Context, input operations.UpdateInput) error
	Delete(ctx context.Context, id string) error
}

type repository struct {
	dbc database.Client
}

func NewRepository(dbc database.Client) Repository {
	return &repository{dbc: dbc}
}

func (r *repository) Create(ctx context.Context, user *domainUser.User, passwordHash string) (string, error) {
	return operations.Create(ctx, r.dbc, operations.CreateInput{
		User:         user,
		PasswordHash: passwordHash,
	})
}

func (r *repository) List(ctx context.Context) ([]domainUser.User, error) {
	return operations.List(ctx, r.dbc)
}

func (r *repository) GetByID(ctx context.Context, id string) (*domainUser.User, error) {
	return operations.GetByID(ctx, r.dbc, operations.GetByIDInput{ID: id})
}

func (r *repository) GetCredentialsByEmail(ctx context.Context, email string) (*domainAuth.Credentials, error) {
	return operations.GetCredentialsByEmail(ctx, r.dbc, operations.GetCredentialsByEmailInput{Email: email})
}

func (r *repository) Update(ctx context.Context, input operations.UpdateInput) error {
	return operations.Update(ctx, r.dbc, input)
}

func (r *repository) Delete(ctx context.Context, id string) error {
	return operations.Delete(ctx, r.dbc, operations.DeleteInput{ID: id})
}

package user

import (
	"context"

	"github.com/DaniilKalts/microservices-course-2023/6-week/internal/clients/database"
	domainUser "github.com/DaniilKalts/microservices-course-2023/6-week/internal/domain/user"
	"github.com/DaniilKalts/microservices-course-2023/6-week/internal/repository/user/operations"
)

type UpdateInput struct {
	ID    string
	Name  *string
	Email *string
}

type Credentials struct {
	ID           string
	PasswordHash string
	Role         domainUser.Role
}

type Repository interface {
	Create(ctx context.Context, user *domainUser.User, passwordHash string) (string, error)
	GetByEmail(ctx context.Context, email string) (*Credentials, error)
	List(ctx context.Context) ([]domainUser.User, error)
	Get(ctx context.Context, id string) (*domainUser.User, error)
	Update(ctx context.Context, input UpdateInput) error
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

func (r *repository) GetByEmail(ctx context.Context, email string) (*Credentials, error) {
	credentials, err := operations.GetByEmail(ctx, r.dbc, operations.GetByEmailInput{Email: email})
	if err != nil {
		return nil, err
	}

	return &Credentials{
		ID:           credentials.ID,
		PasswordHash: credentials.PasswordHash,
		Role:         credentials.Role,
	}, nil
}

func (r *repository) List(ctx context.Context) ([]domainUser.User, error) {
	return operations.List(ctx, r.dbc)
}

func (r *repository) Get(ctx context.Context, id string) (*domainUser.User, error) {
	return operations.Get(ctx, r.dbc, operations.GetInput{ID: id})
}

func (r *repository) Update(ctx context.Context, input UpdateInput) error {
	return operations.Update(ctx, r.dbc, operations.UpdateInput{
		ID:    input.ID,
		Name:  input.Name,
		Email: input.Email,
	})
}

func (r *repository) Delete(ctx context.Context, id string) error {
	return operations.Delete(ctx, r.dbc, operations.DeleteInput{ID: id})
}

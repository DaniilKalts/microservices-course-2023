package service

import (
	"context"

	"github.com/DaniilKalts/microservices-course-2023/4-week/internal/models"
)

type UserService interface {
	Create(ctx context.Context, user *models.User, password, passwordConfirm string) (string, error)
	Get(ctx context.Context, id string) (*models.User, error)
	Update(ctx context.Context, id string, userPatch *models.UpdateUserPatch) error
	Delete(ctx context.Context, id string) error
}
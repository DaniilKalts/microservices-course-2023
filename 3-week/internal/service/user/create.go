package user

import (
	"context"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/DaniilKalts/microservices-course-2023/3-week/internal/models"
)

func (s *service) Create(ctx context.Context, user *models.User, password string) (string, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return "", err
	}
	user.ID = id.String()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	idStr, err := s.userRepo.Create(ctx, user, string(hashedPassword))
	if err != nil {
		return "", err
	}

	return idStr, nil
}

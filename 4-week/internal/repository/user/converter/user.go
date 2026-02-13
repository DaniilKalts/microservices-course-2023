package converter

import (
	"database/sql"

	"github.com/DaniilKalts/microservices-course-2023/4-week/internal/models"
	modelsRepo "github.com/DaniilKalts/microservices-course-2023/4-week/internal/repository/user/models"
)

func ToUserFromRepo(user *modelsRepo.User) *models.User {
	return &models.User{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      models.Role(user.Role),
		CreatedAt: user.CreatedAt,
		UpdatedAt: &user.UpdatedAt.Time,
	}
}

func ToRepoFromUser(user *models.User) *modelsRepo.User {
	var updatedAt sql.NullTime
	if user.UpdatedAt != nil {
		updatedAt = sql.NullTime{
			Time:  *user.UpdatedAt,
			Valid: true,
		}
	}

	return &modelsRepo.User{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      int32(user.Role),
		CreatedAt: user.CreatedAt,
		UpdatedAt: updatedAt,
	}
}
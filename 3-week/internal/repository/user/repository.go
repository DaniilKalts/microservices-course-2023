package user

import (
	"context"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/DaniilKalts/microservices-course-2023/3-week/internal/models"
	"github.com/DaniilKalts/microservices-course-2023/3-week/internal/repository"
	"github.com/DaniilKalts/microservices-course-2023/3-week/internal/repository/user/converter"
	repoModels "github.com/DaniilKalts/microservices-course-2023/3-week/internal/repository/user/models"
)

type repo struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) repository.UserRepository {
	return &repo{
		db: db,
	}
}

func (r *repo) Create(ctx context.Context, user *models.User, passwordHash string) (string, error) {
	repoUser := converter.ToRepoFromUser(user)
	repoUser.PasswordHash = passwordHash

	builderCreate := sq.Insert("users").
		PlaceholderFormat(sq.Dollar).
		Columns("id", "name", "email", "password_hash", "role").
		Values(repoUser.ID, repoUser.Name, repoUser.Email, repoUser.PasswordHash, repoUser.Role).
		Suffix("RETURNING id")

	query, args, err := builderCreate.ToSql()
	if err != nil {
		return "", err
	}

	var userID string
	if err := r.db.QueryRow(ctx, query, args...).Scan(&userID); err != nil {
		return "", err
	}

	return userID, nil
}

func (r *repo) Get(ctx context.Context, id string) (*models.User, error) {
	builderSelect := sq.Select("id", "name", "email", "role", "created_at", "updated_at").
		From("users").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"id": id}).
		Limit(1)

	query, args, err := builderSelect.ToSql()
	if err != nil {
		return nil, err
	}

	var repoUser repoModels.User
	err = r.db.QueryRow(ctx, query, args...).Scan(
		&repoUser.ID,
		&repoUser.Name,
		&repoUser.Email,
		&repoUser.Role,
		&repoUser.CreatedAt,
		&repoUser.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return converter.ToUserFromRepo(&repoUser), nil
}

func (r *repo) Update(ctx context.Context, id string, userPatch *models.UpdateUserPatch) error {
	builderUpdate := sq.Update("users").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"id": id})

	if userPatch.Name != nil {
		builderUpdate = builderUpdate.Set("name", *userPatch.Name)
	}
	if userPatch.Email != nil {
		builderUpdate = builderUpdate.Set("email", *userPatch.Email)
	}
	if userPatch.PasswordHash != nil {
		builderUpdate = builderUpdate.Set("password_hash", *userPatch.PasswordHash)
	}

	builderUpdate = builderUpdate.Set("updated_at", time.Now())

	query, args, err := builderUpdate.ToSql()
	if err != nil {
		return err
	}

	_, err = r.db.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}

func (r *repo) Delete(ctx context.Context, id string) error {
	builderDelete := sq.Delete("users").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"id": id})

	query, args, err := builderDelete.ToSql()
	if err != nil {
		return err
	}

	_, err = r.db.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}

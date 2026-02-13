package user

import (
	"context"
	"time"

	sq "github.com/Masterminds/squirrel"

	"github.com/DaniilKalts/microservices-course-2023/4-week/internal/clients/database"
	"github.com/DaniilKalts/microservices-course-2023/4-week/internal/models"
	"github.com/DaniilKalts/microservices-course-2023/4-week/internal/repository"
	"github.com/DaniilKalts/microservices-course-2023/4-week/internal/repository/user/converter"
	repoModels "github.com/DaniilKalts/microservices-course-2023/4-week/internal/repository/user/models"
)

type repo struct {
	dbc database.Client
}

func NewRepository(dbc database.Client) repository.UserRepository {
	return &repo{
		dbc: dbc,
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
	if err := r.dbc.DB().ScanOneContext(ctx, &userID, database.Query{Name: "user.Create", QueryRaw: query}, args...); err != nil {
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
	q := database.Query{
		Name:     "user.Get",
		QueryRaw: query,
	}
	if err = r.dbc.DB().ScanOneContext(ctx, &repoUser, q, args...); err != nil {
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

	builderUpdate = builderUpdate.Set("updated_at", time.Now())

	query, args, err := builderUpdate.ToSql()
	if err != nil {
		return err
	}

	q := database.Query{
		Name: "user.Update",
		QueryRaw: query,
	}
	_, err = r.dbc.DB().ExecContext(ctx, q, args...)
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

	q := database.Query{
		Name: "user.Delete",
		QueryRaw: query,
	}
	_, err = r.dbc.DB().ExecContext(ctx, q, args...)
	if err != nil {
		return err
	}

	return nil
}

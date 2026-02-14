package user

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	domainUser "github.com/DaniilKalts/microservices-course-2023/5-week/internal/domain/user"
	repositoryMocks "github.com/DaniilKalts/microservices-course-2023/5-week/internal/repository/mocks"
)

func TestList_Success(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	expected := []domainUser.Entity{
		{ID: "u-1", Name: "John", Email: "john@example.com"},
		{ID: "u-2", Name: "Jane", Email: "jane@example.com"},
	}

	repo := repositoryMocks.NewUserRepositoryMock(t)
	repo.ListMock.Set(func(_ context.Context) ([]domainUser.Entity, error) {
		return expected, nil
	})

	svc := NewService(repo)
	got, err := svc.List(ctx)

	require.NoError(t, err)
	require.Equal(t, expected, got)
	require.Equal(t, uint64(1), repo.ListAfterCounter())
}

func TestList_RepositoryError(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	repoErr := errors.New("repository list failed")

	repo := repositoryMocks.NewUserRepositoryMock(t)
	repo.ListMock.Set(func(_ context.Context) ([]domainUser.Entity, error) {
		return nil, repoErr
	})

	svc := NewService(repo)
	got, err := svc.List(ctx)

	require.EqualError(t, err, repoErr.Error())
	require.Nil(t, got)
	require.Equal(t, uint64(1), repo.ListAfterCounter())
}

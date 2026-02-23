package operations

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	domainUser "github.com/DaniilKalts/microservices-course-2023/6-week/internal/domain/user"
	repositoryMocks "github.com/DaniilKalts/microservices-course-2023/6-week/internal/repository/mocks"
)

func TestList_Success(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	expected := []domainUser.User{
		{ID: "u-1", Name: "John", Email: "john@example.com"},
		{ID: "u-2", Name: "Jane", Email: "jane@example.com"},
	}

	repo := repositoryMocks.NewUserRepositoryMock(t)
	repo.ListMock.Set(func(_ context.Context) ([]domainUser.User, error) {
		return expected, nil
	})

	got, err := List(ctx, repo)

	require.NoError(t, err)
	require.Equal(t, expected, got)
	require.Equal(t, uint64(1), repo.ListAfterCounter())
}

func TestList_RepositoryError(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	repoErr := errors.New("repository list failed")

	repo := repositoryMocks.NewUserRepositoryMock(t)
	repo.ListMock.Set(func(_ context.Context) ([]domainUser.User, error) {
		return nil, repoErr
	})

	got, err := List(ctx, repo)

	require.EqualError(t, err, repoErr.Error())
	require.Nil(t, got)
	require.Equal(t, uint64(1), repo.ListAfterCounter())
}

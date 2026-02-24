package operations

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	domainUser "github.com/DaniilKalts/microservices-course-2023/6-week/internal/domain/user"
	repositoryMocks "github.com/DaniilKalts/microservices-course-2023/6-week/internal/repository/mocks"
)

func TestGet_Success(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	id := "u-1"
	expected := &domainUser.User{ID: id, Name: "John", Email: "john@example.com"}
	repo := repositoryMocks.NewUserRepositoryMock(t)
	repo.GetMock.Set(func(_ context.Context, gotID string) (*domainUser.User, error) {
		require.Equal(t, id, gotID)
		return expected, nil
	})

	got, err := Get(ctx, repo, GetInput{ID: id})

	require.NoError(t, err)
	require.Equal(t, expected, got)
	require.Equal(t, uint64(1), repo.GetAfterCounter())
}

func TestGet_RepositoryScenarios(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	id := "u-1"

	tests := []struct {
		name    string
		repoErr error
	}{
		{
			name:    "TestGet_NotFound",
			repoErr: errors.New("not found"),
		},
		{
			name:    "TestGet_RepositoryError",
			repoErr: errors.New("repository get failed"),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := repositoryMocks.NewUserRepositoryMock(t)
			repo.GetMock.Set(func(_ context.Context, _ string) (*domainUser.User, error) {
				return nil, tt.repoErr
			})

			got, err := Get(ctx, repo, GetInput{ID: id})

			require.EqualError(t, err, tt.repoErr.Error())
			require.Nil(t, got)
			require.Equal(t, uint64(1), repo.GetAfterCounter())
		})
	}
}

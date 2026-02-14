package user

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	domainUser "github.com/DaniilKalts/microservices-course-2023/5-week/internal/domain/user"
	repositoryMocks "github.com/DaniilKalts/microservices-course-2023/5-week/internal/repository/mocks"
)

func TestGet_Success(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	id := "u-1"
	expected := &domainUser.Entity{ID: id, Name: "John", Email: "john@example.com"}
	repo := repositoryMocks.NewUserRepositoryMock(t)
	repo.GetMock.Set(func(_ context.Context, gotID string) (*domainUser.Entity, error) {
		require.Equal(t, id, gotID)
		return expected, nil
	})

	svc := NewService(repo)
	got, err := svc.Get(ctx, id)

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
			repo.GetMock.Set(func(_ context.Context, _ string) (*domainUser.Entity, error) {
				return nil, tt.repoErr
			})

			svc := NewService(repo)
			got, err := svc.Get(ctx, id)

			require.EqualError(t, err, tt.repoErr.Error())
			require.Nil(t, got)
			require.Equal(t, uint64(1), repo.GetAfterCounter())
		})
	}
}

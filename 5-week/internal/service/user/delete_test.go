package user

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDelete_Success(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	id := "u-1"
	var gotID string

	repo := &repoStub{}
	repo.deleteFn = func(_ context.Context, deleteID string) error {
		gotID = deleteID
		return nil
	}

	svc := NewService(repo)
	err := svc.Delete(ctx, id)

	require.NoError(t, err)
	require.Equal(t, 1, repo.deleteCalls)
	require.Equal(t, id, gotID)
}

func TestDelete_RepositoryScenarios(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	id := "u-1"

	tests := []struct {
		name    string
		repoErr error
	}{
		{
			name:    "TestDelete_NotFound",
			repoErr: errors.New("not found"),
		},
		{
			name:    "TestDelete_RepositoryError",
			repoErr: errors.New("repository delete failed"),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := &repoStub{}
			repo.deleteFn = func(_ context.Context, _ string) error {
				return tt.repoErr
			}

			svc := NewService(repo)
			err := svc.Delete(ctx, id)

			require.EqualError(t, err, tt.repoErr.Error())
			require.Equal(t, 1, repo.deleteCalls)
		})
	}
}

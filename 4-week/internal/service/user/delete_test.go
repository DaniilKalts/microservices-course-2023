package user

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/DaniilKalts/microservices-course-2023/4-week/internal/service/user/mocks"
)

func TestDelete_Success(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	id := uuid.NewString()
	var gotID string

	repoMock := mocks.NewUserRepositoryMock(t)
	repoMock.DeleteMock.Set(func(_ context.Context, deleteID string) error {
		gotID = deleteID
		return nil
	})

	svc := NewService(repoMock, nil)
	err := svc.Delete(ctx, id)

	require.NoError(t, err)
	require.Equal(t, uint64(1), repoMock.DeleteBeforeCounter())
	require.Equal(t, id, gotID)
}

func TestDelete_InputValidationScenarios(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		name      string
		id        string
		repoCalls uint64
	}{
		{
			name:      "TestDelete_EmptyID",
			id:        "",
			repoCalls: 0,
		},
		{
			name:      "TestDelete_InvalidIDFormat",
			id:        "not-a-uuid",
			repoCalls: 0,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repoMock := mocks.NewUserRepositoryMock(t)
			repoMock.DeleteMock.Optional().Set(func(_ context.Context, _ string) error {
				return errors.New("repository should not be called")
			})

			svc := NewService(repoMock, nil)
			err := svc.Delete(ctx, tt.id)

			require.Error(t, err)
			require.Equal(t, tt.repoCalls, repoMock.DeleteBeforeCounter())
		})
	}
}

func TestDelete_RepositoryScenarios(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	id := uuid.NewString()

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

			repoMock := mocks.NewUserRepositoryMock(t)
			repoMock.DeleteMock.Set(func(_ context.Context, _ string) error {
				return tt.repoErr
			})

			svc := NewService(repoMock, nil)
			err := svc.Delete(ctx, id)

			require.EqualError(t, err, tt.repoErr.Error())
		})
	}
}

package user

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/DaniilKalts/microservices-course-2023/4-week/internal/models"
	"github.com/DaniilKalts/microservices-course-2023/4-week/internal/service/user/mocks"
)

func TestGet_Success(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	id := uuid.NewString()
	now := time.Now().UTC()
	updatedAt := now.Add(time.Minute)
	expected := &models.User{ID: id, Name: "John", Email: "john@example.com", Role: models.RoleAdmin, CreatedAt: now, UpdatedAt: &updatedAt}

	repoMock := mocks.NewUserRepositoryMock(t)
	repoMock.GetMock.Set(func(_ context.Context, _ string) (*models.User, error) {
		return expected, nil
	})

	svc := NewService(repoMock, nil)
	got, err := svc.Get(ctx, id)

	require.NoError(t, err)
	require.NotNil(t, got)
	require.Equal(t, expected.ID, got.ID)
	require.Equal(t, expected.Name, got.Name)
	require.Equal(t, expected.Email, got.Email)
	require.Equal(t, expected.Role, got.Role)
	require.True(t, got.CreatedAt.Equal(expected.CreatedAt))
	require.NotNil(t, got.UpdatedAt)
	require.True(t, got.UpdatedAt.Equal(*expected.UpdatedAt))
}

func TestGet_InputValidationScenarios(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		name      string
		id        string
		repoCalls uint64
	}{
		{
			name:      "TestGet_EmptyID",
			id:        "",
			repoCalls: 0,
		},
		{
			name:      "TestGet_InvalidIDFormat",
			id:        "not-a-uuid",
			repoCalls: 0,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repoMock := mocks.NewUserRepositoryMock(t)
			repoMock.GetMock.Optional().Set(func(_ context.Context, _ string) (*models.User, error) {
				return nil, errors.New("repository should not be called")
			})

			svc := NewService(repoMock, nil)
			got, err := svc.Get(ctx, tt.id)

			require.Error(t, err)
			require.Nil(t, got)
			require.Equal(t, tt.repoCalls, repoMock.GetBeforeCounter())
		})
	}
}

func TestGet_RepositoryScenarios(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	id := uuid.NewString()

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

			repoMock := mocks.NewUserRepositoryMock(t)
			repoMock.GetMock.Set(func(_ context.Context, _ string) (*models.User, error) {
				return nil, tt.repoErr
			})

			svc := NewService(repoMock, nil)
			got, err := svc.Get(ctx, id)

			require.EqualError(t, err, tt.repoErr.Error())
			require.Nil(t, got)
		})
	}
}

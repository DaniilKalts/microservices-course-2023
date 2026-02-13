package user

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/DaniilKalts/microservices-course-2023/4-week/internal/models"
	"github.com/DaniilKalts/microservices-course-2023/4-week/internal/service/user/mocks"
)

func TestUpdate_Success(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	id := uuid.NewString()
	name := "Jane"
	email := "jane@example.com"
	patch := &models.UpdateUserPatch{Name: &name, Email: &email}

	var gotID string
	var gotPatch *models.UpdateUserPatch

	repoMock := mocks.NewUserRepositoryMock(t)
	repoMock.UpdateMock.Set(func(_ context.Context, updateID string, userPatch *models.UpdateUserPatch) error {
		gotID = updateID
		gotPatch = userPatch
		return nil
	})

	svc := NewService(repoMock, nil)
	err := svc.Update(ctx, id, patch)

	require.NoError(t, err)
	require.Equal(t, uint64(1), repoMock.UpdateBeforeCounter())
	require.Equal(t, id, gotID)
	require.Same(t, patch, gotPatch)
}

func TestUpdate_ValidationScenarios(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	id := uuid.NewString()
	badEmail := "not-an-email"

	tests := []struct {
		name      string
		patch     *models.UpdateUserPatch
		repoCalls uint64
	}{
		{
			name:      "TestUpdate_NilPatch",
			patch:     nil,
			repoCalls: 0,
		},
		{
			name:      "TestUpdate_EmptyPatch",
			patch:     &models.UpdateUserPatch{},
			repoCalls: 0,
		},
		{
			name: "TestUpdate_InvalidEmail",
			patch: &models.UpdateUserPatch{
				Email: &badEmail,
			},
			repoCalls: 0,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repoMock := mocks.NewUserRepositoryMock(t)
			repoMock.UpdateMock.Optional().Set(func(_ context.Context, _ string, _ *models.UpdateUserPatch) error {
				return errors.New("repository should not be called")
			})

			svc := NewService(repoMock, nil)
			err := svc.Update(ctx, id, tt.patch)

			require.Error(t, err)
			require.Equal(t, tt.repoCalls, repoMock.UpdateBeforeCounter())
		})
	}
}

func TestUpdate_RepositoryScenarios(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	id := uuid.NewString()
	name := "Jane"
	patch := &models.UpdateUserPatch{Name: &name}
	email := "john@example.com"
	patchWithEmail := &models.UpdateUserPatch{Email: &email}

	tests := []struct {
		name    string
		patch   *models.UpdateUserPatch
		repoErr error
	}{
		{
			name:    "TestUpdate_NotFound",
			patch:   patch,
			repoErr: errors.New("not found"),
		},
		{
			name:    "TestUpdate_DuplicateEmail",
			patch:   patchWithEmail,
			repoErr: errors.New("duplicate email"),
		},
		{
			name:    "TestUpdate_RepositoryError",
			patch:   patch,
			repoErr: errors.New("repository update failed"),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repoMock := mocks.NewUserRepositoryMock(t)
			repoMock.UpdateMock.Set(func(_ context.Context, _ string, _ *models.UpdateUserPatch) error {
				return tt.repoErr
			})

			svc := NewService(repoMock, nil)
			err := svc.Update(ctx, id, tt.patch)

			require.EqualError(t, err, tt.repoErr.Error())
		})
	}
}

func TestUpdate_SetsUpdatedAt(t *testing.T) {
	t.Parallel()
	t.Skip("UpdateUserPatch has no UpdatedAt field in current contract")
}

func TestUpdate_PartialPatch_DoesNotOverwriteOtherFields(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	id := uuid.NewString()
	name := "OnlyNameChanged"
	patch := &models.UpdateUserPatch{Name: &name}

	var gotPatch *models.UpdateUserPatch

	repoMock := mocks.NewUserRepositoryMock(t)
	repoMock.UpdateMock.Set(func(_ context.Context, _ string, userPatch *models.UpdateUserPatch) error {
		gotPatch = userPatch
		return nil
	})

	svc := NewService(repoMock, nil)
	err := svc.Update(ctx, id, patch)

	require.NoError(t, err)
	require.NotNil(t, gotPatch)
	require.NotNil(t, gotPatch.Name)
	require.Equal(t, name, *gotPatch.Name)
	require.Nil(t, gotPatch.Email)
}

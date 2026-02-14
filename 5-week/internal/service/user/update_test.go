package user

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	domainUser "github.com/DaniilKalts/microservices-course-2023/5-week/internal/domain/user"
	repositoryMocks "github.com/DaniilKalts/microservices-course-2023/5-week/internal/repository/mocks"
)

func TestUpdate_Success(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	id := "u-1"
	name := "Jane"
	email := "jane@example.com"
	patch := &domainUser.UpdatePatch{Name: &name, Email: &email}

	var gotID string
	var gotPatch *domainUser.UpdatePatch

	repo := repositoryMocks.NewUserRepositoryMock(t)
	repo.UpdateMock.Set(func(_ context.Context, updateID string, userPatch *domainUser.UpdatePatch) error {
		gotID = updateID
		gotPatch = userPatch
		return nil
	})

	svc := NewService(repo)
	err := svc.Update(ctx, id, patch)

	require.NoError(t, err)
	require.Equal(t, uint64(1), repo.UpdateAfterCounter())
	require.Equal(t, id, gotID)
	require.Same(t, patch, gotPatch)
}

func TestUpdate_RepositoryScenarios(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	id := "u-1"
	name := "Jane"
	patch := &domainUser.UpdatePatch{Name: &name}
	email := "john@example.com"
	patchWithEmail := &domainUser.UpdatePatch{Email: &email}

	tests := []struct {
		name    string
		patch   *domainUser.UpdatePatch
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

			repo := repositoryMocks.NewUserRepositoryMock(t)
			repo.UpdateMock.Set(func(_ context.Context, _ string, _ *domainUser.UpdatePatch) error {
				return tt.repoErr
			})

			svc := NewService(repo)
			err := svc.Update(ctx, id, tt.patch)

			require.EqualError(t, err, tt.repoErr.Error())
			require.Equal(t, uint64(1), repo.UpdateAfterCounter())
		})
	}
}

func TestUpdate_PartialPatch_DoesNotOverwriteOtherFields(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	id := "u-1"
	name := "OnlyNameChanged"
	patch := &domainUser.UpdatePatch{Name: &name}

	var gotPatch *domainUser.UpdatePatch

	repo := repositoryMocks.NewUserRepositoryMock(t)
	repo.UpdateMock.Set(func(_ context.Context, _ string, userPatch *domainUser.UpdatePatch) error {
		gotPatch = userPatch
		return nil
	})

	svc := NewService(repo)
	err := svc.Update(ctx, id, patch)

	require.NoError(t, err)
	require.NotNil(t, gotPatch)
	require.NotNil(t, gotPatch.Name)
	require.Equal(t, name, *gotPatch.Name)
	require.Nil(t, gotPatch.Email)
}

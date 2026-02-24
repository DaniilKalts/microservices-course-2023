package operations

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/DaniilKalts/microservices-course-2023/6-week/internal/repository"
	repositoryMocks "github.com/DaniilKalts/microservices-course-2023/6-week/internal/repository/mocks"
)

func TestUpdate_Success(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	id := "u-1"
	name := "Jane"
	email := "jane@example.com"
	input := UpdateInput{ID: id, Name: &name, Email: &email}

	var gotInput repository.UserUpdateInput

	repo := repositoryMocks.NewUserRepositoryMock(t)
	repo.UpdateMock.Set(func(_ context.Context, updateInput repository.UserUpdateInput) error {
		gotInput = updateInput
		return nil
	})

	err := Update(ctx, repo, input)

	require.NoError(t, err)
	require.Equal(t, uint64(1), repo.UpdateAfterCounter())
	require.Equal(t, repository.UserUpdateInput{ID: input.ID, Name: input.Name, Email: input.Email}, gotInput)
}

func TestUpdate_RepositoryScenarios(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	id := "u-1"
	name := "Jane"
	updateInput := UpdateInput{ID: id, Name: &name}
	email := "john@example.com"
	updateInputWithEmail := UpdateInput{ID: id, Email: &email}

	tests := []struct {
		name    string
		input   UpdateInput
		repoErr error
	}{
		{
			name:    "TestUpdate_NotFound",
			input:   updateInput,
			repoErr: errors.New("not found"),
		},
		{
			name:    "TestUpdate_DuplicateEmail",
			input:   updateInputWithEmail,
			repoErr: errors.New("duplicate email"),
		},
		{
			name:    "TestUpdate_RepositoryError",
			input:   updateInput,
			repoErr: errors.New("repository update failed"),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := repositoryMocks.NewUserRepositoryMock(t)
			repo.UpdateMock.Set(func(_ context.Context, _ repository.UserUpdateInput) error {
				return tt.repoErr
			})

			err := Update(ctx, repo, tt.input)

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
	input := UpdateInput{ID: id, Name: &name}

	var gotInput repository.UserUpdateInput

	repo := repositoryMocks.NewUserRepositoryMock(t)
	repo.UpdateMock.Set(func(_ context.Context, updateInput repository.UserUpdateInput) error {
		gotInput = updateInput
		return nil
	})

	err := Update(ctx, repo, input)

	require.NoError(t, err)
	require.NotNil(t, gotInput.Name)
	require.Equal(t, name, *gotInput.Name)
	require.Nil(t, gotInput.Email)
}

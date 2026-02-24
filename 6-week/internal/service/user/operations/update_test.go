package operations

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	repositoryMocks "github.com/DaniilKalts/microservices-course-2023/6-week/internal/repository/mocks"
	userRepository "github.com/DaniilKalts/microservices-course-2023/6-week/internal/repository/user"
)

func TestUpdate_Success(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	id := "u-1"
	name := "Jane"
	email := "jane@example.com"
	input := UpdateInput{ID: id, Name: &name, Email: &email}

	var gotInput userRepository.UpdateInput

	repo := repositoryMocks.NewUserRepositoryMock(t)
	repo.UpdateMock.Set(func(_ context.Context, updateInput userRepository.UpdateInput) error {
		gotInput = updateInput
		return nil
	})

	err := Update(ctx, repo, input)

	require.NoError(t, err)
	require.Equal(t, uint64(1), repo.UpdateAfterCounter())
	require.Equal(t, userRepository.UpdateInput{ID: input.ID, Name: input.Name, Email: input.Email}, gotInput)
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
			repo.UpdateMock.Set(func(_ context.Context, _ userRepository.UpdateInput) error {
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

	var gotInput userRepository.UpdateInput

	repo := repositoryMocks.NewUserRepositoryMock(t)
	repo.UpdateMock.Set(func(_ context.Context, updateInput userRepository.UpdateInput) error {
		gotInput = updateInput
		return nil
	})

	err := Update(ctx, repo, input)

	require.NoError(t, err)
	require.NotNil(t, gotInput.Name)
	require.Equal(t, name, *gotInput.Name)
	require.Nil(t, gotInput.Email)
}

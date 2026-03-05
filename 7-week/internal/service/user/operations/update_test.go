package operations

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"

	repositoryMocks "github.com/DaniilKalts/microservices-course-2023/7-week/internal/repository/mocks"
	userRepositoryOperations "github.com/DaniilKalts/microservices-course-2023/7-week/internal/repository/user/operations"
)

func TestUpdate_Success(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	id := "u-1"
	name := "Jane"
	email := "jane@example.com"
	input := UpdateInput{ID: id, Name: &name, Email: &email}

	var gotInput userRepositoryOperations.UpdateInput

	repo := repositoryMocks.NewUserRepositoryMock(t)
	repo.UpdateMock.Set(func(_ context.Context, updateInput userRepositoryOperations.UpdateInput) error {
		gotInput = updateInput
		return nil
	})

	err := Update(ctx, repo, input)

	require.NoError(t, err)
	require.Equal(t, uint64(1), repo.UpdateAfterCounter())
	require.Equal(t, userRepositoryOperations.UpdateInput{ID: input.ID, Name: input.Name, Email: input.Email}, gotInput)
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
			repo.UpdateMock.Set(func(_ context.Context, _ userRepositoryOperations.UpdateInput) error {
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

	var gotInput userRepositoryOperations.UpdateInput

	repo := repositoryMocks.NewUserRepositoryMock(t)
	repo.UpdateMock.Set(func(_ context.Context, updateInput userRepositoryOperations.UpdateInput) error {
		gotInput = updateInput
		return nil
	})

	err := Update(ctx, repo, input)

	require.NoError(t, err)
	require.NotNil(t, gotInput.Name)
	require.Equal(t, name, *gotInput.Name)
	require.Nil(t, gotInput.Email)
}

func TestUpdate_WithPassword_HashesBeforeRepositoryCall(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	id := "u-1"
	password := "Secret123"
	input := UpdateInput{ID: id, Password: &password}

	var gotInput userRepositoryOperations.UpdateInput

	repo := repositoryMocks.NewUserRepositoryMock(t)
	repo.UpdateMock.Set(func(_ context.Context, updateInput userRepositoryOperations.UpdateInput) error {
		gotInput = updateInput
		return nil
	})

	err := Update(ctx, repo, input)

	require.NoError(t, err)
	require.NotNil(t, gotInput.PasswordHash)
	require.NoError(t, bcrypt.CompareHashAndPassword([]byte(*gotInput.PasswordHash), []byte(password)))
	require.Equal(t, uint64(1), repo.UpdateAfterCounter())
}

func TestUpdate_WithPasswordHashError(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	id := "u-1"
	password := strings.Repeat("a", 73)
	input := UpdateInput{ID: id, Password: &password}

	repo := repositoryMocks.NewUserRepositoryMock(t)

	err := Update(ctx, repo, input)

	require.Error(t, err)
	require.Equal(t, uint64(0), repo.UpdateAfterCounter())
}

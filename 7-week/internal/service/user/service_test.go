package user

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"

	domainUser "github.com/DaniilKalts/microservices-course-2023/7-week/internal/domain/user"
	repositoryMocks "github.com/DaniilKalts/microservices-course-2023/7-week/internal/repository/mocks"
	userRepository "github.com/DaniilKalts/microservices-course-2023/7-week/internal/repository/user"
)

func newTestService(t *testing.T) (*repositoryMocks.UserRepositoryMock, Service) {
	t.Helper()
	repo := repositoryMocks.NewUserRepositoryMock(t)
	svc := NewService(repo)
	return repo, svc
}

func TestCreate_ValidationScenarios(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		name            string
		password        string
		passwordConfirm string
		wantErr   error
		repoCalls int
	}{
		{
			name:            "TestCreate_PasswordMismatch",
			password:        "P@ssword123",
			passwordConfirm: "P@ssword321",
			wantErr:         ErrPasswordMismatch,
			repoCalls:       0,
		},
		{
			name:            "TestCreate_BcryptTooLongPassword",
			password:        strings.Repeat("a", 73),
			passwordConfirm: strings.Repeat("a", 73),
			repoCalls:       0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo, svc := newTestService(t)
			repo.CreateMock.Optional().Set(func(_ context.Context, _ *domainUser.User, _ string) (string, error) {
				return "", errors.New("repository should not be called")
			})

			gotID, err := svc.Create(ctx, CreateInput{
				User:            &domainUser.User{Name: "John", Email: "john@example.com"},
				Password:        tt.password,
				PasswordConfirm: tt.passwordConfirm,
			})

			require.Error(t, err)
			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
			}
			require.Empty(t, gotID)
			require.Equal(t, uint64(tt.repoCalls), repo.CreateAfterCounter())
		})
	}
}

func TestCreate_RepositoryScenarios(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		name       string
		repoResult string
		repoErr    error
		wantID     string
		wantErr    string
	}{
		{
			name:       "TestCreate_Success",
			repoResult: "new-user-id",
			wantID:     "new-user-id",
		},
		{
			name:    "TestCreate_DuplicateEmail",
			repoErr: errors.New("duplicate email"),
			wantErr: "duplicate email",
		},
		{
			name:    "TestCreate_RepositoryError",
			repoErr: errors.New("repository create failed"),
			wantErr: "repository create failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo, svc := newTestService(t)
			repo.CreateMock.Set(func(_ context.Context, user *domainUser.User, passwordHash string) (string, error) {
				require.NotEmpty(t, user.ID)
				require.NotEmpty(t, passwordHash)
				return tt.repoResult, tt.repoErr
			})

			gotID, err := svc.Create(ctx, CreateInput{
				User:            &domainUser.User{Name: "John", Email: "john@example.com"},
				Password:        "P@ssword123",
				PasswordConfirm: "P@ssword123",
			})

			if tt.wantErr != "" {
				require.EqualError(t, err, tt.wantErr)
				require.Empty(t, gotID)
				require.Equal(t, uint64(1), repo.CreateAfterCounter())
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.wantID, gotID)
			require.Equal(t, uint64(1), repo.CreateAfterCounter())
		})
	}
}

func TestGet_Success(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	id := "u-1"
	expected := &domainUser.User{ID: id, Name: "John", Email: "john@example.com"}

	repo, svc := newTestService(t)
	repo.GetByIDMock.Set(func(_ context.Context, gotID string) (*domainUser.User, error) {
		require.Equal(t, id, gotID)
		return expected, nil
	})

	got, err := svc.Get(ctx, id)

	require.NoError(t, err)
	require.Equal(t, expected, got)
	require.Equal(t, uint64(1), repo.GetByIDAfterCounter())
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
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo, svc := newTestService(t)
			repo.GetByIDMock.Set(func(_ context.Context, _ string) (*domainUser.User, error) {
				return nil, tt.repoErr
			})

			got, err := svc.Get(ctx, id)

			require.EqualError(t, err, tt.repoErr.Error())
			require.Nil(t, got)
			require.Equal(t, uint64(1), repo.GetByIDAfterCounter())
		})
	}
}

func TestList_Success(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	expected := []domainUser.User{
		{ID: "u-1", Name: "John", Email: "john@example.com"},
		{ID: "u-2", Name: "Jane", Email: "jane@example.com"},
	}

	repo, svc := newTestService(t)
	repo.ListMock.Set(func(_ context.Context) ([]domainUser.User, error) {
		return expected, nil
	})

	got, err := svc.List(ctx)

	require.NoError(t, err)
	require.Equal(t, expected, got)
	require.Equal(t, uint64(1), repo.ListAfterCounter())
}

func TestList_RepositoryError(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	repoErr := errors.New("repository list failed")

	repo, svc := newTestService(t)
	repo.ListMock.Set(func(_ context.Context) ([]domainUser.User, error) {
		return nil, repoErr
	})

	got, err := svc.List(ctx)

	require.EqualError(t, err, repoErr.Error())
	require.Nil(t, got)
	require.Equal(t, uint64(1), repo.ListAfterCounter())
}

func TestDelete_Success(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	id := "u-1"
	var gotID string

	repo, svc := newTestService(t)
	repo.DeleteMock.Set(func(_ context.Context, deleteID string) error {
		gotID = deleteID
		return nil
	})

	err := svc.Delete(ctx, id)

	require.NoError(t, err)
	require.Equal(t, uint64(1), repo.DeleteAfterCounter())
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
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo, svc := newTestService(t)
			repo.DeleteMock.Set(func(_ context.Context, _ string) error {
				return tt.repoErr
			})

			err := svc.Delete(ctx, id)

			require.EqualError(t, err, tt.repoErr.Error())
			require.Equal(t, uint64(1), repo.DeleteAfterCounter())
		})
	}
}

func TestUpdate_Success(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	id := "u-1"
	name := "Jane"
	email := "jane@example.com"
	input := UpdateInput{ID: id, Name: &name, Email: &email}

	var gotInput userRepository.UpdateInput

	repo, svc := newTestService(t)
	repo.UpdateMock.Set(func(_ context.Context, updateInput userRepository.UpdateInput) error {
		gotInput = updateInput
		return nil
	})

	err := svc.Update(ctx, input)

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
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo, svc := newTestService(t)
			repo.UpdateMock.Set(func(_ context.Context, _ userRepository.UpdateInput) error {
				return tt.repoErr
			})

			err := svc.Update(ctx, tt.input)

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

	repo, svc := newTestService(t)
	repo.UpdateMock.Set(func(_ context.Context, updateInput userRepository.UpdateInput) error {
		gotInput = updateInput
		return nil
	})

	err := svc.Update(ctx, input)

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

	var gotInput userRepository.UpdateInput

	repo, svc := newTestService(t)
	repo.UpdateMock.Set(func(_ context.Context, updateInput userRepository.UpdateInput) error {
		gotInput = updateInput
		return nil
	})

	err := svc.Update(ctx, input)

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

	repo, svc := newTestService(t)

	err := svc.Update(ctx, input)

	require.Error(t, err)
	require.Equal(t, uint64(0), repo.UpdateAfterCounter())
}

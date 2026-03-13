package user

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"

	domainUser "github.com/DaniilKalts/microservices-course-2023/8-week/internal/domain/user"
	userRepo "github.com/DaniilKalts/microservices-course-2023/8-week/internal/repository/user"
)

// --- Helpers ---

func newTestService(t *testing.T) (*userRepo.UserRepositoryMock, Service) {
	t.Helper()
	repo := userRepo.NewUserRepositoryMock(t)
	return repo, NewService(repo)
}

func ptr[T any](v T) *T { return &v }

// --- Tests ---

func TestCreate(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		repo, svc := newTestService(t)

		repo.CreateMock.Set(func(_ context.Context, user *domainUser.User, hash string) (string, error) {
			assert.NotEmpty(t, user.ID)
			assert.Equal(t, "John", user.Name)
			assert.Equal(t, "john@example.com", user.Email)
			assert.NoError(t, bcrypt.CompareHashAndPassword([]byte(hash), []byte("P@ssword123")))
			return user.ID, nil
		})

		id, err := svc.Create(context.Background(), domainUser.CreateInput{
			User:     &domainUser.User{Name: "John", Email: "john@example.com"},
			Password: "P@ssword123",
		})

		require.NoError(t, err)
		assert.NotEmpty(t, id)
	})

	t.Run("password too long for bcrypt", func(t *testing.T) {
		t.Parallel()
		repo, svc := newTestService(t)
		repo.CreateMock.Optional()

		_, err := svc.Create(context.Background(), domainUser.CreateInput{
			User:     &domainUser.User{Name: "John", Email: "john@example.com"},
			Password: strings.Repeat("a", 73),
		})

		require.Error(t, err)
		assert.Equal(t, uint64(0), repo.CreateAfterCounter())
	})

	t.Run("repo errors", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name    string
			repoErr error
			wantErr error
		}{
			{"duplicate email", domainUser.ErrEmailAlreadyExists, domainUser.ErrEmailAlreadyExists},
			{"generic error", errors.New("connection refused"), nil},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				repo, svc := newTestService(t)
				repo.CreateMock.Return("", tt.repoErr)

				id, err := svc.Create(context.Background(), domainUser.CreateInput{
					User:     &domainUser.User{Name: "John", Email: "john@example.com"},
					Password: "P@ssword123",
				})

				require.Error(t, err)
				if tt.wantErr != nil {
					assert.ErrorIs(t, err, tt.wantErr)
				}
				assert.Empty(t, id)
			})
		}
	})
}

func TestGet(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		repo, svc := newTestService(t)

		expected := &domainUser.User{ID: "u-1", Name: "John", Email: "john@example.com"}
		repo.GetByIDMock.Set(func(_ context.Context, id string) (*domainUser.User, error) {
			assert.Equal(t, "u-1", id)
			return expected, nil
		})

		got, err := svc.Get(context.Background(), "u-1")

		require.NoError(t, err)
		assert.Equal(t, expected, got)
	})

	t.Run("repo errors", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name    string
			repoErr error
			wantErr error
		}{
			{"not found", domainUser.ErrNotFound, domainUser.ErrNotFound},
			{"generic error", errors.New("connection refused"), nil},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				repo, svc := newTestService(t)
				repo.GetByIDMock.Return(nil, tt.repoErr)

				got, err := svc.Get(context.Background(), "u-1")

				require.Error(t, err)
				if tt.wantErr != nil {
					assert.ErrorIs(t, err, tt.wantErr)
				}
				assert.Nil(t, got)
			})
		}
	})
}

func TestGetCredentialsByEmail(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		repo, svc := newTestService(t)

		expected := &domainUser.Credentials{
			ID:           "u-1",
			PasswordHash: "$2a$10$hash",
			Role:         domainUser.RoleUser,
		}
		repo.GetCredentialsByEmailMock.Set(func(_ context.Context, email string) (*domainUser.Credentials, error) {
			assert.Equal(t, "john@example.com", email)
			return expected, nil
		})

		got, err := svc.GetCredentialsByEmail(context.Background(), "john@example.com")

		require.NoError(t, err)
		assert.Equal(t, expected, got)
	})

	t.Run("repo errors", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name    string
			repoErr error
			wantErr error
		}{
			{"not found", domainUser.ErrNotFound, domainUser.ErrNotFound},
			{"generic error", errors.New("connection refused"), nil},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				repo, svc := newTestService(t)
				repo.GetCredentialsByEmailMock.Return(nil, tt.repoErr)

				got, err := svc.GetCredentialsByEmail(context.Background(), "john@example.com")

				require.Error(t, err)
				if tt.wantErr != nil {
					assert.ErrorIs(t, err, tt.wantErr)
				}
				assert.Nil(t, got)
			})
		}
	})
}

func TestList(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		repo, svc := newTestService(t)

		expected := []domainUser.User{
			{ID: "u-1", Name: "John", Email: "john@example.com"},
			{ID: "u-2", Name: "Jane", Email: "jane@example.com"},
		}
		repo.ListMock.Return(expected, nil)

		got, err := svc.List(context.Background())

		require.NoError(t, err)
		assert.Equal(t, expected, got)
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		repo, svc := newTestService(t)
		repo.ListMock.Return([]domainUser.User{}, nil)

		got, err := svc.List(context.Background())

		require.NoError(t, err)
		assert.Empty(t, got)
	})

	t.Run("repo error", func(t *testing.T) {
		t.Parallel()
		repo, svc := newTestService(t)
		repo.ListMock.Return(nil, errors.New("connection refused"))

		got, err := svc.List(context.Background())

		require.Error(t, err)
		assert.Nil(t, got)
	})
}

func TestUpdate(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		repo, svc := newTestService(t)

		repo.UpdateMock.Set(func(_ context.Context, input domainUser.UpdateInput) error {
			assert.Equal(t, "u-1", input.ID)
			assert.Equal(t, ptr("Jane"), input.Name)
			assert.Equal(t, ptr("jane@example.com"), input.Email)
			assert.Nil(t, input.PasswordHash)
			return nil
		})

		err := svc.Update(context.Background(), domainUser.UpdateInput{
			ID:    "u-1",
			Name:  ptr("Jane"),
			Email: ptr("jane@example.com"),
		})

		require.NoError(t, err)
	})

	t.Run("hashes password before repo call", func(t *testing.T) {
		t.Parallel()
		repo, svc := newTestService(t)

		repo.UpdateMock.Set(func(_ context.Context, input domainUser.UpdateInput) error {
			require.NotNil(t, input.PasswordHash)
			assert.NoError(t, bcrypt.CompareHashAndPassword([]byte(*input.PasswordHash), []byte("Secret123")))
			return nil
		})

		err := svc.Update(context.Background(), domainUser.UpdateInput{ID: "u-1", Password: ptr("Secret123")})

		require.NoError(t, err)
	})

	t.Run("partial patch preserves nil fields", func(t *testing.T) {
		t.Parallel()
		repo, svc := newTestService(t)

		repo.UpdateMock.Set(func(_ context.Context, input domainUser.UpdateInput) error {
			assert.Equal(t, ptr("OnlyName"), input.Name)
			assert.Nil(t, input.Email)
			assert.Nil(t, input.PasswordHash)
			return nil
		})

		err := svc.Update(context.Background(), domainUser.UpdateInput{ID: "u-1", Name: ptr("OnlyName")})

		require.NoError(t, err)
	})

	t.Run("password too long for bcrypt", func(t *testing.T) {
		t.Parallel()
		repo, svc := newTestService(t)
		repo.UpdateMock.Optional()

		err := svc.Update(context.Background(), domainUser.UpdateInput{
			ID:       "u-1",
			Password: ptr(strings.Repeat("a", 73)),
		})

		require.Error(t, err)
		assert.Equal(t, uint64(0), repo.UpdateAfterCounter())
	})

	t.Run("repo errors", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name    string
			input   domainUser.UpdateInput
			repoErr error
			wantErr error
		}{
			{"not found", domainUser.UpdateInput{ID: "u-1", Name: ptr("Jane")}, domainUser.ErrNotFound, domainUser.ErrNotFound},
			{"duplicate email", domainUser.UpdateInput{ID: "u-1", Email: ptr("taken@example.com")}, domainUser.ErrEmailAlreadyExists, domainUser.ErrEmailAlreadyExists},
			{"no fields", domainUser.UpdateInput{ID: "u-1", Name: ptr("Jane")}, domainUser.ErrNoFieldsToUpdate, domainUser.ErrNoFieldsToUpdate},
			{"generic error", domainUser.UpdateInput{ID: "u-1", Name: ptr("Jane")}, errors.New("connection refused"), nil},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				repo, svc := newTestService(t)
				repo.UpdateMock.Return(tt.repoErr)

				err := svc.Update(context.Background(), tt.input)

				require.Error(t, err)
				if tt.wantErr != nil {
					assert.ErrorIs(t, err, tt.wantErr)
				}
			})
		}
	})
}

func TestDelete(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		repo, svc := newTestService(t)

		repo.DeleteMock.Set(func(_ context.Context, id string) error {
			assert.Equal(t, "u-1", id)
			return nil
		})

		err := svc.Delete(context.Background(), "u-1")

		require.NoError(t, err)
	})

	t.Run("repo errors", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name    string
			repoErr error
			wantErr error
		}{
			{"not found", domainUser.ErrNotFound, domainUser.ErrNotFound},
			{"generic error", errors.New("connection refused"), nil},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				repo, svc := newTestService(t)
				repo.DeleteMock.Return(tt.repoErr)

				err := svc.Delete(context.Background(), "u-1")

				require.Error(t, err)
				if tt.wantErr != nil {
					assert.ErrorIs(t, err, tt.wantErr)
				}
			})
		}
	})
}

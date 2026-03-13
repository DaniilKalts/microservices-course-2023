package user

import (
	"context"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"go.uber.org/zap"

	pgClient "github.com/DaniilKalts/microservices-course-2023/8-week/internal/adapters/database/postgres"
	"github.com/DaniilKalts/microservices-course-2023/8-week/internal/clients/database"
	"github.com/DaniilKalts/microservices-course-2023/8-week/internal/config"
	domainUser "github.com/DaniilKalts/microservices-course-2023/8-week/internal/domain/user"
)

// --- Helpers ---

func newTestRepoWithTimeout(t *testing.T, queryTimeout time.Duration) (Repository, context.Context) {
	t.Helper()
	ctx := context.Background()

	// Start postgres container
	pgContainer, err := postgres.Run(ctx, "postgres:17",
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("testuser"),
		postgres.WithPassword("testpass"),
		postgres.BasicWaitStrategies(),
	)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, pgContainer.Terminate(ctx)) })

	dsn, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

	host, err := pgContainer.Host(ctx)
	require.NoError(t, err)

	port, err := pgContainer.MappedPort(ctx, "5432")
	require.NoError(t, err)

	// Connect database client
	client, err := pgClient.New(ctx, config.PostgresConfig{
		Host:         host,
		Port:         port.Port(),
		User:         "testuser",
		Password:     "testpass",
		Name:         "testdb",
		SSLMode:      "disable",
		QueryTimeout: queryTimeout,
	}, zap.NewNop())
	require.NoError(t, err)
	t.Cleanup(func() { client.Close() })

	// Run migrations
	_, here, _, _ := runtime.Caller(0)
	migrationsDir := filepath.Join(filepath.Dir(here), "..", "..", "..", "deployments", "migrations")

	db, err := goose.OpenDBWithDriver("pgx", dsn)
	require.NoError(t, err)
	t.Cleanup(func() { db.Close() })

	require.NoError(t, goose.Up(db, migrationsDir))

	return NewRepository(client), ctx
}

func newTestRepo(t *testing.T) (Repository, context.Context) {
	return newTestRepoWithTimeout(t, 5*time.Second)
}

func randomUser() *domainUser.User {
	return &domainUser.User{
		ID:    uuid.New().String(),
		Name:  "John Doe",
		Email: uuid.New().String() + "@test.com",
		Role:  domainUser.RoleUser,
	}
}

func createUser(t *testing.T, repo Repository, ctx context.Context) *domainUser.User {
	t.Helper()
	user := randomUser()
	_, err := repo.Create(ctx, user, "hashedpass")
	require.NoError(t, err)
	return user
}

func ptr[T any](v T) *T { return &v }

// --- Tests ---

func TestCreate(t *testing.T) {
	repo, ctx := newTestRepo(t)

	t.Run("success", func(t *testing.T) {
		user := randomUser()

		id, err := repo.Create(ctx, user, "hashedpass")

		require.NoError(t, err)
		assert.Equal(t, user.ID, id)
	})

	t.Run("duplicate email", func(t *testing.T) {
		existing := createUser(t, repo, ctx)
		duplicate := randomUser()
		duplicate.Email = existing.Email

		_, err := repo.Create(ctx, duplicate, "hashedpass")

		assert.ErrorIs(t, err, domainUser.ErrEmailAlreadyExists)
	})
}

func TestGetByID(t *testing.T) {
	repo, ctx := newTestRepo(t)

	t.Run("success", func(t *testing.T) {
		existing := createUser(t, repo, ctx)

		found, err := repo.GetByID(ctx, existing.ID)

		require.NoError(t, err)
		assert.Equal(t, existing.ID, found.ID)
		assert.Equal(t, existing.Name, found.Name)
		assert.Equal(t, existing.Email, found.Email)
		assert.Equal(t, existing.Role, found.Role)
		assert.WithinDuration(t, time.Now(), found.CreatedAt, 5*time.Second)
	})

	t.Run("not found", func(t *testing.T) {
		_, err := repo.GetByID(ctx, uuid.New().String())

		assert.ErrorIs(t, err, domainUser.ErrNotFound)
	})
}

func TestList(t *testing.T) {
	repo, ctx := newTestRepo(t)

	first := createUser(t, repo, ctx)
	second := createUser(t, repo, ctx)

	users, err := repo.List(ctx)

	require.NoError(t, err)
	assert.Len(t, users, 2)
	assert.Equal(t, second.ID, users[0].ID)
	assert.Equal(t, first.ID, users[1].ID)
}

func TestGetCredentialsByEmail(t *testing.T) {
	repo, ctx := newTestRepo(t)

	t.Run("success", func(t *testing.T) {
		existing := createUser(t, repo, ctx)

		creds, err := repo.GetCredentialsByEmail(ctx, existing.Email)

		require.NoError(t, err)
		assert.Equal(t, existing.ID, creds.ID)
		assert.Equal(t, "hashedpass", creds.PasswordHash)
		assert.Equal(t, existing.Role, creds.Role)
	})

	t.Run("not found", func(t *testing.T) {
		_, err := repo.GetCredentialsByEmail(ctx, "nobody@test.com")

		assert.ErrorIs(t, err, domainUser.ErrNotFound)
	})
}

func TestUpdate(t *testing.T) {
	repo, ctx := newTestRepo(t)

	t.Run("success", func(t *testing.T) {
		existing := createUser(t, repo, ctx)

		err := repo.Update(ctx, domainUser.UpdateInput{
			ID:    existing.ID,
			Name:  ptr("Jane Doe"),
			Email: ptr(uuid.New().String() + "@updated.com"),
		})
		require.NoError(t, err)

		found, err := repo.GetByID(ctx, existing.ID)
		require.NoError(t, err)
		assert.Equal(t, "Jane Doe", found.Name)
		assert.NotEqual(t, existing.Email, found.Email)
	})

	t.Run("not found", func(t *testing.T) {
		err := repo.Update(ctx, domainUser.UpdateInput{ID: uuid.New().String(), Name: ptr("Jane")})

		assert.ErrorIs(t, err, domainUser.ErrNotFound)
	})

	t.Run("no fields", func(t *testing.T) {
		err := repo.Update(ctx, domainUser.UpdateInput{ID: uuid.New().String()})

		assert.ErrorIs(t, err, domainUser.ErrNoFieldsToUpdate)
	})

	t.Run("duplicate email", func(t *testing.T) {
		first := createUser(t, repo, ctx)
		second := createUser(t, repo, ctx)

		err := repo.Update(ctx, domainUser.UpdateInput{ID: second.ID, Email: &first.Email})

		assert.ErrorIs(t, err, domainUser.ErrEmailAlreadyExists)
	})
}

func TestDelete(t *testing.T) {
	repo, ctx := newTestRepo(t)

	t.Run("success", func(t *testing.T) {
		existing := createUser(t, repo, ctx)

		err := repo.Delete(ctx, existing.ID)

		require.NoError(t, err)
		_, err = repo.GetByID(ctx, existing.ID)
		assert.ErrorIs(t, err, domainUser.ErrNotFound)
	})

	t.Run("not found", func(t *testing.T) {
		err := repo.Delete(ctx, uuid.New().String())

		assert.ErrorIs(t, err, domainUser.ErrNotFound)
	})
}

func TestQueryTimeout(t *testing.T) {
	repo, ctx := newTestRepoWithTimeout(t, time.Nanosecond)

	t.Run("create", func(t *testing.T) {
		_, err := repo.Create(ctx, randomUser(), "hashedpass")

		assert.ErrorIs(t, err, database.ErrTimeout)
	})

	t.Run("list", func(t *testing.T) {
		_, err := repo.List(ctx)

		assert.ErrorIs(t, err, database.ErrTimeout)
	})

	t.Run("get by id", func(t *testing.T) {
		_, err := repo.GetByID(ctx, uuid.New().String())

		assert.ErrorIs(t, err, database.ErrTimeout)
	})

	t.Run("get credentials by email", func(t *testing.T) {
		_, err := repo.GetCredentialsByEmail(ctx, "test@test.com")

		assert.ErrorIs(t, err, database.ErrTimeout)
	})

	t.Run("update", func(t *testing.T) {
		err := repo.Update(ctx, domainUser.UpdateInput{ID: uuid.New().String(), Name: ptr("Jane")})

		assert.ErrorIs(t, err, database.ErrTimeout)
	})

	t.Run("delete", func(t *testing.T) {
		err := repo.Delete(ctx, uuid.New().String())

		assert.ErrorIs(t, err, database.ErrTimeout)
	})
}

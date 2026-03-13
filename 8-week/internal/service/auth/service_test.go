package auth

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"

	domainUser "github.com/DaniilKalts/microservices-course-2023/8-week/internal/domain/user"
	userService "github.com/DaniilKalts/microservices-course-2023/8-week/internal/service/user"
	"github.com/DaniilKalts/microservices-course-2023/8-week/pkg/jwt"
)

// --- Helpers ---

func newTestService(t *testing.T) (*userService.UserServiceMock, *jwt.JWTManagerMock, Service) {
	t.Helper()
	userSvc := userService.NewUserServiceMock(t)
	// Mark methods not used by auth as optional.
	userSvc.ListMock.Optional()
	userSvc.GetMock.Optional()
	userSvc.UpdateMock.Optional()
	userSvc.DeleteMock.Optional()
	jwtMgr := jwt.NewJWTManagerMock(t)
	return userSvc, jwtMgr, NewService(userSvc, jwtMgr)
}

func stubTokenGeneration(jwtMgr *jwt.JWTManagerMock, access, refresh string) {
	jwtMgr.GenerateAccessTokenMock.Return(access, nil)
	jwtMgr.GenerateRefreshTokenMock.Return(refresh, nil)
	jwtMgr.AccessTokenTTLMock.Return(15 * time.Minute)
	jwtMgr.RefreshTokenTTLMock.Return(7 * 24 * time.Hour)
}

func skipJWT(jwtMgr *jwt.JWTManagerMock) {
	jwtMgr.GenerateAccessTokenMock.Optional()
	jwtMgr.GenerateRefreshTokenMock.Optional()
	jwtMgr.AccessTokenTTLMock.Optional()
	jwtMgr.RefreshTokenTTLMock.Optional()
}

var registerInput = RegisterInput{Name: "John", Email: "john@example.com", Password: "P@ssword123"}

// --- Tests ---

func TestRegister(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		userSvc, jwtMgr, svc := newTestService(t)

		userSvc.CreateMock.Set(func(_ context.Context, input domainUser.CreateInput) (string, error) {
			assert.Equal(t, "John", input.User.Name)
			assert.Equal(t, "john@example.com", input.User.Email)
			assert.Equal(t, domainUser.RoleUser, input.User.Role)
			assert.Equal(t, "P@ssword123", input.Password)
			return "user-id-123", nil
		})
		stubTokenGeneration(jwtMgr, "access-token", "refresh-token")

		user, tokens, err := svc.Register(context.Background(), registerInput)

		require.NoError(t, err)
		assert.Equal(t, "user-id-123", user.ID)
		assert.Equal(t, "John", user.Name)
		assert.Equal(t, "john@example.com", user.Email)
		assert.Equal(t, "access-token", tokens.AccessToken)
		assert.Equal(t, "refresh-token", tokens.RefreshToken)
		assert.Equal(t, int64(900), tokens.AccessTokenExpiresIn)
		assert.Equal(t, int64(604800), tokens.RefreshTokenExpiresIn)
	})

	t.Run("user creation fails", func(t *testing.T) {
		t.Parallel()
		userSvc, jwtMgr, svc := newTestService(t)
		skipJWT(jwtMgr)

		userSvc.CreateMock.Return("", errors.New("email already exists"))

		_, _, err := svc.Register(context.Background(), registerInput)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "email already exists")
		assert.Equal(t, uint64(0), jwtMgr.GenerateAccessTokenAfterCounter())
	})

	t.Run("empty user id", func(t *testing.T) {
		t.Parallel()
		userSvc, jwtMgr, svc := newTestService(t)
		skipJWT(jwtMgr)

		userSvc.CreateMock.Return("", nil)

		_, _, err := svc.Register(context.Background(), registerInput)

		require.ErrorIs(t, err, ErrIssueTokens)
	})

	t.Run("token generation fails", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name            string
			accessTokenErr  error
			refreshTokenErr error
		}{
			{"access token", errors.New("signing error"), nil},
			{"refresh token", nil, errors.New("signing error")},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				userSvc, jwtMgr, svc := newTestService(t)

				userSvc.CreateMock.Return("user-id-123", nil)
				if tt.accessTokenErr != nil {
					jwtMgr.GenerateAccessTokenMock.Return("", tt.accessTokenErr)
					jwtMgr.GenerateRefreshTokenMock.Optional()
				} else {
					jwtMgr.GenerateAccessTokenMock.Return("access-token", nil)
					jwtMgr.GenerateRefreshTokenMock.Return("", tt.refreshTokenErr)
				}
				jwtMgr.AccessTokenTTLMock.Optional().Return(15 * time.Minute)
				jwtMgr.RefreshTokenTTLMock.Optional().Return(7 * 24 * time.Hour)

				_, _, err := svc.Register(context.Background(), registerInput)

				require.ErrorIs(t, err, ErrIssueTokens)
			})
		}
	})
}

func TestLogin(t *testing.T) {
	t.Parallel()

	validHash, _ := bcrypt.GenerateFromPassword([]byte("P@ssword123"), bcrypt.DefaultCost)
	validCreds := &domainUser.Credentials{
		ID:           "user-id-123",
		PasswordHash: string(validHash),
		Role:         domainUser.RoleUser,
	}

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		userSvc, jwtMgr, svc := newTestService(t)

		userSvc.GetCredentialsByEmailMock.Set(func(_ context.Context, email string) (*domainUser.Credentials, error) {
			assert.Equal(t, "john@example.com", email)
			return validCreds, nil
		})
		stubTokenGeneration(jwtMgr, "access-token", "refresh-token")

		tokens, err := svc.Login(context.Background(), LoginInput{
			Email: "john@example.com", Password: "P@ssword123",
		})

		require.NoError(t, err)
		assert.Equal(t, "access-token", tokens.AccessToken)
		assert.Equal(t, "refresh-token", tokens.RefreshToken)
	})

	t.Run("user not found", func(t *testing.T) {
		t.Parallel()
		userSvc, jwtMgr, svc := newTestService(t)
		skipJWT(jwtMgr)

		userSvc.GetCredentialsByEmailMock.Return(nil, domainUser.ErrNotFound)

		tokens, err := svc.Login(context.Background(), LoginInput{
			Email: "nobody@example.com", Password: "P@ssword123",
		})

		require.ErrorIs(t, err, ErrInvalidCredentials)
		assert.Empty(t, tokens.AccessToken)
	})

	t.Run("wrong password", func(t *testing.T) {
		t.Parallel()
		userSvc, jwtMgr, svc := newTestService(t)
		skipJWT(jwtMgr)

		userSvc.GetCredentialsByEmailMock.Return(validCreds, nil)

		tokens, err := svc.Login(context.Background(), LoginInput{
			Email: "john@example.com", Password: "wrong-password",
		})

		require.ErrorIs(t, err, ErrInvalidCredentials)
		assert.Empty(t, tokens.AccessToken)
	})

	t.Run("internal error", func(t *testing.T) {
		t.Parallel()
		userSvc, jwtMgr, svc := newTestService(t)
		skipJWT(jwtMgr)

		userSvc.GetCredentialsByEmailMock.Return(nil, errors.New("connection refused"))

		tokens, err := svc.Login(context.Background(), LoginInput{
			Email: "john@example.com", Password: "P@ssword123",
		})

		require.ErrorIs(t, err, ErrAuthentication)
		assert.Empty(t, tokens.AccessToken)
	})

	t.Run("token generation fails", func(t *testing.T) {
		t.Parallel()
		userSvc, jwtMgr, svc := newTestService(t)

		userSvc.GetCredentialsByEmailMock.Return(validCreds, nil)
		jwtMgr.GenerateAccessTokenMock.Return("", errors.New("signing error"))
		jwtMgr.GenerateRefreshTokenMock.Optional()
		jwtMgr.AccessTokenTTLMock.Optional()
		jwtMgr.RefreshTokenTTLMock.Optional()

		tokens, err := svc.Login(context.Background(), LoginInput{
			Email: "john@example.com", Password: "P@ssword123",
		})

		require.ErrorIs(t, err, ErrIssueTokens)
		assert.Empty(t, tokens.AccessToken)
	})
}

func TestLogout(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		_, jwtMgr, svc := newTestService(t)

		jwtMgr.VerifyRefreshTokenMock.Set(func(token string) (*jwt.Claims, error) {
			assert.Equal(t, "valid-refresh-token", token)
			return &jwt.Claims{UserID: "user-id-123"}, nil
		})

		err := svc.Logout(context.Background(), LogoutInput{RefreshToken: "valid-refresh-token"})

		require.NoError(t, err)
	})

	t.Run("invalid token", func(t *testing.T) {
		t.Parallel()
		_, jwtMgr, svc := newTestService(t)

		jwtMgr.VerifyRefreshTokenMock.Return(nil, errors.New("token expired"))

		err := svc.Logout(context.Background(), LogoutInput{RefreshToken: "invalid-token"})

		require.ErrorIs(t, err, ErrInvalidRefreshToken)
	})
}

func TestRefresh(t *testing.T) {
	t.Parallel()

	validClaims := &jwt.Claims{UserID: "user-id-123", RoleID: int32(domainUser.RoleUser)}

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		_, jwtMgr, svc := newTestService(t)

		jwtMgr.VerifyRefreshTokenMock.Set(func(token string) (*jwt.Claims, error) {
			assert.Equal(t, "valid-refresh-token", token)
			return validClaims, nil
		})
		stubTokenGeneration(jwtMgr, "new-access-token", "new-refresh-token")

		tokens, err := svc.Refresh(context.Background(), RefreshInput{RefreshToken: "valid-refresh-token"})

		require.NoError(t, err)
		assert.Equal(t, "new-access-token", tokens.AccessToken)
		assert.Equal(t, "new-refresh-token", tokens.RefreshToken)
		assert.Equal(t, int64(900), tokens.AccessTokenExpiresIn)
		assert.Equal(t, int64(604800), tokens.RefreshTokenExpiresIn)
	})

	t.Run("invalid token", func(t *testing.T) {
		t.Parallel()
		_, jwtMgr, svc := newTestService(t)

		jwtMgr.VerifyRefreshTokenMock.Return(nil, errors.New("token expired"))

		tokens, err := svc.Refresh(context.Background(), RefreshInput{RefreshToken: "expired-token"})

		require.ErrorIs(t, err, ErrInvalidRefreshToken)
		assert.Empty(t, tokens.AccessToken)
	})

	t.Run("token generation fails", func(t *testing.T) {
		t.Parallel()
		_, jwtMgr, svc := newTestService(t)

		jwtMgr.VerifyRefreshTokenMock.Return(validClaims, nil)
		jwtMgr.GenerateAccessTokenMock.Return("", errors.New("signing error"))
		jwtMgr.GenerateRefreshTokenMock.Optional()
		jwtMgr.AccessTokenTTLMock.Optional()
		jwtMgr.RefreshTokenTTLMock.Optional()

		tokens, err := svc.Refresh(context.Background(), RefreshInput{RefreshToken: "valid-token"})

		require.ErrorIs(t, err, ErrIssueTokens)
		assert.Empty(t, tokens.AccessToken)
	})
}

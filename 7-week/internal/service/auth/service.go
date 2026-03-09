package auth

import (
	"context"
	"errors"
	"fmt"

	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"

	domainUser "github.com/DaniilKalts/microservices-course-2023/7-week/internal/domain/user"
	userService "github.com/DaniilKalts/microservices-course-2023/7-week/internal/service/user"
	"github.com/DaniilKalts/microservices-course-2023/7-week/pkg/jwt"
	"github.com/DaniilKalts/microservices-course-2023/7-week/pkg/tracing"
)

type Service interface {
	Register(ctx context.Context, input RegisterInput) (domainUser.User, TokenPair, error)
	Login(ctx context.Context, input LoginInput) (TokenPair, error)
	Logout(ctx context.Context, input LogoutInput) error
	Refresh(ctx context.Context, input RefreshInput) (TokenPair, error)
}

type service struct {
	userService userService.Service
	jwtManager  jwt.Manager
	logger      *zap.Logger
}

func NewService(
	userSvc userService.Service,
	jwtManager jwt.Manager,
	logger *zap.Logger,
) Service {
	return &service{
		userService: userSvc,
		jwtManager:  jwtManager,
		logger:      logger,
	}
}

const dummyPasswordHash = "$2a$10$7EqJtq98hPqEX7fNZaFWoO5m6jH9MuNoMNFcQJUO2cMjJ1ytY1/6W"

func (s *service) Register(ctx context.Context, input RegisterInput) (domainUser.User, TokenPair, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "service.auth.Register")
	defer span.Finish()

	userID, err := s.userService.Create(ctx, userService.CreateInput{
		User: &domainUser.User{
			Name:  input.Name,
			Email: input.Email,
			Role:  domainUser.RoleUser,
		},
		Password: input.Password,
	})
	if err != nil {
		return domainUser.User{}, TokenPair{}, err
	}

	tokenPair, err := s.generateTokenPair(userID, int32(domainUser.RoleUser))
	if err != nil {
		tracing.LogError(s.logger, span, "failed to issue tokens after registration", err, zap.String("user_id", userID))
		return domainUser.User{}, TokenPair{}, fmt.Errorf("%w: %v", ErrIssueTokens, err)
	}

	s.logger.Info("user registered", zap.String("user_id", userID), zap.String("email", input.Email))

	return domainUser.User{ID: userID, Name: input.Name, Email: input.Email}, tokenPair, nil
}

func (s *service) Login(ctx context.Context, input LoginInput) (TokenPair, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "service.auth.Login")
	defer span.Finish()

	credentials, err := s.userService.GetCredentialsByEmail(ctx, input.Email)
	if err != nil {
		_ = bcrypt.CompareHashAndPassword([]byte(dummyPasswordHash), []byte(input.Password))
		if errors.Is(err, domainUser.ErrNotFound) {
			tracing.LogWarn(s.logger, span, "login attempt for non-existent email", ErrInvalidCredentials, zap.String("email", input.Email))
			return TokenPair{}, ErrInvalidCredentials
		}

		return TokenPair{}, fmt.Errorf("%w: %v", ErrAuthentication, err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(credentials.PasswordHash), []byte(input.Password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			tracing.LogWarn(s.logger, span, "invalid password attempt", ErrInvalidCredentials, zap.String("email", input.Email))
			return TokenPair{}, ErrInvalidCredentials
		}

		tracing.LogError(s.logger, span, "failed to compare password hash", err, zap.String("email", input.Email))
		return TokenPair{}, fmt.Errorf("%w: %v", ErrAuthentication, err)
	}

	tokenPair, err := s.generateTokenPair(credentials.ID, int32(credentials.Role))
	if err != nil {
		tracing.LogError(s.logger, span, "failed to issue tokens during login", err, zap.String("user_id", credentials.ID))
		return TokenPair{}, fmt.Errorf("%w: %v", ErrIssueTokens, err)
	}

	s.logger.Info("user logged in", zap.String("user_id", credentials.ID), zap.String("email", input.Email))

	return tokenPair, nil
}

func (s *service) Logout(ctx context.Context, input LogoutInput) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "service.auth.Logout")
	defer span.Finish()

	claims, err := s.jwtManager.VerifyRefreshToken(input.RefreshToken)
	if err != nil {
		tracing.LogWarn(s.logger, span, "invalid refresh token on logout", err)
		return fmt.Errorf("%w: %v", ErrInvalidRefreshToken, err)
	}

	s.logger.Info("user logged out", zap.String("user_id", claims.UserID))

	return nil
}

func (s *service) Refresh(ctx context.Context, input RefreshInput) (TokenPair, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "service.auth.Refresh")
	defer span.Finish()

	claims, err := s.jwtManager.VerifyRefreshToken(input.RefreshToken)
	if err != nil {
		tracing.LogWarn(s.logger, span, "invalid refresh token on refresh", err)
		return TokenPair{}, fmt.Errorf("%w: %v", ErrInvalidRefreshToken, err)
	}

	tokenPair, err := s.generateTokenPair(claims.UserID, claims.RoleID)
	if err != nil {
		tracing.LogError(s.logger, span, "failed to issue tokens during refresh", err, zap.String("user_id", claims.UserID))
		return TokenPair{}, fmt.Errorf("%w: %v", ErrIssueTokens, err)
	}

	s.logger.Info("tokens refreshed", zap.String("user_id", claims.UserID))

	return tokenPair, nil
}

func (s *service) generateTokenPair(userID string, roleID int32) (TokenPair, error) {
	if userID == "" {
		return TokenPair{}, ErrUserIDEmpty
	}

	claims := jwt.Claims{
		UserID: userID,
		RoleID: roleID,
	}

	accessToken, err := s.jwtManager.GenerateAccessToken(claims)
	if err != nil {
		return TokenPair{}, err
	}

	refreshToken, err := s.jwtManager.GenerateRefreshToken(claims)
	if err != nil {
		return TokenPair{}, err
	}

	accessTokenExpiresIn := int64(s.jwtManager.AccessTokenTTL().Seconds())
	refreshTokenExpiresIn := int64(s.jwtManager.RefreshTokenTTL().Seconds())

	return TokenPair{
		AccessToken:           accessToken,
		RefreshToken:          refreshToken,
		AccessTokenExpiresIn:  accessTokenExpiresIn,
		RefreshTokenExpiresIn: refreshTokenExpiresIn,
	}, nil
}

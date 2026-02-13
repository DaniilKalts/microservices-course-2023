package user

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/DaniilKalts/microservices-course-2023/4-week/internal/models"
	"github.com/DaniilKalts/microservices-course-2023/4-week/internal/service/user/mocks"
)

func TestCreate_ValidationScenarios(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		name            string
		user            *models.User
		password        string
		passwordConfirm string
		errText         string
		repoCalls       uint64
	}{
		{
			name: "TestCreate_PasswordMismatch",
			user: &models.User{
				Name:  "John",
				Email: "john@example.com",
			},
			password:        "P@ssword123",
			passwordConfirm: "P@ssword321",
			errText:         "Passwords don't match",
			repoCalls:       0,
		},
		{
			name: "TestCreate_WeakPassword",
			user: &models.User{
				Name:  "John",
				Email: "john@example.com",
			},
			password:        "12345",
			passwordConfirm: "12345",
			repoCalls:       0,
		},
		{
			name: "TestCreate_EmptyName",
			user: &models.User{
				Name:  "",
				Email: "john@example.com",
			},
			password:        "P@ssword123",
			passwordConfirm: "P@ssword123",
			repoCalls:       0,
		},
		{
			name: "TestCreate_EmptyEmail",
			user: &models.User{
				Name:  "John",
				Email: "",
			},
			password:        "P@ssword123",
			passwordConfirm: "P@ssword123",
			repoCalls:       0,
		},
		{
			name: "TestCreate_InvalidEmailFormat",
			user: &models.User{
				Name:  "John",
				Email: "not-an-email",
			},
			password:        "P@ssword123",
			passwordConfirm: "P@ssword123",
			repoCalls:       0,
		},
		{
			name: "TestCreate_BcryptTooLongPassword",
			user: &models.User{
				Name:  "John",
				Email: "john@example.com",
			},
			password:        strings.Repeat("a", 73),
			passwordConfirm: strings.Repeat("a", 73),
			repoCalls:       0,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repoMock := mocks.NewUserRepositoryMock(t)
			repoMock.CreateMock.Optional().Set(func(_ context.Context, _ *models.User, _ string) (string, error) {
				return "", errors.New("repository should not be called")
			})

			svc := NewService(repoMock, nil)
			gotID, err := svc.Create(ctx, tt.user, tt.password, tt.passwordConfirm)

			require.Error(t, err)
			if tt.errText != "" {
				require.EqualError(t, err, tt.errText)
			}
			require.Empty(t, gotID)
			require.Equal(t, tt.repoCalls, repoMock.CreateBeforeCounter())
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
		{
			name:       "TestCreate_ReturnsUserID",
			repoResult: "repo-returned-id",
			wantID:     "repo-returned-id",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repoMock := mocks.NewUserRepositoryMock(t)
			repoMock.CreateMock.Set(func(_ context.Context, _ *models.User, _ string) (string, error) {
				return tt.repoResult, tt.repoErr
			})

			svc := NewService(repoMock, nil)
			gotID, err := svc.Create(ctx, &models.User{Name: "John", Email: "john@example.com"}, "P@ssword123", "P@ssword123")

			if tt.wantErr != "" {
				require.EqualError(t, err, tt.wantErr)
				require.Empty(t, gotID)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.wantID, gotID)
		})
	}
}

package user

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	domainUser "github.com/DaniilKalts/microservices-course-2023/4-week/internal/domain/user"
)

func TestCreate_ValidationScenarios(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		name            string
		password        string
		passwordConfirm string
		errText         string
		repoCalls       int
	}{
		{
			name:            "TestCreate_PasswordMismatch",
			password:        "P@ssword123",
			passwordConfirm: "P@ssword321",
			errText:         "Passwords don't match",
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
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := &repoStub{}
			repo.createFn = func(_ context.Context, _ *domainUser.Entity, _ string) (string, error) {
				return "", errors.New("repository should not be called")
			}

			svc := NewService(repo)
			gotID, err := svc.Create(ctx, &domainUser.Entity{Name: "John", Email: "john@example.com"}, tt.password, tt.passwordConfirm)

			require.Error(t, err)
			if tt.errText != "" {
				require.EqualError(t, err, tt.errText)
			}
			require.Empty(t, gotID)
			require.Equal(t, tt.repoCalls, repo.createCalls)
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
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := &repoStub{}
			repo.createFn = func(_ context.Context, user *domainUser.Entity, passwordHash string) (string, error) {
				require.NotEmpty(t, user.ID)
				require.NotEmpty(t, passwordHash)
				return tt.repoResult, tt.repoErr
			}

			svc := NewService(repo)
			gotID, err := svc.Create(ctx, &domainUser.Entity{Name: "John", Email: "john@example.com"}, "P@ssword123", "P@ssword123")

			if tt.wantErr != "" {
				require.EqualError(t, err, tt.wantErr)
				require.Empty(t, gotID)
				require.Equal(t, 1, repo.createCalls)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.wantID, gotID)
			require.Equal(t, 1, repo.createCalls)
		})
	}
}

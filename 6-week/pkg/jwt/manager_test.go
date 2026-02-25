package jwt

import (
	"crypto/rand"
	"crypto/rsa"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestManager_TokenPairHasUniqueJTI(t *testing.T) {
	t.Parallel()

	m := newTestManager(t)
	claims := Claims{UserID: "user-id", RoleID: 1}

	accessToken, err := m.GenerateAccessToken(claims)
	require.NoError(t, err)

	refreshToken, err := m.GenerateRefreshToken(claims)
	require.NoError(t, err)

	accessClaims, err := m.VerifyAccessToken(accessToken)
	require.NoError(t, err)

	refreshClaims, err := m.VerifyRefreshToken(refreshToken)
	require.NoError(t, err)

	require.NotEmpty(t, accessClaims.ID)
	require.NotEmpty(t, refreshClaims.ID)
	require.NotEqual(t, accessClaims.ID, refreshClaims.ID)

	require.Equal(t, claims.UserID, accessClaims.UserID)
	require.Equal(t, claims.UserID, refreshClaims.UserID)
}

func newTestManager(t *testing.T) Manager {
	t.Helper()

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	m, err := NewManager(privateKey, &privateKey.PublicKey, Config{})
	require.NoError(t, err)

	return m
}

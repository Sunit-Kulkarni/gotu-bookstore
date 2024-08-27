package auth

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJWTAuth(t *testing.T) {
	// Generate a valid token
	token, err := GenerateToken(123)
	require.NoError(t, err)

	// Test valid token
	uid, err := JWTAuth(context.Background(), token)
	require.NoError(t, err)
	assert.Equal(t, "123", string(uid))

	// Test invalid token
	_, err = JWTAuth(context.Background(), "invalid-token")
	assert.Error(t, err)
}

func TestGenerateToken(t *testing.T) {
	token, err := GenerateToken(123)
	require.NoError(t, err)
	assert.NotEmpty(t, token)

	// Verify the generated token
	uid, err := JWTAuth(context.Background(), token)
	require.NoError(t, err)
	assert.Equal(t, "123", string(uid))
}

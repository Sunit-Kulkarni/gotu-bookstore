package account

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateAccount(t *testing.T) {
	ctx := context.Background()

	// Test successful account creation
	params := &CreateAccountParams{
		Email:    "test@example.com",
		Password: "password123",
	}
	resp, err := CreateAccount(ctx, params)
	require.NoError(t, err)
	assert.NotEmpty(t, resp.ID)
	assert.Equal(t, params.Email, resp.Email)

	// Test duplicate email
	_, err = CreateAccount(ctx, params)
	assert.Error(t, err)

	// Test invalid email
	invalidParams := &CreateAccountParams{
		Email:    "invalid-email",
		Password: "password123",
	}
	_, err = CreateAccount(ctx, invalidParams)
	assert.Error(t, err)

	// Test short password
	shortPassParams := &CreateAccountParams{
		Email:    "another@example.com",
		Password: "short",
	}
	_, err = CreateAccount(ctx, shortPassParams)
	assert.Error(t, err)
}

func TestLogin(t *testing.T) {
	ctx := context.Background()

	// Create a test account
	createParams := &CreateAccountParams{
		Email:    "login-test@example.com",
		Password: "password123",
	}
	_, err := CreateAccount(ctx, createParams)
	require.NoError(t, err)

	// Test successful login
	loginParams := &LoginParams{
		Email:    "login-test@example.com",
		Password: "password123",
	}
	resp, err := Login(ctx, loginParams)
	require.NoError(t, err)
	assert.NotEmpty(t, resp.ID)
	assert.Equal(t, loginParams.Email, resp.Email)
	assert.NotEmpty(t, resp.Token)

	// Test login with incorrect password
	wrongPassParams := &LoginParams{
		Email:    "login-test@example.com",
		Password: "wrongpassword",
	}
	_, err = Login(ctx, wrongPassParams)
	assert.Error(t, err)

	// Test login with non-existent email
	nonExistentParams := &LoginParams{
		Email:    "nonexistent@example.com",
		Password: "password123",
	}
	_, err = Login(ctx, nonExistentParams)
	assert.Error(t, err)
}

package account

import (
	"context"
	auth2 "encore.app/auth"
	"fmt"
	"net/mail"

	"encore.dev/beta/errs"
	"golang.org/x/crypto/bcrypt"
)

//encore:api public method=POST path=/users
func CreateAccount(ctx context.Context, params *CreateAccountParams) (*CreateAccountResponse, error) {
	if _, err := mail.ParseAddress(params.Email); err != nil {
		return nil, errs.Wrap(err, "invalid email format")
	}

	if len(params.Password) < 8 {
		return nil, fmt.Errorf("password must be at least 8 characters long")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %v", err)
	}

	var exists bool
	err = bookstoredb.QueryRow(ctx, `
        SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)
    `, params.Email).Scan(&exists)
	if err != nil {
		return nil, fmt.Errorf("failed to check email existence: %v", err)
	}
	if exists {
		return nil, fmt.Errorf("email already in use")
	}

	var user User
	err = bookstoredb.QueryRow(ctx, `
        INSERT INTO users (email, password)
        VALUES ($1, $2)
        RETURNING id, email
    `, params.Email, string(hashedPassword)).Scan(&user.ID, &user.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %v", err)
	}

	return &CreateAccountResponse{
		ID:    user.ID,
		Email: user.Email,
	}, nil
}

//encore:api public method=POST path=/login
func Login(ctx context.Context, params *LoginParams) (*LoginResponse, error) {
	var user User
	err := bookstoredb.QueryRow(ctx, `
        SELECT id, email, password FROM users WHERE email = $1
    `, params.Email).Scan(&user.ID, &user.Email, &user.Password)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(params.Password)); err != nil {
		return nil, fmt.Errorf("invalid password")
	}

	token, err := auth2.GenerateToken(user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %v", err)
	}

	return &LoginResponse{
		ID:    user.ID,
		Email: user.Email,
		Token: token,
	}, nil
}

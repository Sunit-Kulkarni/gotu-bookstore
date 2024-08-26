package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"encore.dev/beta/auth"
	"encore.dev/beta/errs"
	"github.com/golang-jwt/jwt"
)

type Claims struct {
	UserID int `json:"user_id"`
	jwt.StandardClaims
}

var secrets struct {
	JWTKey string
}

var jwtKey = []byte(secrets.JWTKey)

//encore:authhandler
func JWTAuth(ctx context.Context, token string) (auth.UID, error) {
	claims := &Claims{}
	tkn, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtKey, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrSignatureInvalid) {
			return "", errs.WrapCode(err, errs.PermissionDenied, "invalid token signature")
		}
		return "", errs.WrapCode(err, errs.PermissionDenied, "invalid token")
	}

	if !tkn.Valid {
		return "", errs.WrapCode(err, errs.PermissionDenied, "invalid token")
	}

	return auth.UID(fmt.Sprintf("%d", claims.UserID)), nil
}

func GenerateToken(userID int) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", fmt.Errorf("failed to create token: %v", err)
	}

	return tokenString, nil
}

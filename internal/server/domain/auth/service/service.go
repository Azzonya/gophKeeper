// Package service handles authentication via JWT, including token creation,
// validation, and user ID extraction from gRPC context metadata.
package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"google.golang.org/grpc/metadata"
	"gophKeeper/internal/server/domain/users/model"
	"time"
)

const (
	defaultJWTCookieExpiration = 24 * time.Hour
)

// Claims represents the custom claims used for JWT tokens, including the user ID (UID)
// and standard JWT registered claims like expiration time.
type Claims struct {
	jwt.RegisteredClaims
	UID string
}

// Auth handles authentication-related operations, such as creating and validating JWT tokens.
type Auth struct {
	JwtSecret string
}

// New creates a new Auth instance with the given JWT secret.
func New(jwtSecret string) *Auth {
	return &Auth{JwtSecret: jwtSecret}
}

// GetUserIDFromContext extracts the user ID from the JWT token found in the incoming gRPC context metadata.
// It returns the user ID if the token is valid or an error if the token is invalid or missing.
func (a *Auth) GetUserIDFromContext(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", errors.New("missing metadata in context")
	}

	mdToken := md["token"]
	if len(mdToken) == 0 {
		return "", errors.New("missing cookies in metadata")
	}

	var jwtToken string
	for _, cookieStr := range mdToken {
		if cookieStr[:7] == "Bearer " {
			jwtToken = cookieStr[7:]
			break
		}
	}

	if jwtToken == "" {
		return "", errors.New("jwt cookie not found")
	}

	token, err := jwt.ParseWithClaims(jwtToken, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(a.JwtSecret), nil
	})
	if err != nil {
		return "", errors.New("invalid jwt token")
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims.UID, nil
	}

	return "", errors.New("invalid token")
}

// CreateToken generates a signed JWT token for a given user, based on their user ID.
func (a *Auth) CreateToken(u *model.User) (string, error) {
	token, err := a.NewToken(u)
	if err != nil {
		return "", fmt.Errorf("cannot create auth token: %w", err)
	}
	return token, nil
}

// NewToken creates a new JWT token with an expiration time and includes the user ID (UID) in the claims.
// The token is signed using the provided JWT secret.
func (a *Auth) NewToken(u *model.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(defaultJWTCookieExpiration)),
		},
		UID: u.UserID,
	})

	signedToken, err := token.SignedString([]byte(a.JwtSecret))
	if err != nil {
		return "", fmt.Errorf("cannot sign jwt token: %w", err)
	}
	return signedToken, nil
}

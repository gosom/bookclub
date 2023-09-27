package bookclub

import (
	"context"
	"time"
)

//go:generate mockgen -package mocks -destination ./mocks/mock_auth_uc.go . AuthUseCases,JWTProvider

type LoginParams struct {
	Email    string
	Password string
}

type AuthUseCases interface {
	Login(ctx context.Context, params LoginParams) (string, string, error)
}

type JWTGenerateParams struct {
	UserID  string
	TTL     time.Duration
	Refresh bool
}

type JWTClaims struct {
	Subject   string
	Refresh   bool
	IssuedAt  time.Time
	ExpiresAt time.Time
}

type JWTProvider interface {
	GenerateToken(ctx context.Context, params JWTGenerateParams) (string, error)
	GenerateRefreshToken(ctx context.Context, params JWTGenerateParams) (string, error)
}

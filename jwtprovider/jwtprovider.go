package jwtprovider

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/gosom/bookclub"
)

var _ bookclub.JWTProvider = (*jwtProvider)(nil)

type jwtProvider struct {
	issuer string
	secret []byte
}

func New(secret, issuer string) bookclub.JWTProvider {
	ans := &jwtProvider{
		issuer: issuer,
		secret: []byte(secret),
	}

	return ans
}

func (o *jwtProvider) GenerateToken(ctx context.Context, params bookclub.JWTGenerateParams) (string, error) {
	if params.Refresh {
		return "", errors.New("cannot generate access token with Refresh true")
	}

	return o.generateToken(ctx, params)
}

func (o *jwtProvider) GenerateRefreshToken(ctx context.Context, params bookclub.JWTGenerateParams) (string, error) {
	if !params.Refresh {
		return "", errors.New("cannot generate refresh token with Refresh false")
	}

	return o.generateToken(ctx, params)
}

func (o *jwtProvider) generateToken(ctx context.Context, params bookclub.JWTGenerateParams) (string, error) {
	if params.UserID == "" {
		return "", errors.New("cannot generate token with empty ID")
	}

	if params.TTL <= 0 {
		return "", errors.New("cannot generate token with negative TTL")
	}

	now := time.Now().UTC()

	claims := bookclub.JWTClaims{
		Subject:   params.UserID,
		Refresh:   params.Refresh,
		IssuedAt:  now,
		ExpiresAt: now.Add(params.TTL),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, customClaims{
		Refresh: params.Refresh,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   claims.Subject,
			Issuer:    o.issuer,
			IssuedAt:  jwt.NewNumericDate(claims.IssuedAt),
			ExpiresAt: jwt.NewNumericDate(claims.ExpiresAt),
		},
	})

	return token.SignedString(o.secret)
}

type customClaims struct {
	Refresh bool
	jwt.RegisteredClaims
}

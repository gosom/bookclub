package authuc

import (
	"context"
	"errors"
	"time"

	"github.com/gosom/bookclub"
)

var _ bookclub.AuthUseCases = (*authUseCases)(nil)

type authUseCases struct {
	store       bookclub.Storage
	jwtProvider bookclub.JWTProvider
}

func NewAuthUseCases(
	store bookclub.Storage,
	jwtProvider bookclub.JWTProvider,
) bookclub.AuthUseCases {
	ans := &authUseCases{
		store:       store,
		jwtProvider: jwtProvider,
	}

	return ans
}

func (o *authUseCases) Login(ctx context.Context, params bookclub.LoginParams) (string, string, error) {
	user, err := o.store.GetUserByEmail(ctx, bookclub.Email(params.Email))
	if err != nil {
		if errors.Is(err, bookclub.ErrNotFound) {
			return "", "", bookclub.ErrInvalidCredentials
		}

		return "", "", bookclub.ErrInternalError
	}

	if err := user.ComparePassword(params.Password); err != nil {
		return "", "", bookclub.ErrInvalidCredentials
	}

	accessToken, err := o.jwtProvider.GenerateToken(ctx, bookclub.JWTGenerateParams{
		UserID: user.ID.String(),
		TTL:    15 * time.Minute,
	})

	if err != nil {
		return "", "", bookclub.ErrInternalError
	}

	refreshToken, err := o.jwtProvider.GenerateRefreshToken(ctx, bookclub.JWTGenerateParams{
		Refresh: true,
		UserID:  user.ID.String(),
		TTL:     1 * time.Hour,
	})
	if err != nil {
		return "", "", bookclub.ErrInternalError
	}

	return accessToken, refreshToken, nil
}

func (o *authUseCases) Refresh(ctx context.Context, refreshToken string) (string, string, error) {
	claims, err := o.jwtProvider.ValidateToken(ctx, refreshToken)
	if err != nil {
		return "", "", bookclub.ErrInvalidCredentials
	}

	if !claims.Refresh {
		return "", "", bookclub.ErrInvalidCredentials
	}

	accessToken, err := o.jwtProvider.GenerateToken(ctx, bookclub.JWTGenerateParams{
		UserID: claims.Subject,
		TTL:    15 * time.Minute,
	})

	if err != nil {
		return "", "", bookclub.ErrInternalError
	}

	refreshToken, err = o.jwtProvider.GenerateRefreshToken(ctx, bookclub.JWTGenerateParams{
		Refresh: true,
		UserID:  claims.Subject,
		TTL:     1 * time.Hour,
	})
	if err != nil {
		return "", "", bookclub.ErrInternalError
	}

	return accessToken, refreshToken, nil
}

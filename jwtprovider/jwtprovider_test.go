package jwtprovider_test

import (
	"context"
	"testing"
	"time"

	"github.com/gosom/bookclub"
	"github.com/gosom/bookclub/jwtprovider"
	"github.com/stretchr/testify/require"
)

func TestGenerateToken(t *testing.T) {
	prov := jwtprovider.New("secret", "issuer")

	token, err := prov.GenerateToken(context.Background(), bookclub.JWTGenerateParams{
		UserID: "user_id",
		TTL:    1 * time.Minute,
	})
	require.NoError(t, err)

	// this should be a valid JWT token

	claims, err := prov.ValidateToken(context.Background(), token)
	require.NoError(t, err)

	require.Equal(t, "user_id", claims.Subject)
	require.False(t, claims.Refresh)
}

func TestGenerateRefreshToken(t *testing.T) {
	prov := jwtprovider.New("secret", "issuer")
	refresh, err := prov.GenerateRefreshToken(context.Background(), bookclub.JWTGenerateParams{
		Refresh: true,
		UserID:  "user_id",
		TTL:     1 * time.Hour,
	})
	require.NoError(t, err)

	// this should be a valid JWT token

	claims, err := prov.ValidateToken(context.Background(), refresh)
	require.NoError(t, err)

	require.Equal(t, "user_id", claims.Subject)
	require.True(t, claims.Refresh)
}

func Test_ValidateToken(t *testing.T) {
	prov := jwtprovider.New("secret", "issuer")

	t.Run("should return ErrInvalidCredentials if token is invalid", func(t *testing.T) {
		_, err := prov.ValidateToken(context.Background(), "invalid_token")
		require.ErrorIs(t, err, bookclub.ErrInvalidCredentials)
	})

	t.Run("should return ErrInvalidCredentials if token is expired", func(t *testing.T) {
		token, err := prov.GenerateToken(context.Background(), bookclub.JWTGenerateParams{
			UserID: "user_id",
			TTL:    1 * time.Microsecond,
		})

		require.NoError(t, err)

		time.Sleep(1 * time.Millisecond)

		_, err = prov.ValidateToken(context.Background(), token)
		require.ErrorIs(t, err, bookclub.ErrInvalidCredentials)
	})

	t.Run("should return ErrInvalidCredentials if token is not valid for this issuer", func(t *testing.T) {
		prov2 := jwtprovider.New("secret", "issuer2")

		token, err := prov2.GenerateToken(context.Background(), bookclub.JWTGenerateParams{
			UserID: "user_id",
			TTL:    1 * time.Minute,
		})

		require.NoError(t, err)

		_, err = prov.ValidateToken(context.Background(), token)
		require.ErrorIs(t, err, bookclub.ErrInvalidCredentials)
	})

	t.Run("should return ErrInvalidCredentials if token is not valid for this secret", func(t *testing.T) {
		prov2 := jwtprovider.New("secret2", "issuer")

		token, err := prov2.GenerateToken(context.Background(), bookclub.JWTGenerateParams{
			UserID: "user_id",
			TTL:    1 * time.Minute,
		})

		require.NoError(t, err)

		_, err = prov.ValidateToken(context.Background(), token)
		require.ErrorIs(t, err, bookclub.ErrInvalidCredentials)
	})
}

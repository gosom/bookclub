package authuc_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/gosom/bookclub"
	"github.com/gosom/bookclub/mocks"
	"github.com/gosom/bookclub/usecases/authuc"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func Test_NewAuthUseCases(t *testing.T) {
	mctrl := gomock.NewController(t)

	store := mocks.NewMockStorage(mctrl)
	jwtProv := mocks.NewMockJWTProvider(mctrl)

	uc := authuc.NewAuthUseCases(store, jwtProv)
	require.NotNil(t, uc)
}

func Test_Login(t *testing.T) {
	t.Run("should return ErrInvalidCredentials if user not found", func(t *testing.T) {
		mctrl := gomock.NewController(t)

		store := mocks.NewMockStorage(mctrl)
		store.EXPECT().GetUserByEmail(gomock.Any(), gomock.Any()).
			Return(bookclub.User{}, bookclub.ErrNotFound)

		uc := authuc.NewAuthUseCases(store, nil)

		_, _, err := uc.Login(context.Background(), bookclub.LoginParams{
			Email:    "j@doe.com",
			Password: "12345678",
		})

		require.ErrorIs(t, err, bookclub.ErrInvalidCredentials)
	})

	t.Run("should return ErrInvalidCredentials if password mismatch", func(t *testing.T) {
		mctrl := gomock.NewController(t)

		store := mocks.NewMockStorage(mctrl)
		store.EXPECT().GetUserByEmail(gomock.Any(), gomock.Any()).
			Return(bookclub.User{
				Email: "j@doe.com",
			}, nil)

		uc := authuc.NewAuthUseCases(store, nil)

		_, _, err := uc.Login(context.Background(), bookclub.LoginParams{
			Email:    "j@doe.com",
			Password: "12345678",
		})

		require.ErrorIs(t, err, bookclub.ErrInvalidCredentials)
	})

	t.Run("should return internal error if store returns error", func(t *testing.T) {
		mctrl := gomock.NewController(t)

		store := mocks.NewMockStorage(mctrl)
		store.EXPECT().GetUserByEmail(gomock.Any(), gomock.Any()).
			Return(bookclub.User{}, bookclub.ErrInternalError)

		uc := authuc.NewAuthUseCases(store, nil)

		_, _, err := uc.Login(context.Background(), bookclub.LoginParams{
			Email:    "j@doe.com",
			Password: "12345678",
		})

		require.ErrorIs(t, err, bookclub.ErrInternalError)
	})

	t.Run("when a valid password is provided", func(t *testing.T) {
		bcryptedPassword, err := bookclub.NewPassword("123Ab!4%8!022m")
		require.NoError(t, err)

		t.Run("should return a valid access token and refresh token", func(t *testing.T) {
			mctrl := gomock.NewController(t)

			user := bookclub.User{
				ID:    uuid.New(),
				Email: "j@doe.com",
			}
			user.SetPassword(bcryptedPassword)

			store := mocks.NewMockStorage(mctrl)
			store.EXPECT().GetUserByEmail(gomock.Any(), gomock.Any()).
				Return(user, nil)

			jwtProv := mocks.NewMockJWTProvider(mctrl)
			jwtProv.EXPECT().GenerateToken(gomock.Any(), bookclub.JWTGenerateParams{
				UserID: user.ID.String(),
				TTL:    15 * time.Minute,
			}).
				Return("access_token", nil)
			jwtProv.EXPECT().GenerateRefreshToken(gomock.Any(), bookclub.JWTGenerateParams{
				Refresh: true,
				UserID:  user.ID.String(),
				TTL:     1 * time.Hour,
			}).
				Return("refresh_token", nil)

			uc := authuc.NewAuthUseCases(store, jwtProv)

			accessToken, refreshToken, err := uc.Login(
				context.Background(),
				bookclub.LoginParams{
					Email:    "j@doe.com",
					Password: "123Ab!4%8!022m",
				})

			require.NoError(t, err)

			require.Equal(t, "access_token", accessToken)
			require.Equal(t, "refresh_token", refreshToken)
		})

		t.Run("should return internal error if jwt provider returns error", func(t *testing.T) {
			mctrl := gomock.NewController(t)

			user := bookclub.User{
				ID:    uuid.New(),
				Email: "j@doe.com",
			}
			user.SetPassword(bcryptedPassword)

			store := mocks.NewMockStorage(mctrl)
			store.EXPECT().GetUserByEmail(gomock.Any(), gomock.Any()).
				Return(user, nil)

			jwtProv := mocks.NewMockJWTProvider(mctrl)
			jwtProv.EXPECT().GenerateToken(gomock.Any(), bookclub.JWTGenerateParams{
				UserID: user.ID.String(),
				TTL:    15 * time.Minute,
			}).
				Return("", errors.New("jwt provider error"))

			uc := authuc.NewAuthUseCases(store, jwtProv)

			_, _, err := uc.Login(
				context.Background(),
				bookclub.LoginParams{
					Email:    "j@doe.com",
					Password: "123Ab!4%8!022m",
				})

			require.ErrorIs(t, err, bookclub.ErrInternalError)
		})
	})
}

func Test_Refresh(t *testing.T) {
	t.Run("should return ErrInvalidCredentials if refresh token is invalid", func(t *testing.T) {
		mctrl := gomock.NewController(t)

		jwtProv := mocks.NewMockJWTProvider(mctrl)
		jwtProv.EXPECT().ValidateToken(gomock.Any(), gomock.Any()).
			Return(bookclub.JWTClaims{}, bookclub.ErrInvalidCredentials)

		uc := authuc.NewAuthUseCases(nil, jwtProv)

		_, _, err := uc.Refresh(
			context.Background(), "refresh_token",
		)

		require.ErrorIs(t, err, bookclub.ErrInvalidCredentials)
	})

	t.Run("should return a new access and refresh token", func(t *testing.T) {
		userID := uuid.New()
		mctrl := gomock.NewController(t)

		jwtProv := mocks.NewMockJWTProvider(mctrl)
		jwtProv.EXPECT().ValidateToken(gomock.Any(), gomock.Any()).
			Return(bookclub.JWTClaims{
				Subject: userID.String(),
				Refresh: true,
			}, nil)

		jwtProv.EXPECT().GenerateToken(gomock.Any(), bookclub.JWTGenerateParams{
			UserID: userID.String(),
			TTL:    15 * time.Minute,
		}).
			Return("access_token", nil)
		jwtProv.EXPECT().GenerateRefreshToken(gomock.Any(), bookclub.JWTGenerateParams{
			Refresh: true,
			UserID:  userID.String(),
			TTL:     1 * time.Hour,
		}).
			Return("refresh_token", nil)

		uc := authuc.NewAuthUseCases(nil, jwtProv)

		accessToken, refreshToken, err := uc.Refresh(
			context.Background(), "refresh_token",
		)

		require.NoError(t, err)

		require.Equal(t, "access_token", accessToken)
		require.Equal(t, "refresh_token", refreshToken)
	})
}

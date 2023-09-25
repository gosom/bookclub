package useruc_test

import (
	"context"
	"testing"

	"go.uber.org/mock/gomock"

	"github.com/gosom/bookclub"
	"github.com/gosom/bookclub/mocks"
	"github.com/gosom/bookclub/usecases/useruc"
	"github.com/stretchr/testify/require"
)

func Test_NewUserUseCases(t *testing.T) {
	mctrl := gomock.NewController(t)

	store := mocks.NewMockStorage(mctrl)

	uc := useruc.NewUserUseCases(store)
	require.NotNil(t, uc)
}

func Test_Register(t *testing.T) {
	t.Run("with invalid email", func(t *testing.T) {
		mctrl := gomock.NewController(t)

		store := mocks.NewMockStorage(mctrl)

		uc := useruc.NewUserUseCases(store)

		_, err := uc.Register(context.Background(), bookclub.RegisterParams{
			Email:    "not an email",
			Password: "password",
		})

		require.Error(t, err)
		require.Equal(t, bookclub.ErrInvalidEmail, err)
	})

	t.Run("with invalid password", func(t *testing.T) {
		mctrl := gomock.NewController(t)

		store := mocks.NewMockStorage(mctrl)

		uc := useruc.NewUserUseCases(store)

		_, err := uc.Register(context.Background(), bookclub.RegisterParams{
			Email:    "john@doe.com",
			Password: "short",
		})

		require.Error(t, err)
		require.Equal(t, bookclub.ErrInvalidPassword, err)
	})

	t.Run("when database returns an internal error", func(t *testing.T) {
		mctrl := gomock.NewController(t)

		store := mocks.NewMockStorage(mctrl)
		store.EXPECT().CreateUser(gomock.Any(), gomock.Any(), gomock.Any()).Return(bookclub.User{}, bookclub.ErrInternalError)

		uc := useruc.NewUserUseCases(store)

		_, err := uc.Register(context.Background(), bookclub.RegisterParams{
			Email:    "john@doe.com",
			Password: "123abc!A#4%",
		})

		require.Error(t, err)
		require.Equal(t, bookclub.ErrInternalError, err)
	})

	t.Run("when database returns an already exists error", func(t *testing.T) {
		mctrl := gomock.NewController(t)

		store := mocks.NewMockStorage(mctrl)
		store.EXPECT().CreateUser(gomock.Any(), gomock.Any(), gomock.Any()).Return(bookclub.User{}, bookclub.ErrAlreadyExists)

		uc := useruc.NewUserUseCases(store)

		_, err := uc.Register(context.Background(), bookclub.RegisterParams{
			Email:    "john@doe.com",
			Password: "123abc!A#4%",
		})

		require.Error(t, err)
		require.Equal(t, bookclub.ErrAlreadyExists, err)
	})
}

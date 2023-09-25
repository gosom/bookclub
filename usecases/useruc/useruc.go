package useruc

import (
	"context"

	"github.com/gosom/bookclub"
)

var _ bookclub.UserUseCases = (*userUseCases)(nil)

type userUseCases struct {
	store bookclub.Storage
}

func NewUserUseCases(store bookclub.Storage) bookclub.UserUseCases {
	ans := &userUseCases{
		store: store,
	}

	return ans
}

func (o *userUseCases) Register(ctx context.Context, params bookclub.RegisterParams) (bookclub.User, error) {
	email, err := bookclub.NewEmail(params.Email)
	if err != nil {
		return bookclub.User{}, err
	}
	passwd, err := bookclub.NewPassword(params.Password)
	if err != nil {
		return bookclub.User{}, err
	}

	ans, err := o.store.CreateUser(ctx, email, passwd)

	return ans, err
}

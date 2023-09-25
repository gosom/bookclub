package bookclub

import (
	"context"
)

//go:generate mockgen -package mocks -destination ./mocks/mock_user_uc.go . UserUseCases

type RegisterParams struct {
	Email    string
	Password string
}

type UserUseCases interface {
	Register(ctx context.Context, params RegisterParams) (User, error)
}

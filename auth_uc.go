package bookclub

import "context"

//go:generate mockgen -package mocks -destination ./mocks/mock_auth_uc.go . AuthUseCases

type LoginParams struct {
	Email    string
	Password string
}

type AuthUseCases interface {
	Login(ctx context.Context, params LoginParams) (string, string, error)
}

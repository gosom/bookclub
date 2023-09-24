package bookclub

import "context"

//go:generate mockgen -package mocks -destination ./mocks/mock_storage.go . Storage

type Storage interface {
	CreateUser(ctx context.Context, email Email, passwd Password) (User, error)
}

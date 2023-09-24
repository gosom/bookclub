package postgres

import (
	"context"

	"github.com/gosom/bookclub"
	"github.com/gosom/bookclub/postgres/db"
)

var _ bookclub.Storage = (*storage)(nil)

type storage struct {
	q *db.Queries
}

func New(conn db.DBTX) bookclub.Storage {
	return &storage{
		q: db.New(conn),
	}
}

func (s *storage) CreateUser(
	ctx context.Context,
	email bookclub.Email,
	passwd bookclub.Password,
) (bookclub.User, error) {
	params := db.CreateUserParams{
		Email:  string(email),
		Passwd: []byte(passwd),
	}

	user, err := s.q.CreateUser(ctx, params)
	if err != nil {
		return bookclub.User{}, err
	}

	return dbUserToUser(user), nil
}

func dbUserToUser(dbUser db.User) bookclub.User {
	ans := bookclub.User{
		ID:        dbUser.ID,
		Email:     bookclub.Email(dbUser.Email),
		CreatedAt: dbUser.CreatedAt.Time,
		UpdatedAt: dbUser.UpdatedAt.Time,
	}

	return ans
}

package postgres

import (
	"context"
	"errors"

	"github.com/gosom/bookclub"
	"github.com/gosom/bookclub/postgres/db"
	"github.com/jackc/pgx/v5/pgconn"
)

const (
	uniqueViolation = "23505"
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
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == uniqueViolation {
				return bookclub.User{}, bookclub.ErrAlreadyExists
			}
		}

		return bookclub.User{}, bookclub.ErrInternalError
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

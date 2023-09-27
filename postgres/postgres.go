package postgres

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/gosom/bookclub"
	"github.com/gosom/bookclub/postgres/db"
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

func (s *storage) GetUserByEmail(
	ctx context.Context,
	email bookclub.Email,
) (bookclub.User, error) {
	user, err := s.q.GetUserByEmail(ctx, string(email))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return bookclub.User{}, bookclub.ErrNotFound
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

	ans.SetPassword(dbUser.Passwd)

	return ans
}

func Migrate(dbAddr, migrationsPath string) error {
	fmt.Println(dbAddr)
	_, second, ok := strings.Cut(dbAddr, "://")
	if !ok {
		return errors.New("invalid dbAddr")
	}

	dburl := "pgx5" + "://" + second
	m, err := migrate.New(fmt.Sprintf("file:%s", migrationsPath), dburl)
	if err != nil {
		return err
	}

	defer m.Close()

	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	return nil
}

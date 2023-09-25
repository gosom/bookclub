package postgres_test

import (
	"context"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"

	"github.com/gosom/bookclub"
	"github.com/gosom/bookclub/postgres"
)

var testDbInstance *pgxpool.Pool

func TestMain(m *testing.M) {
	testDB := setupTestDatabase()
	testDbInstance = testDB.DbInstance
	defer testDB.TearDown()

	if err := migrateDb(testDB.DbAddress); err != nil {
		panic(err)
	}

	os.Exit(m.Run())
}

func Test_Integration_TestQuery(t *testing.T) {
	const q = "SELECT 1"

	var result int
	err := testDbInstance.QueryRow(context.Background(), q).Scan(&result)
	require.NoError(t, err)

	require.Equal(t, 1, result)
}

func Test_Integration_CreateUser(t *testing.T) {
	store := postgres.New(testDbInstance)

	t.Run("should create a user", func(t *testing.T) {
		defer truncateTables(t)

		email, err := bookclub.NewEmail("giorgos@giorgos.com")
		require.NoError(t, err)

		passwd, err := bookclub.NewPassword("123aA@123123")
		require.NoError(t, err)

		user, err := store.CreateUser(context.Background(), email, passwd)
		require.NoError(t, err)

		require.Equal(t, email, user.Email)
		require.NotZero(t, user.ID)

		// check also the database
		q := `SELECT id, email FROM users WHERE id = $1;`
		var id uuid.UUID
		var emailFromDB string

		err = testDbInstance.QueryRow(context.Background(), q, user.ID).Scan(&id, &emailFromDB)
		require.NoError(t, err)

		require.Equal(t, user.ID, id)
		require.Equal(t, string(user.Email), emailFromDB)
	})

	t.Run("will return ErrAlreadyExists when trying to create a user with an existing email", func(t *testing.T) {
		defer truncateTables(t)

		email, err := bookclub.NewEmail("giorgos@giorgos.com")
		require.NoError(t, err)

		passwd, err := bookclub.NewPassword("123aA@123123")
		require.NoError(t, err)

		_, err = store.CreateUser(context.Background(), email, passwd)
		require.NoError(t, err)

		_, err = store.CreateUser(context.Background(), email, passwd)
		require.Error(t, err)

		require.ErrorIs(t, bookclub.ErrAlreadyExists, err)
	})
}

func truncateTables(t *testing.T) {
	t.Helper()
	q := `TRUNCATE TABLE users RESTART IDENTITY CASCADE;`

	_, err := testDbInstance.Exec(context.Background(), q)
	require.NoError(t, err)
}

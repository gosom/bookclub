// Disclaimer: slightly modified from https://medium.com/@dilshataliev/integration-tests-with-golang-test-containers-and-postgres-abb49e8096c5
package postgres_test

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"runtime"
	"time"

	"github.com/golang-migrate/migrate/v4"

	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	dbName = "test_db"
	dbUser = "test_user"
	dbPass = "test_password"
)

type TestDatabase struct {
	DbInstance *pgxpool.Pool
	DbAddress  string
	container  testcontainers.Container
}

func (tdb *TestDatabase) TearDown() {
	tdb.DbInstance.Close()
	// remove test container
	_ = tdb.container.Terminate(context.Background())
}

func setupTestDatabase() *TestDatabase {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	defer cancel()

	container, dbInstance, dbAddr, err := createContainer(ctx)
	if err != nil {
		panic(err)
	}

	return &TestDatabase{
		container:  container,
		DbInstance: dbInstance,
		DbAddress:  dbAddr,
	}
}

func createContainer(ctx context.Context) (testcontainers.Container, *pgxpool.Pool, string, error) {
	var env = map[string]string{
		"POSTGRES_PASSWORD": dbPass,
		"POSTGRES_USER":     dbUser,
		"POSTGRES_DB":       dbName,
	}
	var port = "5432/tcp"

	req := testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "postgres:16.0-alpine",
			ExposedPorts: []string{port},
			Env:          env,
			WaitingFor:   wait.ForLog("database system is ready to accept connections"),
		},
		Started: true,
	}

	container, err := testcontainers.GenericContainer(ctx, req)
	if err != nil {
		return container, nil, "", fmt.Errorf("failed to start container: %v", err)
	}

	p, err := container.MappedPort(ctx, "5432")
	if err != nil {
		return container, nil, "", fmt.Errorf("failed to get container external port: %v", err)
	}

	time.Sleep(time.Second * 1)

	dbAddr := fmt.Sprintf("localhost:%s", p.Port())

	db, err := pgxpool.New(ctx, fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
		dbUser, dbPass, dbAddr, dbName))
	if err != nil {
		return container, db, dbAddr, fmt.Errorf("failed to establish database connection: %v", err)
	}

	return container, db, dbAddr, nil
}

func migrateDb(dbAddr string) error {
	_, path, _, ok := runtime.Caller(0)
	if !ok {
		return fmt.Errorf("failed to get path")
	}
	pathToMigrationFiles := filepath.Dir(path) + "/../scripts/migrations"

	databaseURL := fmt.Sprintf("pgx5://%s:%s@%s/%s?sslmode=disable",
		dbUser, dbPass, dbAddr, dbName)

	m, err := migrate.New(
		fmt.Sprintf("file:%s", pathToMigrationFiles),
		databaseURL,
	)
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

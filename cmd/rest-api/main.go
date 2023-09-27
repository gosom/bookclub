package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/kelseyhightower/envconfig"

	"github.com/gosom/bookclub/api"
	"github.com/gosom/bookclub/api/schema"
	"github.com/gosom/bookclub/jwtprovider"
	"github.com/gosom/bookclub/postgres"
	"github.com/gosom/bookclub/usecases/authuc"
	"github.com/gosom/bookclub/usecases/useruc"
)

type Config struct {
	PGURL      string `default:"postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"`
	JWT_ISSUER string `default:"bookclub"`
	JWT_SECRET string `default:"secret"`
}

func (c *Config) FromEnv() error {
	const appPrefix = "BOOKCLUB"
	return envconfig.Process(appPrefix, c)
}

func main() {
	var cfg Config
	if err := cfg.FromEnv(); err != nil {
		panic(err)
	}

	var migrations string
	flag.StringVar(&migrations, "migrations", "", "run migrations from the set path")
	flag.Parse()

	if migrations != "" {
		fmt.Println("Running migrations...")
		fmt.Println("Migrations path:", migrations)
		fmt.Println("Database URL:", cfg.PGURL)
		if err := postgres.Migrate(cfg.PGURL, migrations); err != nil {
			panic(err)
		}

		return
	}

	dbpool, err := pgxpool.New(context.Background(), cfg.PGURL)

	if err != nil {
		panic(err)
	}

	defer dbpool.Close()

	if err := dbpool.Ping(context.Background()); err != nil {
		panic(err)
	}

	store := postgres.New(dbpool)
	jwtProv := jwtprovider.New(cfg.JWT_SECRET, cfg.JWT_ISSUER)

	userUC := useruc.NewUserUseCases(store)
	authUC := authuc.NewAuthUseCases(store, jwtProv)

	bookclubAPI := api.NewBooklubAPI(userUC, authUC)

	r := chi.NewRouter()

	schema.HandlerFromMux(bookclubAPI, r)

	s := &http.Server{
		Handler: r,
		Addr:    ":8080",
	}

	if err := s.ListenAndServe(); err != nil {
		panic(err)
	}
}

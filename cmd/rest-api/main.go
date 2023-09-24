package main

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/gosom/bookclub/api"
	"github.com/gosom/bookclub/api/schema"
	"github.com/gosom/bookclub/postgres"
	"github.com/gosom/bookclub/usecases/useruc"
)

func main() {
	dburl := "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"
	dbpool, err := pgxpool.New(context.Background(), dburl)

	if err != nil {
		panic(err)
	}

	defer dbpool.Close()

	if err := dbpool.Ping(context.Background()); err != nil {
		panic(err)
	}

	store := postgres.New(dbpool)

	userUC := useruc.NewUserUseCases(store)

	bookclubAPI := api.NewBooklubAPI(userUC)

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

package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/gosom/bookclub/api"
	"github.com/gosom/bookclub/api/schema"
)

func main() {
	bookclubAPI := api.NewBooklubAPI()

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

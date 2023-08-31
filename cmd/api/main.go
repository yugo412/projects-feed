package main

import (
	"projects-rss/cmd/api/handlers"

	"github.com/go-chi/chi/v5"
)

import (
	"net/http"
)

func main() {
	c := chi.NewRouter()

	c.Get("/projects", handlers.GetProjects)

	err := http.ListenAndServe(":3000", c)
	if err != nil {
		panic(err)
	}
}

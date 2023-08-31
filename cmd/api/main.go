package main

import (
	"fmt"
	"log"
	"os"
	"projects-feed/cmd/api/handlers"

	"github.com/go-chi/chi/v5"
)

import (
	"net/http"
)

func main() {
	c := chi.NewRouter()

	c.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/projects/rss", http.StatusSeeOther)
	})

	c.Get("/projects", handlers.GetProjects)
	c.Get("/projects/{type}", handlers.GetProjectsFeed)
	c.Get("/projects/go", handlers.RedirectProject)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	env := os.Getenv("ENV")
	if env != "production" {
		log.Printf("Running [%s] in port: %s", env, port)
	}

	err := http.ListenAndServe(fmt.Sprintf(":%s", port), c)
	if err != nil {
		panic(err)
	}
}

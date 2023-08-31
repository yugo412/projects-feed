package main

import (
	"fmt"
	"log"
	"os"
	"projects-rss/cmd/api/handlers"

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

	if os.Getenv("ENV") != "production" {
		log.Println("Running in port:", port)
	}

	err := http.ListenAndServe(fmt.Sprintf(":%s", port), c)
	if err != nil {
		panic(err)
	}
}

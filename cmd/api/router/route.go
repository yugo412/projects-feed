package router

import (
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

var c *chi.Mux

func init() {
	c = chi.NewRouter()

	env := os.Getenv("ENV")
	if env != "production" {
		c.Use(middleware.RequestID)
		c.Use(middleware.Logger)
	}

	// register built-in middleware
	c.Use(middleware.Heartbeat("/ping"))
	c.Use(middleware.Throttle(100))
	c.Use(middleware.Timeout(10 * time.Second))
	c.Use(middleware.Recoverer)
}

func RegisterRoutes() *chi.Mux {
	c.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
		w.Write([]byte("Page not found."))
	})

	c.Route("/projects", projectsRouter)

	return c
}

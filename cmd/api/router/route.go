package router

import (
	"net/http"
	"os"
	"projects-feed/cmd/api/handler"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"gorm.io/gorm"
)

var c *chi.Mux

type App struct {
	DB *gorm.DB
}

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

	// service static files
	fileServer := http.FileServer(http.Dir("./web/asset"))
	c.Handle("/public/*", http.StripPrefix("/public/", fileServer))
}

func RegisterRoutes(app *App) *chi.Mux {
	c.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
		w.Write([]byte("Page not found."))
	})

	project := handler.NewProject(&handler.Project{
		DB: app.DB,
	})
	c.Get("/", project.Index)
	c.Get("/projects", project.GetProjects)
	c.Get("/projects/{type}", project.GetProjectsFeed)
	c.Get("/projects/go", project.RedirectProject)

	return c
}

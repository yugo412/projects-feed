package router

import (
	"projects-feed/cmd/api/handler"

	"github.com/go-chi/chi/v5"
)

// projectsRouter registers related projects' routes to router.
// All routes are prefixed with "/projects" by default.
func projectsRouter(r chi.Router) {
	r.Get("/", handler.GetProjects)
	r.Get("/{type}", handler.GetProjectsFeed)
	r.Get("/go", handler.RedirectProject)
}

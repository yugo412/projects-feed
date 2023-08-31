package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"projects-feed/pkg/projects"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/feeds"
)

func GetProjects(w http.ResponseWriter, r *http.Request) {
	tag := ""
	if q := r.URL.Query().Get("tag"); q != "" {
		tag = q
	}

	page := 1
	if n, err := strconv.Atoi(r.URL.Query().Get("page")); err == nil {
		page = n
	}

	projects, err := projects.GetProjects(uint(page), tag)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// convert result to json
	body, _ := json.Marshal(projects)
	w.Write(body)

	return
}

func GetProjectsFeed(w http.ResponseWriter, r *http.Request) {
	feed := &feeds.Feed{
		Title:       "Projects.co.id",
		Description: "Kerja Online Hasil Maksimal",
		Author: &feeds.Author{
			Name: "Projects.co.id",
		},
		Link: &feeds.Link{
			Href: "https://projects.co.id/public/browse_projects/listing",
		},
		Created: time.Now().Local(),
	}

	tag := ""
	if q := r.URL.Query().Get("tag"); q != "" {
		tag = q
	}
	projects, err := projects.GetProjects(1, tag)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))

		return
	}

	for _, p := range projects {
		item := &feeds.Item{
			Title: p.Title,
			Link: &feeds.Link{
				Href: fmt.Sprintf("%s/projects/go?to=%s", os.Getenv("URL"), p.URL),
			},
			Description: p.Description,
			Author: &feeds.Author{
				Name: p.Author.Name,
			},
			Created: p.PublishedDate,
		}

		feed.Items = append(feed.Items, item)
	}

	output := "rss"
	if t := chi.URLParam(r, "type"); t != "" {
		output = strings.ToLower(t)
	}

	var body string
	switch output {
	case "json":
		w.Header().Set("Content-Type", "application/json")
		body, _ = feed.ToJSON()
	case "atom":
		w.Header().Set("Content-Type", "application/atom+xml")
		body, _ = feed.ToAtom()
	default:
		w.Header().Set("Content", "application/rss+xml")
		body, _ = feed.ToRss()
	}

	w.Write([]byte(body))

	return
}

func RedirectProject(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, r.URL.Query().Get("to"), http.StatusSeeOther)
}

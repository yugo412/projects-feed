package handler

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path"
	"projects-feed/internal"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/feeds"
)

func Index(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles(path.Join("web", "template", "index.html"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))

		return
	}

	var vendor string
	if v := r.URL.Query().Get("vendor"); v != "" {
		vendor = v
	}

	var tag string
	if t := r.URL.Query().Get("tag"); t != "" {
		tag = t
	}

	projects, err := internal.GetProjects(vendor, 1, tag)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error."))

		log.Println("Failed to get projects:", err)

		return
	}

	title := "Projects Feed"
	if tag != "" {
		title = fmt.Sprintf("%s for %s", title, tag)
	}
	if vendor != "" {
		title = fmt.Sprintf("%s from %s", title, vendor)
	}

	err = t.Execute(w, map[string]interface{}{
		"projects": projects,
		"vendor":   vendor,
		"tag":      tag,
		"query":    template.URL(r.URL.RawQuery),
		"title":    title,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))

		return
	}
}

func GetProjects(w http.ResponseWriter, r *http.Request) {
	tag := ""
	if q := r.URL.Query().Get("tag"); q != "" {
		tag = q
	}

	page := 1
	if n, err := strconv.Atoi(r.URL.Query().Get("page")); err == nil {
		page = n
	}

	vendor := ""
	if v := r.URL.Query().Get("vendor"); v != "" {
		vendor = v
	}

	projects, err := internal.GetProjects(vendor, page, tag)
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

	var vendor string
	if v := r.URL.Query().Get("vendor"); v != "" {
		vendor = v
	}

	tag := ""
	if q := r.URL.Query().Get("tag"); q != "" {
		tag = q
	}

	projects, err := internal.GetProjects(vendor, 1, tag)
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
			Created: p.PublishedAt,
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

package handlers

import (
	"encoding/json"
	"net/http"
	"projects-rss/pkg"
)

func GetProjects(w http.ResponseWriter, r *http.Request) {
	tag := ""
	if q := r.URL.Query().Get("tag"); q != "" {
		tag = q
	}

	projects, err := pkg.GetProjects(1, tag)
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

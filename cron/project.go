package cron

import (
	"projects-feed/pkg/cache"
	"projects-feed/pkg/project"
	"projects-feed/srv"
	"strings"

	"github.com/gookit/slog"
)

func PrefetchProjects() (err error) {
	_, err = srv.GetProjects("", 1, "")

	return
}

func UpdateProject() (err error) {
	c, err := cache.New("memory")
	if err != nil {
		slog.Errorf("failed to initialize cache: %v", err)
	}

	// get all stored data in cache by key
	for k, v := range c.Items() {
		var cacheProjects []project.Project
		if strings.Contains(k, "projects_") {
			projects := v.([]project.Project)
			for _, p := range projects {
				if p.Description == "" {
					detail, err := srv.GetDetail(p.URL, p.Vendor)
					if err == nil {
						p.Description = detail.Description

						// default image from sribu is broken
						//p.Author.AvatarURL = detail.Author.AvatarURL
					}
				}
				cacheProjects = append(cacheProjects, p)
			}
		}

		if ok, _ := c.Set(k, cacheProjects); !ok {
			slog.Errorf("Failed to set cache for \"%s\".", k)
		}
	}

	return err
}

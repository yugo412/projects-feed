package cron

import (
	"projects-feed/pkg/cache"
	"projects-feed/pkg/project"
	"projects-feed/srv"
	"strings"

	"github.com/gookit/slog"
)

func UpdateProject() (err error) {
	c := cache.New("memory")

	// get all stored data in cache by key
	for k, v := range c.Items() {
		var cache []project.Project
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
				cache = append(cache, p)
			}
		}

		if ok, _ := c.Set(k, cache); !ok {
			slog.Errorf("Failed to set cache for \"%s\".", k)
		}
	}

	return err
}

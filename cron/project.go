package cron

import (
	"net/url"
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
				avatar, err := url.Parse(p.Author.AvatarURL)
				if avatar.Host == "" || err != nil {
					_, err := srv.GetDetail(p.URL, p.Vendor)
					if err == nil {
						p.Author.Name = "TEST"
						p.Author.AvatarURL = ""
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

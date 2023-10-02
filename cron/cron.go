package cron

import (
	"github.com/gookit/slog"
	"github.com/robfig/cron/v3"
)

func Run() {
	var err error
	c := cron.New()

	_, err = c.AddFunc("@every 1m", func() {
		// by default, it will try to fetch projects every minute
		// but, since the projects is stored to the cache for 10 minutes
		// the projects won't re-fetch from the source until it expires
		err := PrefetchProjects()
		if err != nil {
			slog.Errorf("Failed to fetch projects: %v", err)
		}
	})
	if err != nil {
		slog.Errorf("Failed to run cron: %v", err)
	}

	_, err = c.AddFunc("@every 30s", func() {
		// same behavior with prefetch-project
		// only fetch project detail from source if there are no additional infos available
		err := UpdateProject()
		if err != nil {
			slog.Errorf("Failed to update existing projects: %v", err)
		}
	})
	if err != nil {
		slog.Errorf("Failed to run cron: %v", err)
	}

	c.Start()
}

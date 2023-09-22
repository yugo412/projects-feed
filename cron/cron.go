package cron

import (
	"github.com/gookit/slog"
	"github.com/robfig/cron/v3"
)

func Run() {
	c := cron.New()
	id, err := c.AddFunc("@every 10m", FetchProjects)
	if err != nil {
		slog.Errorf("Failed to fetch projects with cron ID %s: %v", id, err)
	}

	c.Start()
}

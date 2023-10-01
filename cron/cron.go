package cron

import (
	"github.com/gookit/slog"
	"github.com/robfig/cron/v3"
)

func Run() {
	var err error
	c := cron.New()

	_, err = c.AddFunc("@every 10s", func() {
		err := PrefectProject()
		if err != nil {
			slog.Errorf("Failed to fetch projects: %v", err)
		}
	})
	if err != nil {
		slog.Errorf("Failed to run cron: %v", err)
	}

	_, err = c.AddFunc("@every 1m", func() {
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

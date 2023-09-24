package cron

import (
	"github.com/gookit/slog"
	"github.com/robfig/cron/v3"

	"gorm.io/gorm"
)

type Cron struct {
	DB *gorm.DB
}

func Run(e Cron) {
	var err error

	c := cron.New()

	_, err = c.AddFunc("@every 5m", func() {
		err := FetchProjects(e.DB)
		if err != nil {
			slog.Errorf("Failed to fetch projects: %v", err)
		}
	})
	if err != nil {
		slog.Errorf("Failed to run cron: %v", err)
	}

	_, err = c.AddFunc("@every 1m", func() {
		err := UpdateProject(e.DB)
		if err != nil {
			slog.Errorf("Failed to update existing projects: %v", err)
		}
	})
	if err != nil {
		slog.Errorf("Failed to run cron: %v", err)
	}

	c.Start()
}

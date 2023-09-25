package cron

import (
	"errors"
	"projects-feed/models"
	"projects-feed/pkg/project"
	"sync"
	"time"

	"github.com/gookit/slog"
	"gorm.io/gorm"
)

func vendors() []models.Vendor {
	return []models.Vendor{
		models.Vendor{
			Brand:  "Projects.co.id",
			Name:   "projects",
			URL:    "https://www.projects.co.id",
			Source: "https://projects.co.id",
		},
		models.Vendor{
			Brand:  "Sribu.com",
			Name:   "sribu",
			URL:    "https://www.sribu.com",
			Source: "https://api.sribu.com",
		},
	}
}

func FetchProjects(db *gorm.DB) error {
	var projects []project.Project
	var wg sync.WaitGroup
	for _, v := range vendors() {
		wg.Add(1)

		go func(v models.Vendor) {
			var count int64
			err := db.Model(&models.Vendor{}).
				Where("name = ?", v.Name).
				Count(&count).
				Error
			if err != nil {
				return
			}

			if count <= 0 {
				db.Create(&v)
			}

			vendor := project.New(v.Name)
			latest, err := vendor.GetProjects(1, "")
			if err != nil {
				slog.Errorf("Failed to fetch project from source: %v", err)
				return
			}

			projects = append(projects, latest...)

			wg.Done()
		}(v)
	}

	wg.Wait()

	var records []models.Project
	for _, p := range projects {
		var count int64
		err := db.Model(&models.Project{}).
			Where("title = ?", p.Title).
			Count(&count).
			Error
		if err != nil {
			slog.Errorf("Failed to count project: %v", err)
			continue
		}

		if count <= 0 {
			// vendor should be created before this line
			var vendor models.Vendor
			err = db.Model(&models.Vendor{}).
				Where("brand = ?", p.Vendor).
				First(&vendor).
				Error
			if err != nil {
				slog.Errorf("Vendor with name \"%s\" does not exists.", p.Vendor)
			}

			// create unique author for each vendor
			var author models.Author
			err = db.Model(&models.Author{}).
				Where("username = ? AND vendor_id = ?", p.Author.Username, vendor.Model.ID).
				First(&author).
				Error
			if errors.Is(err, gorm.ErrRecordNotFound) {
				author = models.Author{
					Name:     p.Author.Name,
					Username: p.Author.Username,
					Avatar:   p.Author.AvatarURL,
					URL:      p.Author.URL,
				}
				db.Create(&author)
			}

			records = append(records, models.Project{
				Title:       p.Title,
				Description: p.Description,
				PublishedAt: p.PublishedAt,
				URL:         p.URL,
				MinBudget:   p.Budget.Min,
				MaxBudget:   p.Budget.Max,
				Vendor:      vendor,
				Author:      author,
			})
		}
	}

	if len(records) >= 1 {
		db.Create(&records)
	}

	return nil
}

func UpdateProject(db *gorm.DB) (err error) {
	var projects []models.Project
	err = db.Model(&models.Project{}).
		Preload("Vendor").
		Where("description = '' AND published_at > ?", time.Now().AddDate(0, 0, -7)).
		Take(&projects).
		Error

	return err
}

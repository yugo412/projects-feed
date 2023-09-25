package srv

import (
	"fmt"
	"projects-feed/models"
	"projects-feed/pkg/project"
	"strconv"

	"gorm.io/gorm"
)

type Project struct {
	DB *gorm.DB
}

func NewProject(p *Project) *Project {
	return p
}

func (p Project) GetProjects(vendor string, page int, tag string) (projects []project.Project, err error) {
	var ps []models.Project
	err = p.DB.Model(&models.Project{}).
		Debug().
		Preload("Vendor").
		Preload("Author").
		Where("title LIKE ?", fmt.Sprintf("%%%s%%", tag)).
		Order("published_at DESC").
		Limit(20).
		Find(&ps).
		Error
	if err != nil {
		return
	}

	for _, x := range ps {
		projects = append(projects, project.Project{
			ID:          strconv.Itoa(int(x.ID)),
			Vendor:      x.Vendor.Brand,
			Title:       x.Title,
			URL:         x.URL,
			Description: x.Description,
			PublishedAt: x.PublishedAt,
			Tags:        nil,
			Author: project.Author{
				Name:      x.Author.Name,
				Username:  x.Author.Username,
				URL:       x.Author.URL,
				AvatarURL: x.Author.Avatar,
			},
			Budget: project.Budget{
				Min: x.MinBudget,
				Max: x.MaxBudget,
			},
		})
	}

	return
}

package projects

import (
	"strings"
	"time"

	"github.com/xeonx/timeago"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

type Vendor interface {
	Name() string
	GetProjects(int, string) ([]Project, error)
}

type Author struct {
	Name      string `json:"name"`
	URL       string `json:"url"`
	AvatarURL string `json:"avatar_url"`
}

type Budget struct {
	Min float64 `json:"min"`
	Max float64 `json:"max"`
}

type Project struct {
	Vendor      string    `json:"vendor"`
	Title       string    `json:"title"`
	URL         string    `json:"url"`
	Description string    `json:"description"`
	PublishedAt time.Time `json:"published_at"`
	Tags        []string  `json:"tags"`
	Author      Author    `json:"author"`
	Budget      Budget    `json:"budget"`
}

func New(vendor string) Vendor {
	vendor = strings.TrimSpace(vendor)

	switch vendor {
	case "projects":
		return &Projects{
			BaseURL: "https://projects.co.id",
		}
	}

	return nil
}

func (p Project) LimitedTags(limit int) []string {
	if len(p.Tags) > limit {
		return p.Tags[:limit]
	}

	return p.Tags
}

func (p Project) RemainTag() int {
	if len(p.Tags) > 3 {
		return len(p.Tags) - 3
	}

	return 0
}

func (p Project) Timeago() string {
	return timeago.English.Format(p.PublishedAt)
}

func (p Project) GetBudget() string {
	printer := message.NewPrinter(language.Indonesian)

	return printer.Sprintf("Rp%.0f - Rp%0.f", p.Budget.Min, p.Budget.Max)
}

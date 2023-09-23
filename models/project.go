package models

import (
	"gorm.io/gorm"
	"time"
)

type Project struct {
	gorm.Model
	VendorID    uint      `gorm:"column:vendor_id"`
	AuthorID    uint      `gorm:"column:author_id"`
	Title       string    `gorm:"column:title"`
	Description string    `gorm:"column:description"`
	PublishedAt time.Time `gorm:"column:published_at"`
	MinBudget   float64   `gorm:"column:min_budget"`
	MaxBudget   float64   `gorm:"column:max_budget"`

	Vendor Vendor
	Author Author
}

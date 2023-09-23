package models

import "gorm.io/gorm"

type Vendor struct {
	gorm.Model
	Brand  string `gorm:"column:brand"`
	Name   string `gorm:"column:name"`
	URL    string `gorm:"column:url"`
	Source string `gorm:"column:source"`

	Project []Project
	Author  []Author
}

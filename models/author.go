package models

import "gorm.io/gorm"

type Author struct {
	gorm.Model
	VendorID int    `gorm:"column:vendor_id"`
	Name     string `gorm:"column:name"`
	Username string `gorm:"column:username"`
	URL      string `gorm:"column:url"`
	Avatar   string `gorm:"column:avatar"`

	Vendor Vendor
}

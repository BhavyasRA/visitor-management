package models

import "time"

type Visitor struct {
	ID uint `gorm:"primaryKey"`

	ImageURL string

	Name   string
	Mobile string `gorm:"unique"`
	Email  string

	IsRestricted bool `gorm:"default:false"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

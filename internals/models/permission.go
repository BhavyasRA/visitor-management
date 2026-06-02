package models

import "time"

type Permission struct {
	ID uint `gorm:"primaryKey"`

	Name string `gorm:"unique"`

	CreatedAt time.Time
}
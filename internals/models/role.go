package models

import "time"

type Role struct {
	ID uint `gorm:"primaryKey"`

	Name string `gorm:"unique"`

	Permissions []Permission `gorm:"many2many:role_permissions"`

	CreatedAt time.Time
}
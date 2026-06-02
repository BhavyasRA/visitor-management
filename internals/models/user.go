package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID uint `gorm:"primaryKey"`

	Name  string
	Email string `gorm:"unique"`
	Phone string `gorm:"unique"`

	IsActive bool `gorm:"default:true"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	Roles []Role `gorm:"many2many:user_roles"`
}
package models

import "time"

type Visitor struct {
	ID uint `gorm:"primaryKey" json:"id"`

	Name   string `json:"name"`
	Mobile string `gorm:"unique" json:"mobile"`
	Email  string `json:"email"`

	IsRestricted bool `gorm:"default:false" json:"is_restricted"`

	Documents []VisitorDocument `gorm:"foreignKey:VisitorID" json:"documents"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

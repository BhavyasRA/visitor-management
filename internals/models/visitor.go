package models

import "time"

type Visitor struct {
	ID uint `gorm:"primaryKey"`

	Name   string
	Mobile string
	Email  string

	IsRestricted bool `gorm:"default:false"`

	ToWhom uint

	Purpose string

	VisitingTill *time.Time

	ExitAt *time.Time

	CreatedAt time.Time
}
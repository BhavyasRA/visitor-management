package models

import "time"

type EntryLog struct {
	ID uint `gorm:"primaryKey"`

	VisitorID uint
	Visitor   Visitor `gorm:"foreignKey:VisitorID"`

	PersonToMeet uint

	Purpose string
	Status  string `gorm:"default:'active'"`

	VisitingTill *time.Time
	EnteredAt    time.Time
	ExitedAt     *time.Time

	CreatedAt time.Time
	UpdatedAt time.Time
}

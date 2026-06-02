package models

import "time"

type EntryLog struct {
	ID uint `gorm:"primaryKey"`

	VisitorID uint

	EnteredAt time.Time

	ExitedAt *time.Time

	CreatedBy uint

	CreatedAt time.Time
}
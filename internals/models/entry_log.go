package models

import "time"

type EntryLog struct {
	ID uint `gorm:"primaryKey" json:"id"`

	VisitorID uint    `json:"visitor_id"`
	Visitor   Visitor `gorm:"foreignKey:VisitorID" json:"visitor"`

	PersonToMeet uint `json:"person_to_meet"`

	Purpose string `json:"purpose"`
	Status  string `gorm:"default:'active'" json:"status"`

	VisitingTill *time.Time `json:"visiting_till"`
	EnteredAt    time.Time  `json:"entered_at"`
	ExitedAt     *time.Time `json:"exited_at"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

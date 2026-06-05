package models

import "time"

type VisitorDocument struct {
	ID uint `gorm:"primaryKey" json:"id"`

	VisitorID           uint    `json:"visitor_id"`
	Visitor             Visitor `gorm:"foreignKey:VisitorID" json:"visitor"`
	PhotoURL            string  `json:"photo_url"`
	IdentityDocumentURL string  `json:"identity_document"`
	DocumentType        string  `json:"document_type"`
	DocumentNumber      string  `json:"document_number"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

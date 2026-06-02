package dto

type MakeEntryDTO struct {
	VisitorID uint `json:"visitor_id" validate:"required"`
}
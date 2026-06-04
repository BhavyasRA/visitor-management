package dto

type CreateVisitorDTO struct {
	ImageURL     string `json:"image_url" validate:"omitempty,url"`
	Name         string `json:"name" validate:"required,min=2"`
	Mobile       string `json:"mobile" validate:"required,min=10,max=15"`
	Email        string `json:"email" validate:"omitempty,email"`
	Purpose      string `json:"purpose" validate:"required"`
	PersonToMeet uint   `json:"person_to_meet" validate:"required"`
	VisitingTill string `json:"visiting_till" validate:"omitempty"`
}

type UpdateVisitorDTO struct {
	ImageURL string `json:"image_url" validate:"omitempty,url"`
	Name     string `json:"name" validate:"required,min=2"`
	Mobile   string `json:"mobile" validate:"required,min=10,max=15"`
	Email    string `json:"email" validate:"omitempty,email"`
	Purpose  string `json:"purpose" validate:"required"`
}

type VisitorFilter struct {
	Name     string `query:"name"`
	Mobile   string `query:"mobile"`
	Email    string `query:"email"`
	From     string `query:"from"`
	To       string `query:"to"`
	IsActive string `query:"is_active"`
}

type VisitorListItemDTO struct {
	ID             uint   `json:"id"`
	Name           string `json:"name"`
	PurposeOfVisit string `json:"purpose"`
	Status         string `json:"status"`
}

type VisitorGroupedDTO struct {
	Date     string               `json:"date"`
	Visitors []VisitorListItemDTO `json:"visitors"`
}

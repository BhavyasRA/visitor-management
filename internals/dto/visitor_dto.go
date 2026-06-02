package dto

type CreateVisitorDTO struct {
	Name         string `json:"name" validate:"required,min=2"`
	Mobile       string `json:"mobile" validate:"required,min=10,max=15"`
	Email        string `json:"email" validate:"omitempty,email"`
	Purpose      string `json:"purpose" validate:"required"`
	ToWhom       uint   `json:"to_whom" validate:"required"`
	VisitingTill string `json:"visiting_till"`
}

type UpdateVisitorDTO struct {
	Name    string `json:"name" validate:"required,min=2"`
	Mobile  string `json:"mobile" validate:"required,min=10,max=15"`
	Email   string `json:"email" validate:"omitempty,email"`
	Purpose string `json:"purpose" validate:"required"`
}
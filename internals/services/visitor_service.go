package services

import (
	"time"

	"entry-system/internals/models"
	"entry-system/internals/repositories"
)

type VisitorService struct {
	visitorRepo *repositories.VisitorRepository
}

func NewVisitorService() *VisitorService {
	return &VisitorService{
		visitorRepo: repositories.NewVisitorRepository(),
	}
}

func (s *VisitorService) CreateVisitor(
	name string,
	mobile string,
	email string,
	purpose string,
	toWhom uint,
	visitingTill *time.Time,
) error {
	visitor := models.Visitor{
		Name:         name,
		Mobile:       mobile,
		Email:        email,
		Purpose:      purpose,
		ToWhom:       toWhom,
		VisitingTill: visitingTill,
	}

	return s.visitorRepo.Create(&visitor)
}

func (s *VisitorService) GetVisitors() ([]models.Visitor, error) {
	return s.visitorRepo.FindAll()
}

func (s *VisitorService) UpdateVisitor(id uint, name, mobile, email, purpose string) error {
	visitor, err := s.visitorRepo.FindByID(id)
	if err != nil {
		return err
	}

	visitor.Name = name
	visitor.Mobile = mobile
	visitor.Email = email
	visitor.Purpose = purpose

	return s.visitorRepo.Update(visitor)
}

func (s *VisitorService) RestrictVisitor(id uint) error {
	return s.visitorRepo.Restrict(id)
}

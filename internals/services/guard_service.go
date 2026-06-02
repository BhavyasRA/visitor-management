package services

import (
	"errors"
	"time"

	"entry-system/internals/models"
	"entry-system/internals/repositories"
)

type GuardService struct {
	visitorRepo *repositories.VisitorRepository
	entryRepo   *repositories.EntryRepository
	userRepo    *repositories.UserRepository
}

func NewGuardService() *GuardService {
	return &GuardService{
		visitorRepo: repositories.NewVisitorRepository(),
		entryRepo:   repositories.NewEntryRepository(),
		userRepo:    repositories.NewUserRepository(),
	}
}

func (s *GuardService) MakeEntry(visitorID uint, guardID uint) error {
	visitor, err := s.visitorRepo.FindByID(visitorID)
	if err != nil {
		return errors.New("visitor not found")
	}

	if visitor.IsRestricted {
		return errors.New("visitor is restricted")
	}

	entry := models.EntryLog{
		VisitorID: visitorID,
		EnteredAt: time.Now(),
		CreatedBy: guardID,
	}

	return s.entryRepo.Create(&entry)
}

func (s *GuardService) SeeEntries(
	filter string,
	from string,
	to string,
) ([]models.EntryLog, error) {

	return s.entryRepo.FindEntriesByFilter(filter, from, to)
}
func (s *GuardService) ExitVisitor(visitorID uint) error {
	return s.entryRepo.ExitVisitor(visitorID)
}

func (s *GuardService) RestrictVisitor(visitorID uint) error {
	return s.visitorRepo.Restrict(visitorID)
}

func (s *GuardService) RestrictEmployee(userID uint) error {
	return s.userRepo.Deactivate(userID)
}

package services

import (
	"errors"
	"time"

	"entry-system/internals/dto"
	"entry-system/internals/models"
	"entry-system/internals/repositories"
)

var AllowedPersonsToMeet = map[uint]string{
	1: "Deepak Swain",
	2: "Adithya",
	3: "Rahul Sharma",
	4: "Amit Kumar",
	5: "Ravi Kumar",
}

type VisitorService struct {
	visitorRepo  *repositories.VisitorRepository
	entryLogRepo *repositories.EntryRepository
}

func NewVisitorService() *VisitorService {
	return &VisitorService{
		visitorRepo:  repositories.NewVisitorRepository(),
		entryLogRepo: repositories.NewEntryRepository(),
	}
}

func (s *VisitorService) CreateVisitor(
	name string,
	mobile string,
	email string,
	photo string,
	identityDocument string,
	purpose string,
	personToMeet uint,
	visitingTill *time.Time,
) error {

	if _, ok := AllowedPersonsToMeet[personToMeet]; !ok {
		return errors.New("invalid person to meet")
	}

	visitor, err := s.visitorRepo.FindByMobile(mobile)

	if err != nil {
		visitor = &models.Visitor{
			Name:             name,
			Mobile:           mobile,
			Email:            email,
			Photo:            photo,
			IdentityDocument: identityDocument,
		}

		if err := s.visitorRepo.Create(visitor); err != nil {
			return err
		}
	} else {
		visitor.Name = name
		visitor.Email = email

		if photo != "" {
			visitor.Photo = photo
		}

		if identityDocument != "" {
			visitor.IdentityDocument = identityDocument
		}

		if err := s.visitorRepo.Update(visitor); err != nil {
			return err
		}
	}

	entry := models.EntryLog{
		VisitorID:    visitor.ID,
		PersonToMeet: personToMeet,
		Purpose:      purpose,
		Status:       "active",
		VisitingTill: visitingTill,
		EnteredAt:    time.Now(),
	}

	return s.entryLogRepo.Create(&entry)
}

func (s *VisitorService) GetVisitors(filter dto.VisitorFilter) ([]models.Visitor, error) {
	return s.visitorRepo.FindAllWithFilters(filter)
}

func (s *VisitorService) UpdateVisitor(id uint, name, mobile, email string) error {
	visitor, err := s.visitorRepo.FindByID(id)
	if err != nil {
		return err
	}

	visitor.Name = name
	visitor.Mobile = mobile
	visitor.Email = email

	return s.visitorRepo.Update(visitor)
}

func (s *VisitorService) RestrictVisitor(id uint) error {
	return s.visitorRepo.Restrict(id)
}

func (s *VisitorService) GetGroupedVisitorEntries(
	status string,
	filter string,
	from string,
	to string,
) ([]dto.VisitorGroupedDTO, error) {

	entries, err := s.entryLogRepo.GetVisitorEntriesByStatusAndDate(
		"active",
		filter,
		from,
		to,
	)

	if err != nil {
		return nil, err
	}

	groupMap := make(map[string][]dto.VisitorListItemDTO)

	for _, entry := range entries {
		date := entry.EnteredAt.Format("2006-01-02")

		statusText := "exited"
		if entry.ExitedAt == nil {
			statusText = "active"
		}

		item := dto.VisitorListItemDTO{
			ID:             entry.VisitorID,
			EntryID:        entry.ID,
			Name:           entry.Visitor.Name,
			PurposeOfVisit: entry.Purpose,
			Status:         statusText,
		}

		groupMap[date] = append(groupMap[date], item)
	}

	var result []dto.VisitorGroupedDTO

	for date, visitors := range groupMap {
		result = append(result, dto.VisitorGroupedDTO{
			Date:     date,
			Visitors: visitors,
		})
	}

	return result, nil
}

func (s *VisitorService) GetVisitorByMobile(mobile string) (*models.Visitor, error) {
	return s.visitorRepo.FindByMobile(mobile)
}

func (s *VisitorService) GetVisitorStats() (map[string]int64, error) {
	return s.entryLogRepo.GetVisitorStats()
}

func (s *VisitorService) ExitVisitor(entryID uint) error {
	return s.entryLogRepo.ExitVisitor(entryID)
}
func (s *VisitorService) GetPersonsDropdown() []map[string]any {
	return []map[string]any{
		{"id": 1, "name": "Deepak Swain"},
		{"id": 2, "name": "Adithya"},
		{"id": 3, "name": "Rahul Sharma"},
		{"id": 4, "name": "Amit Kumar"},
		{"id": 5, "name": "Ravi Kumar"},
	}
}

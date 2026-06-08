package services

import (
	"errors"
	"sort"
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
	visitorRepo         *repositories.VisitorRepository
	entryLogRepo        *repositories.EntryRepository
	visitorDocumentRepo *repositories.VisitorDocumentRepository
}

func NewVisitorService() *VisitorService {
	return &VisitorService{
		visitorRepo:         repositories.NewVisitorRepository(),
		entryLogRepo:        repositories.NewEntryRepository(),
		visitorDocumentRepo: repositories.NewVisitorDocumentRepository(),
	}
}

func (s *VisitorService) CreateVisitor(
	name string,
	mobile string,
	email string,
	photoURL string,
	identityDocumentURL string,
	purpose string,
	personToMeet uint,
	visitingTill *time.Time,
) (*models.EntryLog, *models.Visitor, *models.VisitorDocument, error) {

	if _, ok := AllowedPersonsToMeet[personToMeet]; !ok {
		return nil, nil, nil, errors.New("invalid person to meet")
	}

	visitor, err := s.visitorRepo.FindByMobile(mobile)

	if err != nil {
		visitor = &models.Visitor{
			Name:   name,
			Mobile: mobile,
			Email:  email,
		}

		if err := s.visitorRepo.Create(visitor); err != nil {
			return nil, nil, nil, err
		}
	} else {
		visitor.Name = name
		visitor.Email = email

		if err := s.visitorRepo.Update(visitor); err != nil {
			return nil, nil, nil, err
		}
	}

	document := models.VisitorDocument{
		VisitorID:           visitor.ID,
		PhotoURL:            photoURL,
		IdentityDocumentURL: identityDocumentURL,
	}

	if err := s.visitorDocumentRepo.Create(&document); err != nil {
		return nil, nil, nil, err
	}

	entry := models.EntryLog{
		VisitorID:    visitor.ID,
		PersonToMeet: personToMeet,
		Purpose:      purpose,
		Status:       "active",
		VisitingTill: visitingTill,
		EnteredAt:    time.Now(),
	}

	if err := s.entryLogRepo.Create(&entry); err != nil {
		return nil, nil, nil, err
	}

	return &entry, visitor, &document, nil
}

func (s *VisitorService) GetVisitors(filter dto.VisitorFilter) ([]models.Visitor, error) {
	return s.visitorRepo.FindAllWithFilters(filter)
}

func (s *VisitorService) UpdateVisitor(
	id uint,
	name string,
	mobile string,
	email string,
) error {

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

func (s *VisitorService) ExitVisitor(entryID uint) error {
	return s.entryLogRepo.ExitVisitor(entryID)
}

func (s *VisitorService) GetGroupedVisitorEntries(
	status string,
	filter string,
	from string,
	to string,
) ([]dto.VisitorGroupedDTO, error) {

	entries, err := s.entryLogRepo.GetVisitorEntriesByStatusAndDate(
		status,
		filter,
		from,
		to,
	)

	if err != nil {
		return nil, err
	}

	groupMap := make(map[string][]dto.VisitorListItemDTO)

	ist, _ := time.LoadLocation("Asia/Kolkata")

	for _, entry := range entries {
		date := entry.EnteredAt.In(ist).Format("2006-01-02")

		statusText := entry.Status
		if statusText == "" {
			statusText = "active"

			if entry.ExitedAt != nil {
				statusText = "exited"
			}
		}

		photoURL := ""

		if len(entry.Visitor.Documents) > 0 {
			photoURL = entry.Visitor.Documents[len(entry.Visitor.Documents)-1].PhotoURL
		}

		item := dto.VisitorListItemDTO{
			ID:             entry.VisitorID,
			EntryID:        entry.ID,
			Name:           entry.Visitor.Name,
			PurposeOfVisit: entry.Purpose,
			Status:         statusText,
			EnteredAt: entry.EnteredAt.
				In(ist).
				Format("2006-01-02T15:04:05Z07:00"),
			PhotoURL: photoURL,
		}

		groupMap[date] = append(groupMap[date], item)
	}

	var dates []string

	for date := range groupMap {
		dates = append(dates, date)
	}

	sort.Sort(sort.Reverse(sort.StringSlice(dates)))

	var result []dto.VisitorGroupedDTO

	for _, date := range dates {
		result = append(result, dto.VisitorGroupedDTO{
			Date:     date,
			Visitors: groupMap[date],
		})
	}

	return result, nil
}

func (s *VisitorService) GetActiveEntries() ([]dto.VisitorGroupedDTO, error) {
	return s.entryLogRepo.GetActiveEntries()
}

func (s *VisitorService) GetVisitorByMobile(
	mobile string,
) (*models.Visitor, error) {

	return s.visitorRepo.FindByMobile(mobile)
}

func (s *VisitorService) GetVisitorStats() (map[string]int64, error) {
	return s.entryLogRepo.GetVisitorStats()
}

func (s *VisitorService) GetVisitorDocumentForAI(
	documentID uint,
) (*models.VisitorDocument, error) {

	return s.visitorDocumentRepo.FindByID(documentID)
}

func (s *VisitorService) UpdateDocumentAIResponse(
	documentID uint,
	documentType string,
	documentNumber string,
) error {

	document, err := s.visitorDocumentRepo.FindByID(documentID)
	if err != nil {
		return errors.New("visitor document not found")
	}

	document.DocumentType = documentType
	document.DocumentNumber = documentNumber

	return s.visitorDocumentRepo.Update(document)
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

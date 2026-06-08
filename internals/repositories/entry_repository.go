package repositories

import (
	"errors"
	"sort"
	"time"

	"entry-system/internals/database"
	"entry-system/internals/dto"
	"entry-system/internals/models"
)

type EntryRepository struct{}

func NewEntryRepository() *EntryRepository {
	return &EntryRepository{}
}

func (r *EntryRepository) Create(entry *models.EntryLog) error {
	return database.DB.Create(entry).Error
}

func (r *EntryRepository) FindAll() ([]models.EntryLog, error) {
	var entries []models.EntryLog

	err := database.DB.
		Preload("Visitor").
		Preload("User").
		Order("entered_at DESC").
		Find(&entries).Error

	return entries, err
}

func (r *EntryRepository) GetVisitorEntriesByStatusAndDate(
	status string,
	filter string,
	from string,
	to string,
) ([]models.EntryLog, error) {
	var entries []models.EntryLog

	query := database.DB.
		Preload("Visitor").
		Preload("Visitor.Documents").
		Model(&models.EntryLog{}).
		Order("entered_at DESC")
	today := time.Now()
	start := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, today.Location())
	end := start.AddDate(0, 0, 1)

	query = query.Where("entered_at >= ? AND entered_at < ?", start, end)

	if status == "active" {
		query = query.Where("status = ?", "active")
	}

	if status == "exited" {
		query = query.Where("status = ?", "exited")
	}

	err := query.Find(&entries).Error
	return entries, err
}

func (r *EntryRepository) GetVisitorStats() (map[string]int64, error) {

	var totalActiveVisitors int64
	var todayActiveVisitors int64
	var todayExitedVisitors int64

	today := time.Now().Format("2006-01-02")

	err := database.DB.
		Model(&models.EntryLog{}).
		Where("exited_at IS NULL").
		Count(&totalActiveVisitors).Error

	if err != nil {
		return nil, err
	}

	err = database.DB.
		Model(&models.EntryLog{}).
		Where("DATE(entered_at) = ?", today).
		Where("exited_at IS NULL").
		Count(&todayActiveVisitors).Error

	if err != nil {
		return nil, err
	}

	err = database.DB.
		Model(&models.EntryLog{}).
		Where("DATE(exited_at) = ?", today).
		Count(&todayExitedVisitors).Error

	if err != nil {
		return nil, err
	}

	return map[string]int64{
		"active_visitors": totalActiveVisitors,
		"entered_today":   todayActiveVisitors,
		"exited_today":    todayExitedVisitors,
	}, nil
}

func (r *EntryRepository) ExitVisitor(entryID uint) error {
	now := time.Now()

	result := database.DB.
		Model(&models.EntryLog{}).
		Where("id = ? AND exited_at IS NULL", entryID).
		Updates(map[string]any{
			"exited_at": now,
			"status":    "exited",
		})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("entry not found or already exited")
	}

	return nil
}

func (r *EntryRepository) GetActiveEntries() ([]dto.VisitorGroupedDTO, error) {

	var entries []models.EntryLog

	err := database.DB.
		Preload("Visitor").
		Preload("Visitor.Documents").
		Where("exited_at IS NULL").
		Order("entered_at DESC").
		Find(&entries).Error

	if err != nil {
		return nil, err
	}

	ist, _ := time.LoadLocation("Asia/Kolkata")

	groupMap := make(map[string][]dto.VisitorListItemDTO)

	for _, entry := range entries {

		date := entry.EnteredAt.
			In(ist).
			Format("2006-01-02")

		statusText := entry.Status
		if statusText == "" {
			statusText = "active"
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

		groupMap[date] = append(
			groupMap[date],
			item,
		)
	}

	var dates []string

	for date := range groupMap {
		dates = append(dates, date)
	}

	sort.Sort(sort.Reverse(sort.StringSlice(dates)))

	var result []dto.VisitorGroupedDTO

	for _, date := range dates {
		result = append(
			result,
			dto.VisitorGroupedDTO{
				Date:     date,
				Visitors: groupMap[date],
			},
		)
	}

	return result, nil
}

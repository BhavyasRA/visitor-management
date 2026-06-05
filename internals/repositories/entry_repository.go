package repositories

import (
	"time"

	"entry-system/internals/database"
	"entry-system/internals/models"
)

type EntryRepository struct{}

func NewEntryRepository() *EntryRepository {
	return &EntryRepository{}
}

func (r *EntryRepository) Create(entry *models.EntryLog) error {
	return database.DB.Create(entry).Error
}

func (r *EntryRepository) ExitVisitor(entryID uint) error {
	now := time.Now()

	return database.DB.
		Model(&models.EntryLog{}).
		Where("id = ? AND exited_at IS NULL", entryID).
		Updates(map[string]any{
			"exited_at": now,
			"status":    "exited",
		}).Error
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

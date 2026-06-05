package repositories

import (
	"errors"
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

func (r *EntryRepository) FindAll() ([]models.EntryLog, error) {
	var entries []models.EntryLog

	err := database.DB.
		Preload("Visitor").
		Preload("User").
		Order("entered_at DESC").
		Find(&entries).Error

	return entries, err
}

func (r *EntryRepository) ExitVisitor(visitorID uint) error {
	now := time.Now()

	return database.DB.
		Model(&models.EntryLog{}).
		Where("visitor_id = ? AND exited_at IS NULL", visitorID).
		Updates(map[string]interface{}{
			"exited_at": now,
			"status":    "exited",
		}).Error
}

func (r *EntryRepository) GetVisitorEntriesByStatusAndDate(
	status string,
	filter string,
	from string,
	to string,
) ([]models.EntryLog, error) {

	var entries []models.EntryLog

	query := database.DB.
		Model(&models.EntryLog{}).
		Preload("Visitor").
		Order("entered_at DESC")

	if status == "active" {
		query = query.Where("exited_at IS NULL")
	}

	if status == "exited" {
		query = query.Where("exited_at IS NOT NULL")
	}

	now := time.Now()

	switch filter {

	case "today":
		start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		end := start.AddDate(0, 0, 1)

		query = query.Where("entered_at >= ? AND entered_at < ?", start, end)

	case "yesterday":
		todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		yesterdayStart := todayStart.AddDate(0, 0, -1)

		query = query.Where("entered_at >= ? AND entered_at < ?", yesterdayStart, todayStart)

	case "last_2_days":
		start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).
			AddDate(0, 0, -2)

		query = query.Where("entered_at >= ?", start)

	case "week":
		start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).
			AddDate(0, 0, -7)

		query = query.Where("entered_at >= ?", start)

	case "month":
		start := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		end := start.AddDate(0, 1, 0)

		query = query.Where("entered_at >= ? AND entered_at < ?", start, end)

	case "custom":
		if from == "" || to == "" {
			return nil, errors.New("from and to are required")
		}

		start, err := time.Parse("2006-01-02", from)
		if err != nil {
			return nil, errors.New("invalid from date")
		}

		end, err := time.Parse("2006-01-02", to)
		if err != nil {
			return nil, errors.New("invalid to date")
		}

		end = end.AddDate(0, 0, 1)

		query = query.Where("entered_at >= ? AND entered_at < ?", start, end)

	case "", "all":

	default:
		return nil, errors.New("invalid filter")
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

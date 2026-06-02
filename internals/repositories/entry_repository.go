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
		Order("created_at DESC").
		Find(&entries).Error

	return entries, err
}

func (r *EntryRepository) ExitVisitor(visitorID uint) error {
	now := time.Now()

	return database.DB.
		Model(&models.EntryLog{}).
		Where("visitor_id = ? AND exited_at IS NULL", visitorID).
		Update("exited_at", now).Error
}

func (r *EntryRepository) FindEntriesByFilter(
	filter string,
	from string,
	to string,
) ([]models.EntryLog, error) {

	var entries []models.EntryLog

	query := database.DB.
		Model(&models.EntryLog{}).
		Order("created_at DESC")

	now := time.Now()

	switch filter {

	case "today":
		start := time.Date(
			now.Year(),
			now.Month(),
			now.Day(),
			0, 0, 0, 0,
			now.Location(),
		)

		end := start.AddDate(0, 0, 1)

		query = query.Where(
			"created_at >= ? AND created_at < ?",
			start,
			end,
		)

	case "yesterday":
		todayStart := time.Date(
			now.Year(),
			now.Month(),
			now.Day(),
			0, 0, 0, 0,
			now.Location(),
		)

		yesterdayStart := todayStart.AddDate(0, 0, -1)

		query = query.Where(
			"created_at >= ? AND created_at < ?",
			yesterdayStart,
			todayStart,
		)

	case "last_2_days":
		start := now.AddDate(0, 0, -2)

		query = query.Where(
			"created_at >= ?",
			start,
		)

	case "week":
		start := now.AddDate(0, 0, -7)

		query = query.Where(
			"created_at >= ?",
			start,
		)

	case "month":
		start := now.AddDate(0, -1, 0)

		query = query.Where(
			"created_at >= ?",
			start,
		)

	case "custom":
		if from == "" || to == "" {
			return nil, errors.New("from and to dates are required for custom filter")
		}

		start, err := time.Parse("2006-01-02", from)
		if err != nil {
			return nil, errors.New("invalid from date format, use YYYY-MM-DD")
		}

		end, err := time.Parse("2006-01-02", to)
		if err != nil {
			return nil, errors.New("invalid to date format, use YYYY-MM-DD")
		}

		end = end.AddDate(0, 0, 1)

		query = query.Where(
			"created_at >= ? AND created_at < ?",
			start,
			end,
		)

	case "all":
	
	default:
		return nil, errors.New("invalid filter")
	}

	err := query.Find(&entries).Error

	return entries, err
}

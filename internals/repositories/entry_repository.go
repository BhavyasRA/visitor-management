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
	query := `
		INSERT INTO entry_logs (
			visitor_id,
			person_to_meet,
			purpose,
			status,
			visiting_till,
			entered_at,
			exited_at,
			created_at,
			updated_at
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, NOW(), NOW())
	`

	return database.DB.Exec(
		query,
		entry.VisitorID,
		entry.PersonToMeet,
		entry.Purpose,
		entry.Status,
		entry.VisitingTill,
		entry.EnteredAt,
		entry.ExitedAt,
	).Error
}

func (r *EntryRepository) FindAll() ([]models.EntryLog, error) {
	var entries []models.EntryLog

	query := `
		SELECT *
		FROM entry_logs
		ORDER BY entered_at DESC
	`

	err := database.DB.Raw(query).Scan(&entries).Error
	if err != nil {
		return nil, err
	}

	r.loadVisitor(&entries)

	return entries, nil
}

func (r *EntryRepository) GetVisitorEntriesByStatusAndDate(
	status string,
	filter string,
	from string,
	to string,
) ([]models.EntryLog, error) {

	var entries []models.EntryLog

	today := time.Now()

	start := time.Date(
		today.Year(),
		today.Month(),
		today.Day(),
		0, 0, 0, 0,
		today.Location(),
	)

	end := start.AddDate(0, 0, 1)

	query := `
		SELECT *
		FROM entry_logs
		WHERE entered_at >= ?
		AND entered_at < ?
		AND exited_at IS NULL
		ORDER BY entered_at DESC
	`

	err := database.DB.Raw(query, start, end).Scan(&entries).Error
	if err != nil {
		return nil, err
	}

	r.loadVisitorAndDocuments(&entries)

	return entries, nil
}

func (r *EntryRepository) GetVisitorStats() (map[string]int64, error) {
	var totalActiveVisitors int64
	var todayActiveVisitors int64
	var todayExitedVisitors int64

	today := time.Now().Format("2006-01-02")

	err := database.DB.Raw(`
		SELECT COUNT(*)
		FROM entry_logs
		WHERE exited_at IS NULL
	`).Scan(&totalActiveVisitors).Error

	if err != nil {
		return nil, err
	}

	err = database.DB.Raw(`
		SELECT COUNT(*)
		FROM entry_logs
		WHERE DATE(entered_at) = ?
		AND exited_at IS NULL
	`, today).Scan(&todayActiveVisitors).Error

	if err != nil {
		return nil, err
	}

	err = database.DB.Raw(`
		SELECT COUNT(*)
		FROM entry_logs
		WHERE DATE(exited_at) = ?
	`, today).Scan(&todayExitedVisitors).Error

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

	result := database.DB.Exec(`
		UPDATE entry_logs
		SET
			exited_at = ?,
			status = ?,
			updated_at = NOW()
		WHERE id = ?
		AND exited_at IS NULL
	`, now, "exited", entryID)

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

	err := database.DB.Raw(`
		SELECT *
		FROM entry_logs
		WHERE exited_at IS NULL
		ORDER BY entered_at DESC
	`).Scan(&entries).Error

	if err != nil {
		return nil, err
	}

	r.loadVisitorAndDocuments(&entries)

	ist, _ := time.LoadLocation("Asia/Kolkata")

	groupMap := make(map[string][]dto.VisitorListItemDTO)

	for _, entry := range entries {
		date := entry.EnteredAt.In(ist).Format("2006-01-02")

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

func (r *EntryRepository) FindActiveEntryByVisitorID(
	visitorID uint,
) (*models.EntryLog, error) {

	var entry models.EntryLog

	err := database.DB.Raw(`
		SELECT *
		FROM entry_logs
		WHERE visitor_id = ?
		AND exited_at IS NULL
		AND status = ?
		LIMIT 1
	`, visitorID, "active").Scan(&entry).Error

	if err != nil {
		return nil, err
	}

	if entry.ID == 0 {
		return nil, errors.New("active entry not found")
	}

	return &entry, nil
}

func (r *EntryRepository) loadVisitor(entries *[]models.EntryLog) {
	for i := range *entries {
		database.DB.Raw(`
			SELECT *
			FROM visitors
			WHERE id = ?
			LIMIT 1
		`, (*entries)[i].VisitorID).Scan(&(*entries)[i].Visitor)
	}
}

func (r *EntryRepository) loadVisitorAndDocuments(entries *[]models.EntryLog) {
	for i := range *entries {
		database.DB.Raw(`
			SELECT *
			FROM visitors
			WHERE id = ?
			LIMIT 1
		`, (*entries)[i].VisitorID).Scan(&(*entries)[i].Visitor)

		database.DB.Raw(`
			SELECT *
			FROM visitor_documents
			WHERE visitor_id = ?
			ORDER BY id ASC
		`, (*entries)[i].VisitorID).Scan(&(*entries)[i].Visitor.Documents)
	}
}

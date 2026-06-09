package repositories

import (
	"errors"

	"entry-system/internals/database"
	"entry-system/internals/dto"
	"entry-system/internals/models"
)

type VisitorRepository struct{}

func NewVisitorRepository() *VisitorRepository {
	return &VisitorRepository{}
}

func (r *VisitorRepository) Create(visitor *models.Visitor) error {
	query := `
		INSERT INTO visitors (
			name,
			mobile,
			email,
			is_restricted,
			created_at,
			updated_at
		)
		VALUES (?, ?, ?, ?, NOW(), NOW())
	`

	return database.DB.Exec(
		query,
		visitor.Name,
		visitor.Mobile,
		visitor.Email,
		visitor.IsRestricted,
	).Error
}

func (r *VisitorRepository) Update(visitor *models.Visitor) error {
	query := `
		UPDATE visitors
		SET
			name = ?,
			mobile = ?,
			email = ?,
			is_restricted = ?,
			updated_at = NOW()
		WHERE id = ?
	`

	result := database.DB.Exec(
		query,
		visitor.Name,
		visitor.Mobile,
		visitor.Email,
		visitor.IsRestricted,
		visitor.ID,
	)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("visitor not found")
	}

	return nil
}

func (r *VisitorRepository) FindByID(id uint) (*models.Visitor, error) {
	var visitor models.Visitor

	query := `
		SELECT *
		FROM visitors
		WHERE id = ?
		LIMIT 1
	`

	err := database.DB.Raw(query, id).Scan(&visitor).Error
	if err != nil {
		return nil, err
	}

	if visitor.ID == 0 {
		return nil, errors.New("visitor not found")
	}

	r.loadDocuments(&visitor)

	return &visitor, nil
}

func (r *VisitorRepository) FindByMobile(mobile string) (*models.Visitor, error) {
	var visitor models.Visitor

	query := `
		SELECT *
		FROM visitors
		WHERE mobile = ?
		LIMIT 1
	`

	err := database.DB.Raw(query, mobile).Scan(&visitor).Error
	if err != nil {
		return nil, err
	}

	if visitor.ID == 0 {
		return nil, errors.New("visitor not found")
	}

	r.loadDocuments(&visitor)

	return &visitor, nil
}

func (r *VisitorRepository) FindAllWithFilters(filter dto.VisitorFilter) ([]models.Visitor, error) {
	var visitors []models.Visitor

	query := `
		SELECT *
		FROM visitors
		WHERE 1 = 1
	`

	args := []any{}

	if filter.Name != "" {
		query += ` AND name ILIKE ?`
		args = append(args, "%"+filter.Name+"%")
	}

	if filter.Mobile != "" {
		query += ` AND mobile = ?`
		args = append(args, filter.Mobile)
	}

	if filter.Email != "" {
		query += ` AND email ILIKE ?`
		args = append(args, "%"+filter.Email+"%")
	}

	query += ` ORDER BY created_at DESC`

	err := database.DB.Raw(query, args...).Scan(&visitors).Error
	if err != nil {
		return nil, err
	}

	for i := range visitors {
		r.loadDocuments(&visitors[i])
	}

	return visitors, nil
}

func (r *VisitorRepository) Restrict(id uint) error {
	query := `
		UPDATE visitors
		SET
			is_restricted = true,
			updated_at = NOW()
		WHERE id = ?
	`

	result := database.DB.Exec(query, id)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("visitor not found")
	}

	return nil
}

func (r *VisitorRepository) VisitorHistory(userID uint) ([]models.Visitor, error) {
	var visitors []models.Visitor

	query := `
		SELECT v.*
		FROM visitors v
		JOIN entry_logs e ON e.visitor_id = v.id
		WHERE e.person_to_meet = ?
		ORDER BY v.created_at DESC
	`

	err := database.DB.Raw(query, userID).Scan(&visitors).Error
	if err != nil {
		return nil, err
	}

	return visitors, nil
}

func (r *VisitorRepository) loadDocuments(visitor *models.Visitor) {
	query := `
		SELECT *
		FROM visitor_documents
		WHERE visitor_id = ?
		ORDER BY created_at DESC
	`

	database.DB.Raw(query, visitor.ID).Scan(&visitor.Documents)
}

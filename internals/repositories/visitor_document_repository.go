package repositories

import (
	"errors"

	"entry-system/internals/database"
	"entry-system/internals/models"
)

type VisitorDocumentRepository struct{}

func NewVisitorDocumentRepository() *VisitorDocumentRepository {
	return &VisitorDocumentRepository{}
}

func (r *VisitorDocumentRepository) Create(
	document *models.VisitorDocument,
) error {

	query := `
		INSERT INTO visitor_documents (
			visitor_id,
			photo_url,
			identity_document_url,
			document_type,
			document_number,
			created_at,
			updated_at
		)
		VALUES (?, ?, ?, ?, ?, NOW(), NOW())
	`

	return database.DB.Exec(
		query,
		document.VisitorID,
		document.PhotoURL,
		document.IdentityDocumentURL,
		document.DocumentType,
		document.DocumentNumber,
	).Error
}

func (r *VisitorDocumentRepository) FindByID(
	id uint,
) (*models.VisitorDocument, error) {

	var document models.VisitorDocument

	query := `
		SELECT *
		FROM visitor_documents
		WHERE id = ?
		LIMIT 1
	`

	err := database.DB.Raw(query, id).Scan(&document).Error
	if err != nil {
		return nil, err
	}

	if document.ID == 0 {
		return nil, errors.New("visitor document not found")
	}

	return &document, nil
}

func (r *VisitorDocumentRepository) Update(
	document *models.VisitorDocument,
) error {

	query := `
		UPDATE visitor_documents
		SET
			visitor_id = ?,
			photo_url = ?,
			identity_document_url = ?,
			document_type = ?,
			document_number = ?,
			updated_at = NOW()
		WHERE id = ?
	`

	result := database.DB.Exec(
		query,
		document.VisitorID,
		document.PhotoURL,
		document.IdentityDocumentURL,
		document.DocumentType,
		document.DocumentNumber,
		document.ID,
	)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("visitor document not found")
	}

	return nil
}
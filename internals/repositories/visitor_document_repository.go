package repositories

import (
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
	return database.DB.Create(document).Error
}

func (r *VisitorDocumentRepository) FindByID(
	id uint,
) (*models.VisitorDocument, error) {

	var document models.VisitorDocument

	err := database.DB.
		Where("id = ?", id).
		First(&document).Error

	return &document, err
}

func (r *VisitorDocumentRepository) Update(
	document *models.VisitorDocument,
) error {

	return database.DB.
		Save(document).
		Error
}

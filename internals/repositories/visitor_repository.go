package repositories

import (
	"entry-system/internals/database"
	"entry-system/internals/models"
)

type VisitorRepository struct{}

func NewVisitorRepository() *VisitorRepository {
	return &VisitorRepository{}
}

func (r *VisitorRepository) Create(visitor *models.Visitor) error {
	return database.DB.Create(visitor).Error
}

func (r *VisitorRepository) FindByID(id uint) (*models.Visitor, error) {
	var visitor models.Visitor

	err := database.DB.First(&visitor, id).Error

	return &visitor, err
}

func (r *VisitorRepository) FindAll() ([]models.Visitor, error) {
	var visitors []models.Visitor

	err := database.DB.Find(&visitors).Error

	return visitors, err
}

func (r *VisitorRepository) Update(visitor *models.Visitor) error {
	return database.DB.Save(visitor).Error
}

func (r *VisitorRepository) Restrict(id uint) error {
	return database.DB.
		Model(&models.Visitor{}).
		Where("id = ?", id).
		Update("is_restricted", true).Error
}

func (r *VisitorRepository) VisitorHistory(userID uint) ([]models.Visitor, error) {
	var visitors []models.Visitor

	err := database.DB.
		Where("to_whom = ?", userID).
		Order("created_at DESC").
		Find(&visitors).Error

	return visitors, err
}

func (r *VisitorRepository) GetVisitorsByUserID(userID uint) ([]models.Visitor, error) {
	var visitors []models.Visitor

	err := database.DB.
		Where("to_whom = ?", userID).
		Order("created_at DESC").
		Find(&visitors).Error

	return visitors, err
}

package repositories

import (
	"entry-system/internals/database"
	"entry-system/internals/dto"
	"entry-system/internals/models"
)

type VisitorRepository struct{}

func NewVisitorRepository() *VisitorRepository {
	return &VisitorRepository{}
}

func (r *VisitorRepository) Create(visitor *models.Visitor) error {
	return database.DB.Create(visitor).Error
}

func (r *VisitorRepository) Update(visitor *models.Visitor) error {
	return database.DB.Save(visitor).Error
}

func (r *VisitorRepository) FindByID(id uint) (*models.Visitor, error) {
	var visitor models.Visitor

	err := database.DB.
		Preload("Documents").
		First(&visitor, id).Error

	return &visitor, err
}

func (r *VisitorRepository) FindByMobile(mobile string) (*models.Visitor, error) {
	var visitor models.Visitor

	err := database.DB.
		Preload("Documents").
		Where("mobile = ?", mobile).
		First(&visitor).Error

	return &visitor, err
}

func (r *VisitorRepository) FindAllWithFilters(
	filter dto.VisitorFilter,
) ([]models.Visitor, error) {
	var visitors []models.Visitor

	query := database.DB.
		Preload("Documents").
		Model(&models.Visitor{})

	if filter.Name != "" {
		query = query.Where("name ILIKE ?", "%"+filter.Name+"%")
	}

	if filter.Mobile != "" {
		query = query.Where("mobile = ?", filter.Mobile)
	}

	if filter.Email != "" {
		query = query.Where("email ILIKE ?", "%"+filter.Email+"%")
	}

	err := query.
		Order("created_at DESC").
		Find(&visitors).Error

	return visitors, err
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
		Where("person_to_meet = ?", userID).
		Order("created_at DESC").Find(&visitors).Error
	return visitors, err
}

package repositories

import (
	"entry-system/internals/database"
	"entry-system/internals/models"
)

type AuthRepository struct{}

func NewAuthRepository() *AuthRepository {
	return &AuthRepository{}
}

func (r *AuthRepository) Create(auth *models.Authentication) error {
	return database.DB.Create(auth).Error
}

func (r *AuthRepository) FindByUserID(userID uint) (*models.Authentication, error) {
	var auth models.Authentication

	err := database.DB.
		Where("user_id = ?", userID).
		First(&auth).Error

	return &auth, err
}

func (r *AuthRepository) Update(auth *models.Authentication) error {
	return database.DB.Save(auth).Error
}

func (r *AuthRepository) FindByVerificationToken(token string) (*models.Authentication, error) {
	var auth models.Authentication

	err := database.DB.
		Where("verification_token = ?", token).
		First(&auth).Error

	return &auth, err
}

func (r *AuthRepository) FindByResetToken(token string) (*models.Authentication, error) {
	var auth models.Authentication

	err := database.DB.
		Where("reset_token = ?", token).
		First(&auth).Error

	return &auth, err
}
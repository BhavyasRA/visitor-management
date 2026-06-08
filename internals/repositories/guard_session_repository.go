package repositories

import (
	"errors"
	"time"

	"entry-system/internals/database"
	"entry-system/internals/models"
)

type GuardSessionRepository struct{}

func NewGuardSessionRepository() *GuardSessionRepository {
	return &GuardSessionRepository{}
}

func (r *GuardSessionRepository) Create(
	session *models.GuardSession,
) error {
	return database.DB.Create(session).Error
}

func (r *GuardSessionRepository) FindActiveByGuardID(
	guardID uint,
) (*models.GuardSession, error) {

	var session models.GuardSession

	err := database.DB.
		Where("guard_id = ? AND status = ?", guardID, "active").
		Order("created_at DESC").
		First(&session).Error

	return &session, err
}

func (r *GuardSessionRepository) Logout(
	guardID uint,
	logoutPhotoURL string,
) error {

	now := time.Now()

	result := database.DB.
		Model(&models.GuardSession{}).
		Where("guard_id = ? AND status = ?", guardID, "active").
		Updates(map[string]any{
			"logout_photo_url": logoutPhotoURL,
			"logout_at":        now,
			"status":           "completed",
		})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("no active guard session found")
	}

	return nil
}

func (r *GuardSessionRepository) FindAll() ([]models.GuardSession, error) {
	var sessions []models.GuardSession

	err := database.DB.
		Preload("Guard").
		Order("created_at DESC").
		Find(&sessions).Error

	return sessions, err
}

func (r *GuardSessionRepository) FindByGuardID(
	guardID uint,
) ([]models.GuardSession, error) {

	var sessions []models.GuardSession

	err := database.DB.
		Preload("Guard").
		Where("guard_id = ?", guardID).
		Order("created_at DESC").
		Find(&sessions).Error

	return sessions, err
}
func (r *GuardSessionRepository) Update(
	session *models.GuardSession,
) error {

	return database.DB.
		Save(session).
		Error
}

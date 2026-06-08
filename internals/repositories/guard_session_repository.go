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

	var session models.GuardSession

	err := database.DB.
		Where("guard_id = ? AND logout_at IS NULL", guardID).
		Order("login_at DESC").
		First(&session).
		Error

	if err != nil {
		return err
	}

	now := time.Now()

	session.LogoutPhotoURL = logoutPhotoURL
	session.LogoutAt = &now
	session.Status = "logged_out"

	return database.DB.Save(&session).Error
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

func (r *GuardSessionRepository) CheckoutBySessionID(
	sessionID string,
	logoutPhotoURL string,
) error {

	var session models.GuardSession

	err := database.DB.
		Where("id = ?", sessionID).
		First(&session).
		Error

	if err != nil {
		return errors.New("guard session not found")
	}

	if session.LogoutAt != nil {
		return errors.New("guard session already checked out")
	}

	now := time.Now()

	session.LogoutAt = &now
	session.LogoutPhotoURL = logoutPhotoURL
	session.Status = "checked_out"

	return database.DB.Save(&session).Error
}

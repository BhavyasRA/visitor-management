package repositories

import (
	"errors"
	"time"

	"entry-system/internals/database"
	"entry-system/internals/dto"
	"entry-system/internals/models"
)

type GuardSessionRepository struct{}

func NewGuardSessionRepository() *GuardSessionRepository {
	return &GuardSessionRepository{}
}

func (r *GuardSessionRepository) Create(
	session *models.GuardSession,
) error {

	query := `
		INSERT INTO guard_sessions (
			guard_id,login_photo_url,
			logout_photo_url,login_at,
			logout_at,status,
			created_at,updated_at
		)
		VALUES (?, ?, ?, ?, ?, ?, NOW(), NOW())
	`

	return database.DB.Exec(
		query,
		session.GuardID,
		session.LoginPhotoURL,
		session.LogoutPhotoURL,
		session.LoginAt,
		session.LogoutAt,
		session.Status,
	).Error
}

func (r *GuardSessionRepository) FindActiveByGuardID(
	guardID uint,
) (*models.GuardSession, error) {

	var session models.GuardSession

	query := `
		SELECT
			id,
			guard_id,
			login_photo_url,
			logout_photo_url,
			login_at,
			logout_at,
			status,
			created_at,
			updated_at
		FROM guard_sessions
		WHERE guard_id = ?
		AND status = ?
		ORDER BY created_at DESC
		LIMIT 1
	`

	err := database.DB.Raw(query, guardID, "active").Scan(&session).Error
	if err != nil {
		return nil, err
	}

	if session.ID == 0 {
		return nil, errors.New("active guard session not found")
	}

	return &session, nil
}

func (r *GuardSessionRepository) Logout(
	guardID uint,
	logoutPhotoURL string,
) error {

	var session models.GuardSession

	findQuery := `
		SELECT
			id,
			guard_id,
			login_photo_url,
			logout_photo_url,
			login_at,
			logout_at,
			status,
			created_at,
			updated_at
		FROM guard_sessions
		WHERE guard_id = ?
		AND logout_at IS NULL
		ORDER BY login_at DESC
		LIMIT 1
	`

	err := database.DB.Raw(findQuery, guardID).Scan(&session).Error
	if err != nil {
		return err
	}

	if session.ID == 0 {
		return errors.New("guard session not found")
	}

	now := time.Now()

	updateQuery := `
		UPDATE guard_sessions
		SET
			logout_photo_url = ?,
			logout_at = ?,
			status = ?,
			updated_at = NOW()
		WHERE id = ?
	`

	result := database.DB.Exec(
		updateQuery,
		logoutPhotoURL,
		now,
		"logged_out",
		session.ID,
	)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("guard session not updated")
	}

	return nil
}

func (r *GuardSessionRepository) FindAll() ([]models.GuardSession, error) {

	var sessions []models.GuardSession

	query := `
		SELECT
			gs.id,
			gs.guard_id,
			gs.login_photo_url,
			gs.logout_photo_url,
			gs.login_at,
			gs.logout_at,
			gs.status,
			gs.created_at,
			gs.updated_at
		FROM guard_sessions gs
		ORDER BY gs.created_at DESC
	`

	err := database.DB.Raw(query).Scan(&sessions).Error
	if err != nil {
		return nil, err
	}

	return sessions, nil
}

func (r *GuardSessionRepository) FindByGuardID(guardID uint) (*dto.GuardSessionsByGuardDTO, error) {

	var result dto.GuardSessionsByGuardDTO

	guardQuery := `
		SELECT
			id AS guard_id,
			name,
			email,
			phone
		FROM users
		WHERE id = ?
		LIMIT 1
	`

	err := database.DB.Raw(guardQuery, guardID).Scan(&result).Error
	if err != nil {
		return nil, err
	}

	if result.GuardID == 0 {
		return nil, errors.New("guard not found")
	}

	var sessions []dto.GuardSessionItemDTO

	sessionQuery := `
		SELECT
			id,
			login_photo_url,
			logout_photo_url,
			login_at,
			logout_at,
			status,
			created_at,
			updated_at
		FROM guard_sessions
		WHERE guard_id = ?
		ORDER BY created_at DESC
	`

	err = database.DB.Raw(sessionQuery, guardID).Scan(&sessions).Error
	if err != nil {
		return nil, err
	}

	result.Sessions = sessions

	return &result, nil
}

func (r *GuardSessionRepository) Update(
	session *models.GuardSession,
) error {

	query := `
		UPDATE guard_sessions
		SET
			guard_id = ?,
			login_photo_url = ?,
			logout_photo_url = ?,
			login_at = ?,
			logout_at = ?,
			status = ?,
			updated_at = NOW()
		WHERE id = ?
	`

	result := database.DB.Exec(
		query,
		session.GuardID,
		session.LoginPhotoURL,
		session.LogoutPhotoURL,
		session.LoginAt,
		session.LogoutAt,
		session.Status,
		session.ID,
	)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("guard session not found")
	}

	return nil
}

func (r *GuardSessionRepository) CheckoutBySessionID(
	sessionID string,
	logoutPhotoURL string,
) error {

	var session models.GuardSession

	findQuery := `
		SELECT
			id,
			guard_id,
			login_photo_url,
			logout_photo_url,
			login_at,
			logout_at,
			status,
			created_at,
			updated_at
		FROM guard_sessions
		WHERE id = ?
		LIMIT 1
	`

	err := database.DB.Raw(findQuery, sessionID).Scan(&session).Error
	if err != nil {
		return err
	}

	if session.ID == 0 {
		return errors.New("guard session not found")
	}

	if session.LogoutAt != nil {
		return errors.New("guard session already checked out")
	}

	now := time.Now()

	updateQuery := `
		UPDATE guard_sessions
		SET
			logout_at = ?,
			logout_photo_url = ?,
			status = ?,
			updated_at = NOW()
		WHERE id = ?
	`

	result := database.DB.Exec(
		updateQuery,
		now,
		logoutPhotoURL,
		"checked_out",
		session.ID,
	)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("guard session not updated")
	}

	return nil
}

func (r *GuardSessionRepository) GetGuardInfoByUserID(
	userID uint,
) (*dto.GuardSessionInfoDTO, error) {

	var data dto.GuardSessionInfoDTO

	query := `
		SELECT
			u.id AS user_id,
			u.name AS name,
			gs.login_photo_url AS login_photo_url,
			gs.login_at AS entered_at
		FROM users u
		JOIN guard_sessions gs ON gs.guard_id = u.id
		WHERE u.id = ?
		ORDER BY gs.login_at DESC
		LIMIT 1
	`

	err := database.DB.Raw(query, userID).Scan(&data).Error
	if err != nil {
		return nil, err
	}

	if data.UserID == 0 {
		return nil, errors.New("guard session not found")
	}

	loc, err := time.LoadLocation("Asia/Kolkata")
	if err == nil {
		data.EnteredAt = data.EnteredAt.In(loc)
	}

	return &data, nil
}

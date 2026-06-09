package repositories

import (
	"database/sql"
	"errors"

	"entry-system/internals/database"
	"entry-system/internals/models"
)

type AuthRepository struct{}

func NewAuthRepository() *AuthRepository {
	return &AuthRepository{}
}

func (r *AuthRepository) Create(auth *models.Authentication) error {
	query := `
		INSERT INTO authentications (
			user_id,
			password,
			verification_token,
			reset_token,
			otp,
			token_expires_at,
			otp_expires_at,
			verified_at,
			created_at,
			updated_at
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW())
	`

	err := database.DB.Exec(
		query,
		auth.UserID,
		auth.Password,
		nullableString(auth.VerificationToken),
		nullableString(auth.ResetToken),
		nullableString(auth.OTP),
		auth.TokenExpiresAt,
		auth.OTPExpiresAt,
		auth.VerifiedAt,
	).Error

	return err
}

func (r *AuthRepository) FindByUserID(userID uint) (*models.Authentication, error) {
	var auth models.Authentication

	query := `
		SELECT
			id,
			user_id,
			password,
			COALESCE(verification_token, ''),
			COALESCE(reset_token, ''),
			COALESCE(otp, ''),
			token_expires_at,
			otp_expires_at,
			verified_at,
			created_at,
			updated_at
		FROM authentications
		WHERE user_id = ?
		LIMIT 1
	`

	err := database.DB.Raw(query, userID).Row().Scan(
		&auth.ID,
		&auth.UserID,
		&auth.Password,
		&auth.VerificationToken,
		&auth.ResetToken,
		&auth.OTP,
		&auth.TokenExpiresAt,
		&auth.OTPExpiresAt,
		&auth.VerifiedAt,
		&auth.CreatedAt,
		&auth.UpdatedAt,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, errors.New("auth record not found")
	}

	if err != nil {
		return nil, err
	}

	return &auth, nil
}

func (r *AuthRepository) Update(auth *models.Authentication) error {
	query := `
		UPDATE authentications
		SET
			user_id = ?,
			password = ?,
			verification_token = ?,
			reset_token = ?,
			otp = ?,
			token_expires_at = ?,
			otp_expires_at = ?,
			verified_at = ?,
			updated_at = NOW()
		WHERE id = ?
	`

	result := database.DB.Exec(
		query,
		auth.UserID,
		auth.Password,
		nullableString(auth.VerificationToken),
		nullableString(auth.ResetToken),
		nullableString(auth.OTP),
		auth.TokenExpiresAt,
		auth.OTPExpiresAt,
		auth.VerifiedAt,
		auth.ID,
	)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("auth record not found")
	}

	return nil
}

func (r *AuthRepository) FindByVerificationToken(token string) (*models.Authentication, error) {
	var auth models.Authentication

	query := `
		SELECT
			id,
			user_id,
			password,
			COALESCE(verification_token, ''),
			COALESCE(reset_token, ''),
			COALESCE(otp, ''),
			token_expires_at,
			otp_expires_at,
			verified_at,
			created_at,
			updated_at
		FROM authentications
		WHERE verification_token = ?
		LIMIT 1
	`

	err := database.DB.Raw(query, token).Row().Scan(
		&auth.ID,
		&auth.UserID,
		&auth.Password,
		&auth.VerificationToken,
		&auth.ResetToken,
		&auth.OTP,
		&auth.TokenExpiresAt,
		&auth.OTPExpiresAt,
		&auth.VerifiedAt,
		&auth.CreatedAt,
		&auth.UpdatedAt,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, errors.New("invalid verification token")
	}

	if err != nil {
		return nil, err
	}

	return &auth, nil
}

func (r *AuthRepository) FindByResetToken(token string) (*models.Authentication, error) {
	var auth models.Authentication

	query := `
		SELECT
			id,
			user_id,
			password,
			COALESCE(verification_token, ''),
			COALESCE(reset_token, ''),
			COALESCE(otp, ''),
			token_expires_at,
			otp_expires_at,
			verified_at,
			created_at,
			updated_at
		FROM authentications
		WHERE reset_token = ?
		LIMIT 1
	`

	err := database.DB.Raw(query, token).Row().Scan(
		&auth.ID,
		&auth.UserID,
		&auth.Password,
		&auth.VerificationToken,
		&auth.ResetToken,
		&auth.OTP,
		&auth.TokenExpiresAt,
		&auth.OTPExpiresAt,
		&auth.VerifiedAt,
		&auth.CreatedAt,
		&auth.UpdatedAt,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, errors.New("invalid reset token")
	}

	if err != nil {
		return nil, err
	}

	return &auth, nil
}

func nullableString(value string) any {
	if value == "" {
		return nil
	}

	return value
}

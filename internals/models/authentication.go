package models

import "time"

type Authentication struct {
	ID uint `gorm:"primaryKey"`

	UserID uint `gorm:"unique"`

	Password string

	VerificationToken string
	ResetToken        string

	OTP string

	TokenExpiresAt *time.Time
	OTPExpiresAt   *time.Time

	VerifiedAt *time.Time

	CreatedAt time.Time
	UpdatedAt time.Time
}
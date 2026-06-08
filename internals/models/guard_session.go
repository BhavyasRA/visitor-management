package models

import "time"

type GuardSession struct {
	ID uint `gorm:"primaryKey" json:"id"`

	GuardID uint `json:"guard_id"`
	Guard   User `gorm:"foreignKey:GuardID" json:"guard"`

	LoginPhotoURL  string `json:"login_photo_url"`
	LogoutPhotoURL string `json:"logout_photo_url"`

	LoginAt  time.Time  `json:"login_at"`
	LogoutAt *time.Time `json:"logout_at"`

	Status string `gorm:"default:'active'" json:"status"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

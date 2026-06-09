package dto

import "time"

type GuardSessionInfoDTO struct {
	UserID        uint      `json:"user_id"`
	Name          string    `json:"name"`
	LoginPhotoURL string    `json:"login_photo_url"`
	EnteredAt     time.Time `json:"entered_at"`
}

type GuardSessionItemDTO struct {
	ID             uint       `json:"id"`
	LoginPhotoURL  string     `json:"login_photo_url"`
	LogoutPhotoURL string     `json:"logout_photo_url"`
	LoginAt        time.Time  `json:"login_at"`
	LogoutAt       *time.Time `json:"logout_at"`
	Status         string     `json:"status"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

type GuardSessionsByGuardDTO struct {
	GuardID  uint                  `json:"guard_id"`
	Name     string                `json:"name"`
	Email    string                `json:"email"`
	Phone    string                `json:"phone"`
	Sessions []GuardSessionItemDTO `json:"sessions" gorm:"-"`
}

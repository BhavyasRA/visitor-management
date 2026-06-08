package services

import (
	"errors"
	"time"

	"entry-system/internals/models"
	"entry-system/internals/repositories"
)

type GuardSessionService struct {
	guardSessionRepo *repositories.GuardSessionRepository
	userRepo         *repositories.UserRepository
}

func NewGuardSessionService() *GuardSessionService {
	return &GuardSessionService{
		guardSessionRepo: repositories.NewGuardSessionRepository(),
		userRepo:         repositories.NewUserRepository(),
	}
}

func (s *GuardSessionService) StartGuardSession(
	guardID uint,
	loginPhotoURL string,
) (*models.GuardSession, error) {

	_, err := s.userRepo.FindByID(guardID)
	if err != nil {
		return nil, errors.New("guard not found")
	}

	session := models.GuardSession{
		GuardID:       guardID,
		LoginPhotoURL: loginPhotoURL,
		LoginAt:       time.Now(),
		Status:        "active",
	}

	if err := s.guardSessionRepo.Create(&session); err != nil {
		return nil, err
	}

	return &session, nil
}

// func (s *GuardSessionService) StartGuardSession(
// 	guardID uint,
// 	loginPhotoURL string,
// ) (*models.GuardSession, error) {

// 	activeSession, err := s.guardSessionRepo.FindActiveByGuardID(guardID)

// 	if err == nil && activeSession != nil {

// 		now := time.Now()

// 		activeSession.LogoutAt = &now
// 		activeSession.Status = "auto_closed"

// 		if err := s.guardSessionRepo.Update(activeSession); err != nil {
// 			return nil, err
// 		}
// 	}

// 	session := models.GuardSession{
// 		GuardID:       guardID,
// 		LoginPhotoURL: loginPhotoURL,
// 		LoginAt:       time.Now(),
// 		Status:        "active",
// 	}

// 	if err := s.guardSessionRepo.Create(&session); err != nil {
// 		return nil, err
// 	}

// 	return &session, nil
// }

func (s *GuardSessionService) EndGuardSession(
	guardID uint,
	logoutPhotoURL string,
) error {

	return s.guardSessionRepo.Logout(
		guardID,
		logoutPhotoURL,
	)
}

func (s *GuardSessionService) GetAllGuardSessions() ([]models.GuardSession, error) {
	return s.guardSessionRepo.FindAll()
}

func (s *GuardSessionService) GetGuardSessionsByGuardID(
	guardID uint,
) ([]models.GuardSession, error) {

	return s.guardSessionRepo.FindByGuardID(guardID)
}

func (s *GuardSessionService) CheckoutGuardSession(
	sessionID string,
	logoutPhotoURL string,
) error {

	return s.guardSessionRepo.CheckoutBySessionID(
		sessionID,
		logoutPhotoURL,
	)
}

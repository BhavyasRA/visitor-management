package services

import (
	"errors"
	"time"

	"entry-system/internals/models"
	"entry-system/internals/repositories"
	"entry-system/internals/utils"

	"github.com/google/uuid"
)

type AuthService struct {
	userRepo *repositories.UserRepository
	authRepo *repositories.AuthRepository
}

func NewAuthService() *AuthService {
	return &AuthService{
		userRepo: repositories.NewUserRepository(),
		authRepo: repositories.NewAuthRepository(),
	}
}

func (s *AuthService) Signup(
	name string,
	email string,
	phone string,
	password string,
) (string, error) {

	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return "", err
	}

	user := models.User{
		Name:     name,
		Email:    email,
		Phone:    phone,
		IsActive: true,
	}

	if err := s.userRepo.Create(&user); err != nil {
		return "", err
	}

	token := uuid.New().String()
	expires := time.Now().Add(30 * time.Minute)

	auth := models.Authentication{
		UserID:            user.ID,
		Password:          hashedPassword,
		VerificationToken: token,
		TokenExpiresAt:    &expires,
	}

	err = s.authRepo.Create(&auth)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *AuthService) Login(email, password string) (string, error) {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	if !user.IsActive {
		return "", errors.New("user is inactive")
	}

	auth, err := s.authRepo.FindByUserID(user.ID)
	if err != nil {
		return "", errors.New("auth record not found")
	}

	if auth.VerifiedAt == nil {
		return "", errors.New("account not verified")
	}

	if err := utils.ComparePassword(auth.Password, password); err != nil {
		return "", errors.New("invalid credentials")
	}

	return utils.GenerateToken(user.ID)
}

func (s *AuthService) VerifyAccount(token string) error {
	auth, err := s.authRepo.FindByVerificationToken(token)
	if err != nil {
		return errors.New("invalid token")
	}

	if auth.TokenExpiresAt != nil && auth.TokenExpiresAt.Before(time.Now()) {
		return errors.New("verification token expired")
	}

	now := time.Now()

	auth.VerifiedAt = &now
	auth.VerificationToken = ""
	auth.TokenExpiresAt = nil

	return s.authRepo.Update(auth)
}

func (s *AuthService) ForgotPassword(email string) (string, error) {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return "", errors.New("user not found")
	}

	auth, err := s.authRepo.FindByUserID(user.ID)
	if err != nil {
		return "", err
	}

	resetToken := uuid.New().String()
	expires := time.Now().Add(15 * time.Minute)

	auth.ResetToken = resetToken
	auth.TokenExpiresAt = &expires

	if err := s.authRepo.Update(auth); err != nil {
		return "", err
	}

	return resetToken, nil
}

func (s *AuthService) ResetPassword(token, newPassword string) error {
	auth, err := s.authRepo.FindByResetToken(token)
	if err != nil {
		return errors.New("invalid reset token")
	}

	if auth.TokenExpiresAt != nil && auth.TokenExpiresAt.Before(time.Now()) {
		return errors.New("reset token expired")
	}

	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return err
	}

	auth.Password = hashedPassword
	auth.ResetToken = ""
	auth.TokenExpiresAt = nil

	return s.authRepo.Update(auth)
}

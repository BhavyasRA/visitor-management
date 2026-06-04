package services

import (
	"entry-system/internals/models"
	"entry-system/internals/repositories"
)

type UserService struct {
	userRepo    *repositories.UserRepository
	visitorRepo *repositories.VisitorRepository
}

func NewUserService() *UserService {
	return &UserService{
		userRepo:    repositories.NewUserRepository(),
		visitorRepo: repositories.NewVisitorRepository(),
	}
}

func (s *UserService) CreateUser(user *models.User) error {
	user.IsActive = true
	return s.userRepo.Create(user)
}

func (s *UserService) GetUsers() ([]models.User, error) {
	return s.userRepo.FindAll()
}

func (s *UserService) GetUser(id uint) (*models.User, error) {
	return s.userRepo.FindByID(id)
}

func (s *UserService) UpdateUser(id uint, name, email, phone string) error {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return err
	}

	user.Name = name
	user.Email = email
	user.Phone = phone

	return s.userRepo.Update(user)
}

func (s *UserService) DeleteUser(id uint) error {
	return s.userRepo.Delete(id)
}

func (s *UserService) DeactivateUser(id uint) error {
	return s.userRepo.Deactivate(id)
}

func (s *UserService) VisitorHistory(userID uint) ([]models.Visitor, error) {
	return s.visitorRepo.VisitorHistory(userID)
}


func (s *UserService) GetUsersDropdown() ([]map[string]any, error) {
	return s.userRepo.GetUsersDropdown()
}

package repositories

import (
	"entry-system/internals/database"
	"entry-system/internals/models"
)

type UserRepository struct{}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

func (r *UserRepository) Create(user *models.User) error {
	return database.DB.Create(user).Error
}

func (r *UserRepository) FindByID(id uint) (*models.User, error) {
	var user models.User

	err := database.DB.
		Preload("Roles").
		Preload("Roles.Permissions").
		First(&user, id).Error

	return &user, err
}

func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User

	err := database.DB.
		Where("email = ?", email).
		First(&user).Error

	return &user, err
}

func (r *UserRepository) FindAll() ([]models.User, error) {
	var users []models.User

	err := database.DB.
		Preload("Roles").
		Find(&users).Error

	return users, err
}

func (r *UserRepository) Update(user *models.User) error {
	return database.DB.Save(user).Error
}

func (r *UserRepository) Delete(id uint) error {
	return database.DB.Delete(&models.User{}, id).Error
}

func (r *UserRepository) Deactivate(id uint) error {
	return database.DB.
		Model(&models.User{}).
		Where("id = ?", id).
		Update("is_active", false).Error
}

func (r *UserRepository) AssignRoleByName(userID uint, roleName string) error {
	var role models.Role

	err := database.DB.
		Where("name = ?", roleName).
		First(&role).Error

	if err != nil {
		return err
	}

	return database.DB.Exec(
		"INSERT INTO user_roles (user_id, role_id) VALUES (?, ?) ON CONFLICT DO NOTHING",
		userID,
		role.ID,
	).Error
}

func (r *UserRepository) GetUsersDropdown() ([]map[string]any, error) {
	var users []models.User

	err := database.DB.
		Where("is_active = ?", true).
		Order("name ASC").
		Find(&users).Error

	if err != nil {
		return nil, err
	}

	var result []map[string]any

	for _, user := range users {
		result = append(result, map[string]any{
			"id":   user.ID,
			"name": user.Name,
		})
	}

	return result, nil
}

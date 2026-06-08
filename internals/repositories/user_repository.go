package repositories

import (
	"errors"

	"entry-system/internals/database"
	"entry-system/internals/models"
)

type UserRepository struct{}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

func (r *UserRepository) Create(user *models.User) error {
query := `
	INSERT INTO users (
		name,
		email,
		phone,
		is_active,
		created_at,
		updated_at
	)
	VALUES (?, ?, ?, ?, NOW(), NOW())

		RETURNING id, created_at, updated_at
	`

	return database.DB.Raw(
		query,
		user.Name,
		user.Email,
		user.Phone,
		user.IsActive,
	).Scan(user).Error
}

func (r *UserRepository) FindByID(id uint) (*models.User, error) {
	var user models.User

	query := `
		SELECT *
		FROM users
		WHERE id = ?
		LIMIT 1
	`

	err := database.DB.Raw(query, id).Scan(&user).Error
	if err != nil {
		return nil, err
	}

	if user.ID == 0 {
		return nil, errors.New("user not found")
	}

	r.loadRolesAndPermissions(&user)

	return &user, nil
}

func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User

	query := `
		SELECT *
		FROM users
		WHERE email = ?
		LIMIT 1
	`

	err := database.DB.Raw(query, email).Scan(&user).Error
	if err != nil {
		return nil, err
	}

	if user.ID == 0 {
		return nil, errors.New("user not found")
	}

	return &user, nil
}

func (r *UserRepository) FindAll() ([]models.User, error) {
	var users []models.User

	query := `
		SELECT *
		FROM users
		ORDER BY id ASC
	`

	err := database.DB.Raw(query).Scan(&users).Error
	if err != nil {
		return nil, err
	}

	for i := range users {
		r.loadRoles(&users[i])
	}

	return users, nil
}

func (r *UserRepository) Update(user *models.User) error {
	query := `
		UPDATE users
		SET
			name = ?,
			email = ?,
			phone_no = ?,
			password = ?,
			is_active = ?,
			updated_at = NOW()
		WHERE id = ?
	`

	result := database.DB.Exec(
		query,
		user.Name,
		user.Email,
		user.Phone,
		user.IsActive,
		user.ID,
	)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("user not found")
	}

	return nil
}

func (r *UserRepository) Delete(id uint) error {
	query := `
		DELETE FROM users
		WHERE id = ?
	`

	result := database.DB.Exec(query, id)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("user not found")
	}

	return nil
}

func (r *UserRepository) Deactivate(id uint) error {
	query := `
		UPDATE users
		SET
			is_active = false,
			updated_at = NOW()
		WHERE id = ?
	`

	result := database.DB.Exec(query, id)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("user not found")
	}

	return nil
}

func (r *UserRepository) AssignRoleByName(userID uint, roleName string) error {
	var role models.Role

	findRoleQuery := `
		SELECT *
		FROM roles
		WHERE name = ?
		LIMIT 1
	`

	err := database.DB.Raw(findRoleQuery, roleName).Scan(&role).Error
	if err != nil {
		return err
	}

	if role.ID == 0 {
		return errors.New("role not found")
	}

	insertQuery := `
		INSERT INTO user_roles (user_id, role_id)
		VALUES (?, ?)
		ON CONFLICT DO NOTHING
	`

	return database.DB.Exec(
		insertQuery,
		userID,
		role.ID,
	).Error
}

func (r *UserRepository) GetUsersDropdown() ([]map[string]any, error) {
	var users []models.User

	query := `
		SELECT *
		FROM users
		WHERE is_active = true
		ORDER BY name ASC
	`

	err := database.DB.Raw(query).Scan(&users).Error
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

func (r *UserRepository) loadRoles(user *models.User) {
	query := `
		SELECT r.*
		FROM roles r
		JOIN user_roles ur ON ur.role_id = r.id
		WHERE ur.user_id = ?
	`

	database.DB.Raw(query, user.ID).Scan(&user.Roles)
}

func (r *UserRepository) loadRolesAndPermissions(user *models.User) {
	r.loadRoles(user)

	for i := range user.Roles {
		query := `
			SELECT p.*
			FROM permissions p
			JOIN role_permissions rp ON rp.permission_id = p.id
			WHERE rp.role_id = ?
		`

		database.DB.Raw(query, user.Roles[i].ID).Scan(&user.Roles[i].Permissions)
	}
}

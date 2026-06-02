package dto

type CreateRoleDTO struct {
	Name string `json:"name" validate:"required"`
}

type CreatePermissionDTO struct {
	Name string `json:"name" validate:"required"`
}

type AssignPermissionDTO struct {
	RoleID       uint `json:"role_id" validate:"required"`
	PermissionID uint `json:"permission_id" validate:"required"`
}
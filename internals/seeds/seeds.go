package seeds

import (
	"log"

	"entry-system/internals/database"
	"entry-system/internals/models"
)

func RunSeeders() {
	seedRoles()
	seedPermissions()
	assignPermissionsToRoles()
}

func seedRoles() {
	roles := []string{
		"admin",
		"guard",
		"user",
		"manager",
	}

	for _, roleName := range roles {
		var role models.Role

		database.DB.FirstOrCreate(
			&role,
			models.Role{Name: roleName},
		)
	}

	log.Println("Roles seeded successfully")
}

func seedPermissions() {
	permissions := []string{
		"create_user",
		"view_users",
		"update_user",
		"delete_user",

		"view_own_profile",
		"update_own_profile",
		"view_own_history",

		"create_visitor",
		"view_visitors",
		"update_visitor",
		"restrict_visitor",

		"make_entry",
		"see_entry",

		"manage_roles",
		"manage_permissions",
	}

	for _, permissionName := range permissions {
		var permission models.Permission

		database.DB.FirstOrCreate(
			&permission,
			models.Permission{Name: permissionName},
		)
	}

	log.Println("Permissions seeded successfully")
}

func assignPermissionsToRoles() {
	adminPermissions := []string{
		"create_user",
		"view_users",
		"update_user",
		"delete_user",

		"view_own_profile",
		"update_own_profile",
		"view_own_history",

		"create_visitor",
		"view_visitors",
		"update_visitor",
		"restrict_visitor",

		"make_entry",
		"see_entry",

		"manage_roles",
		"manage_permissions",
	}

	guardPermissions := []string{
		"create_visitor",
		"view_visitors",
		"make_entry",
		"see_entry",
		"restrict_visitor",
	}

	userPermissions := []string{
		"view_own_profile",
		"update_own_profile",
		"view_own_history",
	}

	managerPermissions := []string{
		"view_users",
		"view_visitors",
		"see_entry",
	}

	assignRolePermissions("admin", adminPermissions)
	assignRolePermissions("guard", guardPermissions)
	assignRolePermissions("user", userPermissions)
	assignRolePermissions("manager", managerPermissions)

	log.Println("Role permissions assigned successfully")
}

func assignRolePermissions(
	roleName string,
	permissionNames []string,
) {
	var role models.Role

	err := database.DB.
		Where("name = ?", roleName).
		First(&role).Error

	if err != nil {
		log.Println("Role not found:", roleName)
		return
	}

	for _, permissionName := range permissionNames {
		var permission models.Permission

		err := database.DB.
			Where("name = ?", permissionName).
			First(&permission).Error

		if err != nil {
			log.Println("Permission not found:", permissionName)
			continue
		}

		database.DB.
			Model(&role).
			Association("Permissions").
			Append(&permission)
	}
}

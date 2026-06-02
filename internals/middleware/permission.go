package middleware

import (
	"entry-system/internals/database"
	"entry-system/internals/models"

	"github.com/gofiber/fiber/v3"
)

func RequirePermission(permission string) fiber.Handler {

	return func(c fiber.Ctx) error {

		userID := c.Locals("user_id")

		var user models.User

		err := database.DB.
			Preload("Roles.Permissions").
			First(&user, userID).Error

		if err != nil {
			return c.Status(401).JSON(
				fiber.Map{
					"message": "user not found",
				},
			)
		}

		for _, role := range user.Roles {

			for _, perm := range role.Permissions {

				if perm.Name == permission {
					return c.Next()
				}
			}
		}

		return c.Status(403).JSON(
			fiber.Map{
				"message": "permission denied",
			},
		)
	}
}
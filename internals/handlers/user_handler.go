package handlers

import (
	"strconv"

	"entry-system/internals/helpers"
	"entry-system/internals/models"
	"entry-system/internals/services"

	"github.com/gofiber/fiber/v3"
)

var userService = services.NewUserService()

func CreateUser(c fiber.Ctx) error {
	var body models.User

	if err := c.Bind().Body(&body); err != nil {
		return helpers.Error(c, 400, "invalid request body")
	}

	if err := userService.CreateUser(&body); err != nil {
		return helpers.Error(c, 400, err.Error())
	}
	data := fiber.Map{
		"name":  body.Name,
		"email": body.Email,
		"phone": body.Phone,
	}

	return helpers.Success(c, "user created successfully", data)
}

func GetUsers(c fiber.Ctx) error {
	users, err := userService.GetUsers()

	if err != nil {
		return helpers.Error(c, 500, err.Error())
	}

	return helpers.Success(c, "users fetched successfully", users)
}

func GetUser(c fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))

	user, err := userService.GetUser(uint(id))

	if err != nil {
		return helpers.Error(c, 404, "user not found")
	}

	return helpers.Success(c, "user fetched successfully", user)
}

func UpdateUser(c fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))

	var body struct {
		Name  string `json:"name"`
		Email string `json:"email"`
		Phone string `json:"phone"`
	}

	if err := c.Bind().Body(&body); err != nil {
		return helpers.Error(c, 400, "invalid request body")
	}

	err := userService.UpdateUser(
		uint(id),
		body.Name,
		body.Email,
		body.Phone,
	)

	if err != nil {
		return helpers.Error(c, 400, err.Error())
	}
	data := fiber.Map{
		"id":    id,
		"name":  body.Name,
		"email": body.Email,
		"phone": body.Phone,
	}

	return helpers.Success(c, "user updated successfully", data)
}

func DeleteUser(c fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))

	if err := userService.DeleteUser(uint(id)); err != nil {
		return helpers.Error(c, 400, err.Error())
	}
	data := fiber.Map{
		"id": id,
	}

	return helpers.Success(c, "user deleted successfully", data)
}

func DeactivateUser(c fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))

	if err := userService.DeactivateUser(uint(id)); err != nil {
		return helpers.Error(c, 400, err.Error())
	}

	data := fiber.Map{
		"id":     id,
		"status": "inactive",
	}
	return helpers.Success(c, "user deactivated successfully", data)
}

func VisitorHistory(c fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))

	history, err := userService.VisitorHistory(uint(id))

	if err != nil {
		return helpers.Error(c, 400, err.Error())
	}

	return helpers.Success(c, "visitor history fetched successfully", history)
}

func GetMe(c fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		return helpers.Error(c, 401, "invalid user id")
	}

	user, err := userService.GetUser(userID)
	if err != nil {
		return helpers.Error(c, 404, "user not found")
	}

	return helpers.Success(
		c,
		"profile fetched successfully",
		user,
	)
}

func UpdateMe(c fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		return helpers.Error(c, 401, "invalid user id")
	}

	var body struct {
		Name  string `json:"name"`
		Email string `json:"email"`
		Phone string `json:"phone"`
	}

	if err := c.Bind().Body(&body); err != nil {
		return helpers.Error(c, 400, "invalid request body")
	}

	err := userService.UpdateUser(
		userID,
		body.Name,
		body.Email,
		body.Phone,
	)

	if err != nil {
		return helpers.Error(c, 400, err.Error())
	}

	return helpers.Success(
		c,
		"profile updated successfully",
		nil,
	)
}

func GetUsersDropdown(c fiber.Ctx) error {
	users, err := userService.GetUsersDropdown()
	if err != nil {
		return helpers.Error(c, 400, err.Error())
	}

	return helpers.Success(
		c,
		"users dropdown fetched successfully",
		users,
	)
}

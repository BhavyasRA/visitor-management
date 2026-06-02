package handlers

import (
	"strconv"

	"entry-system/internals/helpers"
	"entry-system/internals/services"

	"github.com/gofiber/fiber/v3"
)

var guardService = services.NewGuardService()

func MakeEntry(c fiber.Ctx) error {
	var body struct {
		VisitorID uint `json:"visitor_id"`
	}

	if err := c.Bind().Body(&body); err != nil {
		return helpers.Error(c, 400, "invalid request body")
	}

	guardID, ok := c.Locals("user_id").(uint)
	if !ok {
		return helpers.Error(c, 401, "invalid user id in token")
	}

	err := guardService.MakeEntry(body.VisitorID, guardID)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}
	data := fiber.Map{
		"visitor_id": body.VisitorID,
		"guard_id":   guardID,
	}
	return helpers.Success(c, "entry created successfully", data)
}

func SeeEntries(c fiber.Ctx) error {
	filter := c.Query("filter", "all")
	from := c.Query("from")
	to := c.Query("to")

	entries, err := guardService.SeeEntries(filter, from, to)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return helpers.Success(c, "entries fetched successfully", entries)
}

func ExitVisitor(c fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("visitorId"))

	if err := guardService.ExitVisitor(uint(id)); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}
	data := fiber.Map{
		"id":     id,
		"status": "exited",
	}

	return helpers.Success(c, "visitor exited successfully", data)
}

func GuardRestrictVisitor(c fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("visitorId"))

	if err := guardService.RestrictVisitor(uint(id)); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}
	data := fiber.Map{
		"id":     id,
		"status": "restricted",
	}
	return helpers.Success(c, "visitor restricted successfully", data)
}

func GuardRestrictEmployee(c fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("userId"))

	if err := guardService.RestrictEmployee(uint(id)); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}
	data := fiber.Map{
		"id":     id,
		"status": "restricted",
	}
	return helpers.Success(c, "employee restricted successfully", data)
}

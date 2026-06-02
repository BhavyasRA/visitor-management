package handlers

import (
	"strconv"
	"time"

	"entry-system/internals/helpers"
	"entry-system/internals/services"

	"github.com/gofiber/fiber/v3"
)

var visitorService = services.NewVisitorService()

func CreateVisitor(c fiber.Ctx) error {
	var body struct {
		Name         string `json:"name"`
		Mobile       string `json:"mobile"`
		Email        string `json:"email"`
		Purpose      string `json:"purpose"`
		ToWhom       uint   `json:"to_whom"`
		VisitingTill string `json:"visiting_till"`
	}

	if err := c.Bind().Body(&body); err != nil {
		return helpers.Error(c, 400, "invalid request body")
	}

	var visitingTill *time.Time

	if body.VisitingTill != "" {
		parsed, err := time.Parse(time.RFC3339, body.VisitingTill)
		if err != nil {
			return helpers.Error(c, 400, "invalid visiting_till format, use RFC3339")
		}

		visitingTill = &parsed
	}

	err := visitorService.CreateVisitor(
		body.Name,
		body.Mobile,
		body.Email,
		body.Purpose,
		body.ToWhom,
		visitingTill,
	)

	if err != nil {
		return helpers.Error(c, 400, err.Error())
	}
	data := fiber.Map{
		"name":          body.Name,
		"mobile":        body.Mobile,
		"email":         body.Email,
		"purpose":       body.Purpose,
		"to_whom":       body.ToWhom,
		"visiting_till": body.VisitingTill,
	}

	return helpers.Success(c, "visitor created successfully", data)
}

func GetVisitors(c fiber.Ctx) error {
	visitors, err := visitorService.GetVisitors()

	if err != nil {
		return helpers.Error(c, 500, err.Error())
	}

	return helpers.Success(c, "visitors fetched successfully", visitors)
}

func UpdateVisitor(c fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))

	var body struct {
		Name    string `json:"name"`
		Mobile  string `json:"mobile"`
		Email   string `json:"email"`
		Purpose string `json:"purpose"`
	}

	if err := c.Bind().Body(&body); err != nil {
		return helpers.Error(c, 400, "invalid request body")
	}

	err := visitorService.UpdateVisitor(
		uint(id),
		body.Name,
		body.Mobile,
		body.Email,
		body.Purpose,
	)

	if err != nil {
		return helpers.Error(c, 400, err.Error())
	}

	data := fiber.Map{
		"name":    body.Name,
		"mobile":  body.Mobile,
		"email":   body.Email,
		"purpose": body.Purpose,
	}
	return helpers.Success(c, "visitor updated successfully", data)
}

func RestrictVisitor(c fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))

	if err := visitorService.RestrictVisitor(uint(id)); err != nil {
		return helpers.Error(c, 400, err.Error())
	}
	data := fiber.Map{
		"id":     id,
		"status": "restricted",
	}

	return helpers.Success(c, "visitor restricted successfully", data)
}

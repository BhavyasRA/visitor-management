package handlers

import (
	"fmt"
	"strconv"
	"time"

	"entry-system/internals/dto"
	"entry-system/internals/helpers"
	"entry-system/internals/services"

	"github.com/gofiber/fiber/v3"
)

var visitorService = services.NewVisitorService()

func CreateVisitor(c fiber.Ctx) error {
	var body struct {
		ImageURL     string `json:"image_url"`
		Name         string `json:"name"`
		Mobile       string `json:"mobile"`
		Email        string `json:"email"`
		Purpose      string `json:"purpose"`
		PersonToMeet uint   `json:"person_to_meet"`
		VisitingTill string `json:"visiting_till"`
	}

	if err := c.Bind().Body(&body); err != nil {
		return helpers.Error(c, 400, "invalid request body")
	}

	fmt.Printf("CREATE VISITOR PAYLOAD: %+v\n", body)
	fmt.Println("RAW BODY:", string(c.Body()))

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
		body.ImageURL,
		body.Purpose,
		body.PersonToMeet,
		visitingTill,
	)

	if err != nil {
		return helpers.Error(c, 400, err.Error())
	}

	data := fiber.Map{
		"image_url":      body.ImageURL,
		"name":           body.Name,
		"mobile":         body.Mobile,
		"email":          body.Email,
		"purpose":        body.Purpose,
		"person_to_meet": body.PersonToMeet,
		"visiting_till":  body.VisitingTill,
	}

	return helpers.Success(c, "visitor created successfully", data)
}

func GetVisitors(c fiber.Ctx) error {
	var filter dto.VisitorFilter
	filter.Name = c.Query("name")
	filter.Mobile = c.Query("mobile")
	filter.Email = c.Query("email")
	filter.From = c.Query("from")
	filter.To = c.Query("to")

	visitors, err := visitorService.GetVisitors(filter)
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

func GetVisitorEntriesGrouped(c fiber.Ctx) error {
	status := c.Query("status", "all")
	filter := c.Query("filter", "all")
	from := c.Query("from")
	to := c.Query("to")

	data, err := visitorService.GetGroupedVisitorEntries(
		status,
		filter,
		from,
		to,
	)

	if err != nil {
		return helpers.Error(c, 400, err.Error())
	}

	return helpers.Success(
		c,
		"visitor entries fetched successfully",
		data,
	)
}

func GetVisitorByMobile(c fiber.Ctx) error {
	mobile := c.Params("mobile")

	if mobile == "" {
		return helpers.Error(c, 400, "mobile number is required")
	}

	visitor, err := visitorService.GetVisitorByMobile(mobile)
	if err != nil {
		return helpers.Error(c, 404, "visitor not found")
	}

	data := fiber.Map{
		"id":         visitor.ID,
		"name":       visitor.Name,
		"mobile":     visitor.Mobile,
		"email":      visitor.Email,
		"image_url":  visitor.ImageURL,
		"restricted": visitor.IsRestricted,
	}

	return helpers.Success(
		c,
		"visitor fetched successfully",
		data,
	)
}
func GetVisitorStats(c fiber.Ctx) error {

	data, err := visitorService.GetVisitorStats()

	if err != nil {
		return helpers.Error(c, 400, err.Error())
	}

	return helpers.Success(
		c,
		"visitor stats fetched successfully",
		data,
	)
}

// func ExitVisitor(c fiber.Ctx) error {

// 	visitorID, err := strconv.Atoi(
// 		c.Params("visitorId"),
// 	)

// 	if err != nil {
// 		return helpers.Error(
// 			c,
// 			400,
// 			"invalid visitor id",
// 		)
// 	}

// 	if err := visitorService.ExitVisitor(
// 		uint(visitorID),
// 	); err != nil {
// 		return helpers.Error(
// 			c,
// 			400,
// 			err.Error(),
// 		)
// 	}

//		return helpers.Success(
//			c,
//			"visitor exited successfully",
//			fiber.Map{
//				"visitor_id": visitorID,
//				"status":     "exited",
//			},
//		)
//	}
func ExitVisitor(c fiber.Ctx) error {

	entryID, err := strconv.Atoi(c.Params("entryId"))
	if err != nil {
		return helpers.Error(c, 400, "invalid entry id")
	}

	if err := visitorService.ExitVisitor(uint(entryID)); err != nil {
		return helpers.Error(c, 400, err.Error())
	}

	return helpers.Success(
		c,
		"visitor exited successfully",
		fiber.Map{
			"entry_id": entryID,
			"status":   "exited",
		},
	)
}

func GetPersonsDropdown(c fiber.Ctx) error {

	data := visitorService.GetPersonsDropdown()

	return helpers.Success(
		c,
		"persons dropdown fetched successfully",
		data,
	)
}

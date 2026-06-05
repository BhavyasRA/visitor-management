package handlers

import (
	"strconv"
	"time"

	"entry-system/internals/dto"
	"entry-system/internals/helpers"
	"entry-system/internals/services"

	"github.com/gofiber/fiber/v3"
)

var visitorService = services.NewVisitorService()

func CreateVisitor(c fiber.Ctx) error {
	name := c.FormValue("name")
	mobile := c.FormValue("mobile")
	email := c.FormValue("email")
	purpose := c.FormValue("purpose")
	personToMeetStr := c.FormValue("person_to_meet")
	visitingTillStr := c.FormValue("visiting_till")

	if name == "" {
		return helpers.Error(c, 400, "name is required")
	}

	if mobile == "" {
		return helpers.Error(c, 400, "mobile is required")
	}

	if purpose == "" {
		return helpers.Error(c, 400, "purpose is required")
	}

	if personToMeetStr == "" {
		return helpers.Error(c, 400, "person_to_meet is required")
	}

	personToMeet, err := strconv.ParseUint(personToMeetStr, 10, 64)
	if err != nil {
		return helpers.Error(c, 400, "invalid person_to_meet")
	}

	var visitingTill *time.Time

	if visitingTillStr != "" {
		parsed, err := time.Parse(time.RFC3339, visitingTillStr)
		if err != nil {
			return helpers.Error(c, 400, "invalid visiting_till format, use RFC3339")
		}

		visitingTill = &parsed
	}

	photoFile, err := c.FormFile("photo")
	if err != nil {
		return helpers.Error(c, 400, "photo file is required")
	}

	photoURL, err := services.NewS3Service().UploadVisitorDocument(photoFile)
	if err != nil {
		return helpers.Error(c, 400, "failed to upload photo: "+err.Error())
	}

	documentFile, err := c.FormFile("identity_document")
	if err != nil {
		return helpers.Error(c, 400, "identity document file is required")
	}

	documentURL, err := services.NewS3Service().UploadVisitorDocument(documentFile)
	if err != nil {
		return helpers.Error(c, 400, "failed to upload identity document: "+err.Error())
	}

	entry, visitor, document, err := visitorService.CreateVisitor(
		name,
		mobile,
		email,
		photoURL,
		documentURL,
		purpose,
		uint(personToMeet),
		visitingTill,
	)

	if err != nil {
		return helpers.Error(c, 400, err.Error())
	}

	return helpers.Success(
		c,
		"visitor created successfully",
		fiber.Map{
			"visitor_id":            visitor.ID,
			"entry_id":              entry.ID,
			"document_id":           document.ID,
			"name":                  visitor.Name,
			"mobile":                visitor.Mobile,
			"email":                 visitor.Email,
			"photo_url":             document.PhotoURL,
			"identity_document_url": document.IdentityDocumentURL,
			"document_type":         document.DocumentType,
			"document_number":       document.DocumentNumber,
			"purpose":               entry.Purpose,
			"person_to_meet":        entry.PersonToMeet,
			"status":                entry.Status,
			"visiting_till":         entry.VisitingTill,
		},
	)
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
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return helpers.Error(c, 400, "invalid visitor id")
	}

	var body struct {
		Name   string `json:"name"`
		Mobile string `json:"mobile"`
		Email  string `json:"email"`
	}

	if err := c.Bind().Body(&body); err != nil {
		return helpers.Error(c, 400, "invalid request body")
	}

	err = visitorService.UpdateVisitor(
		uint(id),
		body.Name,
		body.Mobile,
		body.Email,
	)

	if err != nil {
		return helpers.Error(c, 400, err.Error())
	}

	return helpers.Success(
		c,
		"visitor updated successfully",
		fiber.Map{
			"id":     id,
			"name":   body.Name,
			"mobile": body.Mobile,
			"email":  body.Email,
		},
	)
}

func RestrictVisitor(c fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return helpers.Error(c, 400, "invalid visitor id")
	}

	if err := visitorService.RestrictVisitor(uint(id)); err != nil {
		return helpers.Error(c, 400, err.Error())
	}

	return helpers.Success(
		c,
		"visitor restricted successfully",
		fiber.Map{
			"id":     id,
			"status": "restricted",
		},
	)
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

	return helpers.Success(
		c,
		"visitor fetched successfully",
		fiber.Map{
			"id":            visitor.ID,
			"name":          visitor.Name,
			"mobile":        visitor.Mobile,
			"email":         visitor.Email,
			"is_restricted": visitor.IsRestricted,
		},
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

func ExitVisitor(c fiber.Ctx) error {
	entryIDParam := c.Params("entryId")

	if entryIDParam == "" {
		return helpers.Error(c, 400, "entry id is required")
	}

	entryID, err := strconv.Atoi(entryIDParam)
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

func GetVisitorDocumentForAI(c fiber.Ctx) error {
	documentID, err := strconv.Atoi(c.Params("documentId"))
	if err != nil {
		return helpers.Error(c, 400, "invalid document id")
	}

	document, err := visitorService.GetVisitorDocumentForAI(uint(documentID))
	if err != nil {
		return helpers.Error(c, 404, "visitor document not found")
	}

	return helpers.Success(
		c,
		"visitor document fetched successfully",
		fiber.Map{
			"visitor_document_id": document.ID,
			"visitor_id":          document.VisitorID,
			"document_url":        document.IdentityDocumentURL,
		},
	)
}

func UpdateDocumentAIResponse(c fiber.Ctx) error {

	var body struct {
		VisitorDocumentID uint   `json:"visitor_document_id"`
		DocumentType      string `json:"document_type"`
		DocumentNumber    string `json:"document_number"`
	}

	if err := c.Bind().Body(&body); err != nil {
		return helpers.Error(
			c,
			400,
			"invalid request body",
		)
	}

	err := visitorService.UpdateDocumentAIResponse(
		body.VisitorDocumentID,
		body.DocumentType,
		body.DocumentNumber,
	)

	if err != nil {
		return helpers.Error(
			c,
			400,
			err.Error(),
		)
	}

	return helpers.Success(
		c,
		"document updated successfully",
		nil,
	)
}

func GetActiveEntries(c fiber.Ctx) error {
	data, err := visitorService.GetActiveEntries()
	if err != nil {
		return helpers.Error(c, 400, err.Error())
	}

	return helpers.Success(c, "active entries fetched successfully", data)
}

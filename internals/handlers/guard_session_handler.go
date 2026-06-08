package handlers

import (
	"strconv"

	"entry-system/internals/helpers"
	"entry-system/internals/services"

	"github.com/gofiber/fiber/v3"
)

var guardSessionService = services.NewGuardSessionService()

func StartGuardSession(c fiber.Ctx) error {
	guardID, ok := c.Locals("user_id").(uint)
	if !ok {
		return helpers.Error(c, 401, "invalid guard id")
	}

	photoFile, err := c.FormFile("photo")
	if err != nil {
		return helpers.Error(c, 400, "login photo is required: "+err.Error())
	}

	photoURL, err := services.NewS3Service().UploadFile(photoFile, "guards")
	if err != nil {
		return helpers.Error(c, 400, "failed to upload login photo: "+err.Error())
	}

	session, err := guardSessionService.StartGuardSession(
		guardID,
		photoURL,
	)

	if err != nil {
		return helpers.Error(c, 400, err.Error())
	}

	return helpers.Success(
		c,
		"guard session started successfully",
		fiber.Map{
			"session_id":      session.ID,
			"guard_id":        session.GuardID,
			"login_photo_url": session.LoginPhotoURL,
			"login_at":        session.LoginAt,
			"status":          session.Status,
		},
	)
}

func EndGuardSession(c fiber.Ctx) error {
	guardID, ok := c.Locals("user_id").(uint)
	if !ok {
		return helpers.Error(c, 401, "invalid guard id")
	}

	photoFile, err := c.FormFile("photo")
	if err != nil {
		return helpers.Error(c, 400, "logout photo is required")
	}

	photoURL, err := services.NewS3Service().UploadFile(photoFile, "guards")
	if err != nil {
		return helpers.Error(c, 400, "failed to upload logout photo: "+err.Error())
	}

	err = guardSessionService.EndGuardSession(
		guardID,
		photoURL,
	)

	if err != nil {
		return helpers.Error(c, 400, err.Error())
	}

	return helpers.Success(
		c,
		"guard session ended successfully",
		fiber.Map{
			"guard_id":         guardID,
			"logout_photo_url": photoURL,
			"status":           "completed",
		},
	)
}

func GetAllGuardSessions(c fiber.Ctx) error {
	data, err := guardSessionService.GetAllGuardSessions()
	if err != nil {
		return helpers.Error(c, 400, err.Error())
	}

	return helpers.Success(
		c,
		"guard sessions fetched successfully",
		data,
	)
}

func GetGuardSessionsByGuardID(c fiber.Ctx) error {
	guardID, err := strconv.Atoi(c.Params("guardId"))
	if err != nil {
		return helpers.Error(c, 400, "invalid guard id")
	}

	data, err := guardSessionService.GetGuardSessionsByGuardID(uint(guardID))
	if err != nil {
		return helpers.Error(c, 400, err.Error())
	}

	return helpers.Success(
		c,
		"guard sessions fetched successfully",
		data,
	)
}

func CheckoutGuardSession(c fiber.Ctx) error {
	id := c.Params("id")

	logoutPhotoURL := c.FormValue("logout_photo_url")

	if logoutPhotoURL == "" {
		logoutPhotoURL = c.FormValue("logoutPhotoURL")
	}

	err := guardSessionService.CheckoutGuardSession(
		id,
		logoutPhotoURL,
	)

	if err != nil {
		return helpers.Error(c, 400, err.Error())
	}

	return helpers.Success(
		c,
		"guard session checked out successfully",
		nil,
	)
}

func GetGuardSessionInfo(c fiber.Ctx) error {

	userID, err := strconv.Atoi(c.Params("userID"))
	if err != nil {
		return helpers.Error(c, 400, "invalid user id")
	}

	data, err := guardSessionService.GetGuardSessionInfoByUserID(uint(userID))
	if err != nil {
		return helpers.Error(c, 404, "guard session not found")
	}

	return helpers.Success(
		c,
		"guard session info fetched successfully",
		data,
	)
}

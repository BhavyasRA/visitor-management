package handlers

import (
	"entry-system/internals/dto"
	"entry-system/internals/helpers"
	"entry-system/internals/services"

	//"entry-system/internals/utils"

	"github.com/gofiber/fiber/v3"
)

var authService = services.NewAuthService()

func Signup(c fiber.Ctx) error {
	var body dto.SignupDTO

	if err := c.Bind().Body(&body); err != nil {
		return helpers.Error(c, 400, "invalid request body")
	}

	// if msg := utils.ValidateStruct(body); msg != "" {
	// 	return helpers.Error(c, 400, msg)
	// }

	token, err := authService.Signup(
		body.Name,
		body.Email,
		body.Phone,
		body.Password,
	)
	data := fiber.Map{
		"name":               body.Name,
		"email":              body.Email,
		"phone":              body.Phone,
		"verification_token": token,
	}

	if err != nil {
		return helpers.Error(c, 400, err.Error())
	}

	return helpers.Success(
		c,
		"signup successful, verify your account",
		data,
	)
}

func Login(c fiber.Ctx) error {
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.Bind().Body(&body); err != nil {
		return helpers.Error(c, 400, "invalid request body")
	}

	token, err := authService.Login(body.Email, body.Password)

	if err != nil {
		return helpers.Error(c, 401, err.Error())
	}
	data := fiber.Map{
		"token": token,
		"role":  "guard",
		"email": body.Email,
	}

	return helpers.Success(
		c,
		"login successful",
		data,
	)

}
func VerifyAccount(c fiber.Ctx) error {
	var body struct {
		Token string `json:"token"`
	}

	if err := c.Bind().Body(&body); err != nil {
		return helpers.Error(c, 400, "invalid request body")
	}

	if err := authService.VerifyAccount(body.Token); err != nil {
		return helpers.Error(c, 400, err.Error())
	}

	return helpers.Success(
		c,
		"account verified successfully",
		nil,
	)
}

func ForgotPassword(c fiber.Ctx) error {
	var body struct {
		Email string `json:"email"`
	}

	if err := c.Bind().Body(&body); err != nil {
		return helpers.Error(c, 400, "invalid request body")
	}

	resetToken, err := authService.ForgotPassword(body.Email)

	if err != nil {
		return helpers.Error(c, 400, err.Error())
	}

	return helpers.Success(
		c,
		"reset token generated",
		fiber.Map{
			"reset_token": resetToken,
		},
	)
}

func ResetPassword(c fiber.Ctx) error {
	var body struct {
		Token       string `json:"token"`
		NewPassword string `json:"new_password"`
	}

	if err := c.Bind().Body(&body); err != nil {
		return helpers.Error(c, 400, "invalid request body")
	}

	if err := authService.ResetPassword(body.Token, body.NewPassword); err != nil {
		return helpers.Error(c, 400, err.Error())
	}

	return helpers.Success(c, "password reset successful", nil)
}

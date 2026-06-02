package middleware

import (
	"log"
	"os"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(c fiber.Ctx) error {

	log.Println("AUTH MIDDLEWARE HIT")

	authHeader := c.Get("Authorization")

	log.Println("AUTH HEADER:", authHeader)

	if authHeader == "" {
		return c.Status(401).JSON(fiber.Map{
			"success": false,
			"message": "missing authorization token",
		})
	}

	if !strings.HasPrefix(authHeader, "Bearer ") {
		return c.Status(401).JSON(fiber.Map{
			"success": false,
			"message": "invalid authorization format",
		})
	}

	tokenString := strings.TrimPrefix(
		authHeader,
		"Bearer ",
	)

	tokenString = strings.TrimSpace(tokenString)

	log.Println("TOKEN STRING:", tokenString)
	log.Println("JWT SECRET:", os.Getenv("JWT_SECRET"))

	claims := jwt.MapClaims{}

	token, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		},
	)

	if err != nil {

		log.Println("JWT PARSE ERROR:", err)

		return c.Status(401).JSON(fiber.Map{
			"success": false,
			"message": "invalid token",
			"error":   err.Error(),
		})
	}

	if !token.Valid {
		return c.Status(401).JSON(fiber.Map{
			"success": false,
			"message": "token not valid",
		})
	}

	userIDValue, ok := claims["user_id"]

	if !ok {
		return c.Status(401).JSON(fiber.Map{
			"success": false,
			"message": "user_id missing in token",
		})
	}

	userIDFloat, ok := userIDValue.(float64)

	if !ok {
		return c.Status(401).JSON(fiber.Map{
			"success": false,
			"message": "invalid user_id type",
		})
	}

	userID := uint(userIDFloat)

	log.Println("USER ID:", userID)

	c.Locals("user_id", userID)

	return c.Next()
}

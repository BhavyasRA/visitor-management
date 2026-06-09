package config

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/joho/godotenv"
)

func LoadEnv() {

	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func GetEnv(key string, fallback string) string {

	value := os.Getenv(key)

	if value == "" {
		return fallback
	}

	return value
}

func SetupCors(app *fiber.App) {

	app.Use(cors.New(cors.Config{

		Next:                  nil,
		AllowOriginsFunc:      nil,
		AllowOrigins:          []string{"*"},
		DisableValueRedaction: false,

		AllowMethods: []string{
			fiber.MethodGet,
			fiber.MethodHead,
			fiber.MethodPost,
			fiber.MethodPut,
			fiber.MethodPatch,
			fiber.MethodConnect,
			fiber.MethodOptions,
			fiber.MethodTrace,
		},

		AllowHeaders:        []string{},
		AllowCredentials:    false,
		ExposeHeaders:       []string{},
		MaxAge:              0,
		AllowPrivateNetwork: false,
	}))
}

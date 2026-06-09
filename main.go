package main

import (
	"log"
	"os"

	"entry-system/internals/config"
	"entry-system/internals/database"
	"entry-system/internals/routes"
	"entry-system/internals/seeds"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/static"
)

func main() {
	config.LoadEnv()

	database.ConnectDB()

	seedDb := config.GetEnv("SEED_DB", "true")

	if seedDb == "true" {
		seeds.RunSeeders()
	}

	_ = os.MkdirAll("./uploads/visitors", os.ModePerm)

	app := fiber.New(fiber.Config{
		BodyLimit:         50 * 1024 * 1024,
		StreamRequestBody: false,
	})
	config.SetupCors(app)

	app.Use("/uploads", static.New("./uploads"))

	routes.Setup(app)

	port := config.GetEnv("APP_PORT", "3000")

	log.Println("Application is running on port", port)

	log.Fatal(app.Listen(":" + port))
}

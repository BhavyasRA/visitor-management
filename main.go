package main

import (
	"log"

	"entry-system/internals/config"
	"entry-system/internals/database"
	"entry-system/internals/routes"
	"entry-system/internals/seeds"

	"github.com/gofiber/fiber/v3"
)

func main() {

	config.LoadEnv()

	database.ConnectDB()

	seeds.RunSeeders()

	app := fiber.New()

	routes.Setup(app)

	port := config.GetEnv(
		"APP_PORT",
		"3000",
	)

	log.Println("Application is running on port", port)
	log.Fatal(
		app.Listen(":" + port),
	)
}

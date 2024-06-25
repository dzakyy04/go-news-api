package main

import (
	"go-news-api/database"
	"go-news-api/routes"

	"github.com/gofiber/fiber/v2"
)

func main() {
	// Connect to database
	database.ConnectDatabase()

	// Migrate database
	database.MigrateDatabase()

	// Initialize fiber app
	app := fiber.New()

	// Initialize route
	routes.RouteInit(app)

	// Listen app on port 3000
	app.Listen(":3000")
}

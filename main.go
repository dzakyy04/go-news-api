package main

import (
	"go-news-api/routes"

	"github.com/gofiber/fiber/v2"
)

func main() {
	// Initialize fiber app
	app := fiber.New()

	// Initialize route
	routes.RouteInit(app)

	// Listen app on port 3000
	app.Listen(":3000")
}

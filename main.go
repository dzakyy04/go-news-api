package main

import (
	"go-news-api/cmd"
	"go-news-api/database"
	"go-news-api/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"

	_ "go-news-api/docs"
)

// @title News API
// @version 1.0
// @description This is a news API for GDSC final project
// @termsOfService http://swagger.io/terms/
// @contact.name Dewa Sheva Dzaky
// @contact.email dzakylinggau@gmail.com
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:3000
// @BasePath /api
func main() {
	// Connect to database
	database.ConnectDatabase()

	// Migrate database
	database.MigrateDatabase()

	// Initialize fiber app
	app := fiber.New()

	// Cobra for cli
	cmd.Execute()

	// Swagger for api docs
	app.Get("/swagger/*", swagger.HandlerDefault)

	// Initialize route
	routes.RouteInit(app)

	// Listen app on port 3000
	app.Listen(":3000")
}

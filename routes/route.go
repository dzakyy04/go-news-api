package routes

import (
	"go-news-api/controllers"

	"github.com/gofiber/fiber/v2"
)

func RouteInit(route *fiber.App) {
	route.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.JSON(fiber.Map{
			"message": "Hello World",
		})
	})

	// Category routes
	route.Get("/categories", controllers.GetAllCategories)
	route.Get("/categories/:id", controllers.GetCategoryById)
	route.Post("/categories", controllers.CreateCategory)
}

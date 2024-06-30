package routes

import (
	"go-news-api/controllers"
	"go-news-api/middleware"

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
	route.Put("/categories/:id", controllers.UpdateCategory)
	route.Delete("/categories/:id", controllers.DeleteCategory)

	// Auth routes
	route.Post("/register", controllers.Register)
	route.Post("/login", controllers.Login)
	route.Post("/verification-email", controllers.SendVerificationEmail)
	route.Post("/verify-email", controllers.VerifyEmail)
	route.Get("/profile", middleware.AuthMiddleware, controllers.GetProfile)
}

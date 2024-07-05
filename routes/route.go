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

	// Static asset
	route.Static("/public", "./public")

	// Prefix /api
	api := route.Group("/api")

	// Category routes
	api.Get("/categories", controllers.GetAllCategories)
	api.Get("/categories/:id", controllers.GetCategoryById)
	api.Post("/categories", controllers.CreateCategory)
	api.Put("/categories/:id", controllers.UpdateCategory)
	api.Delete("/categories/:id", controllers.DeleteCategory)

	// Auth routes
	api.Post("/register", controllers.Register)
	api.Post("/login", controllers.Login)
	api.Post("/email-verification/request", controllers.SendVerificationEmail)
	api.Post("/email-verification/verify", controllers.VerifyEmail)
	api.Get("/profile", middleware.AuthMiddleware, controllers.GetProfile)
	api.Post("/reset-password/request", controllers.SendResetPasswordEmail)
	api.Post("/reset-password/verify", controllers.VerifyOtpReset)
	api.Post("/reset-password", controllers.ResetPassword)

	// Article routes
	api.Get("/articles", controllers.GetAllArticles)
	api.Get("/articles/:id", controllers.GetArticleById)
	api.Post("/articles", controllers.CreateArticle)
	api.Put("/articles/:slug", controllers.UpdateArticle)
}

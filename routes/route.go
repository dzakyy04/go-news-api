package routes

import "github.com/gofiber/fiber/v2"

func RouteInit(route *fiber.App) {
	route.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.JSON(fiber.Map{
			"message": "Hello World",
		})
	})
}

package utils

import "github.com/gofiber/fiber/v2"

func SendErrorResponse(ctx *fiber.Ctx, status int, message string, err error) error {
	return ctx.Status(status).JSON(fiber.Map{
		"success": false,
		"message": message,
		"error":   err.Error(),
	})
}

func SendSuccessResponse(ctx *fiber.Ctx, status int, message string, data interface{}) error {
	return ctx.Status(status).JSON(fiber.Map{
		"success": true,
		"message": message,
		"data":    data,
	})
}

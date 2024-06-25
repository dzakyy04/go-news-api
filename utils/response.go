package utils

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func SendErrorResponse(ctx *fiber.Ctx, status int, message string, err error) error {
	return ctx.Status(status).JSON(fiber.Map{
		"success": false,
		"message": message,
		"error":   err.Error(),
	})
}

func SendValidationErrorResponse(ctx *fiber.Ctx, err error) error {
	var errors []string
	for _, err := range err.(validator.ValidationErrors) {
		errors = append(errors, FormatValidationError(err))
	}
	return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
		"success": false,
		"message": "Validation failed",
		"errors":  errors,
	})
}

func SendSuccessResponse(ctx *fiber.Ctx, status int, message string, data interface{}) error {
	return ctx.Status(status).JSON(fiber.Map{
		"success": true,
		"message": message,
		"data":    data,
	})
}

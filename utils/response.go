package utils

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func SendErrorResponse(ctx *fiber.Ctx, status int, message string, err error) error {
	response := fiber.Map{
		"success": false,
		"message": message,
	}

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		var errors []string
		for _, err := range validationErrors {
			errors = append(errors, FormatValidationError(err))
		}
		response["errors"] = errors
	} else if err != nil {
		response["errors"] = []string{err.Error()}
	}

	return ctx.Status(status).JSON(response)
}

func SendSuccessResponse(ctx *fiber.Ctx, status int, message string) error {
	return ctx.Status(status).JSON(fiber.Map{
		"success": true,
		"message": message,
	})
}

func SendSuccessResponseWithData(ctx *fiber.Ctx, status int, message string, data interface{}) error {
	return ctx.Status(status).JSON(fiber.Map{
		"success": true,
		"message": message,
		"data":    data,
	})
}

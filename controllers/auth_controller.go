package controllers

import (
	"go-news-api/database"
	"go-news-api/models/entity"
	"go-news-api/models/request"
	"go-news-api/utils"

	"github.com/gofiber/fiber/v2"
)

func Register(ctx *fiber.Ctx) error {
	request := new(request.RegisterRequest)

	// Parse request body
	if err := ctx.BodyParser(request); err != nil {
		return utils.SendErrorResponse(ctx, fiber.StatusBadRequest, "Invalid request body", err)
	}

	// Validate request
	if err := validate.Struct(request); err != nil {
		return utils.SendValidationErrorResponse(ctx, err)
	}

	// Create user
	user := entity.User{
		Name:  request.Name,
		Email: request.Email,
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(request.Password)
	if err != nil {
		return utils.SendErrorResponse(ctx, fiber.StatusInternalServerError, "Failed to hash password", err)
	}
	user.Password = hashedPassword

	// Save user
	if err := database.DB.Create(&user).Error; err != nil {
		return utils.SendErrorResponse(ctx, fiber.StatusInternalServerError, "Failed to create user", err)
	}

	return utils.SendSuccessResponse(ctx, fiber.StatusCreated, "User created successfully", fiber.Map{
		"user": user,
	})
}

package controllers

import (
	"go-news-api/database"
	"go-news-api/models/entity"
	"go-news-api/models/request"
	"go-news-api/utils"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func Register(ctx *fiber.Ctx) error {
	request := new(request.RegisterRequest)

	// Parse request body
	if err := ctx.BodyParser(request); err != nil {
		return utils.SendErrorResponse(ctx, fiber.StatusBadRequest, "Invalid request body", err)
	}

	// Validate request
	if err := utils.Validate.Struct(request); err != nil {
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

	return utils.SendSuccessResponse(ctx, fiber.StatusCreated, "Sucessfully registered", fiber.Map{
		"user": user,
	})
}

func Login(ctx *fiber.Ctx) error {
	request := new(request.LoginRequest)

	// Parse request body
	if err := ctx.BodyParser(request); err != nil {
		return utils.SendErrorResponse(ctx, fiber.StatusBadRequest, "Invalid request body", err)
	}

	// Validate request
	if err := utils.Validate.Struct(request); err != nil {
		return utils.SendValidationErrorResponse(ctx, err)
	}

	// Find user
	var user entity.User
	if err := database.DB.Where("email = ?", request.Email).First(&user).Error; err != nil {
		return utils.SendErrorResponse(ctx, fiber.StatusUnauthorized, "Invalid email or password", err)
	}

	// Check password
	if err := utils.VerifyPassword(request.Password, user.Password); err != nil {
		return utils.SendErrorResponse(ctx, fiber.StatusUnauthorized, "Invalid email or password", err)
	}

	// Generate JWT token
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"iat":     time.Now().Unix(),
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	}

	token, err := utils.GenerateToken(&claims)
	if err != nil {
		return utils.SendErrorResponse(ctx, fiber.StatusNotFound, "Invalid email or password", err)
	}

	return utils.SendSuccessResponse(ctx, fiber.StatusOK, "Sucessfully logged in", fiber.Map{
		"token": token,
		"user":  user,
	})
}

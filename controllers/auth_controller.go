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
		return utils.SendErrorResponse(ctx, fiber.StatusBadRequest, "Failed to register", err)
	}

	// Validate request
	if err := utils.Validate.Struct(request); err != nil {
		return utils.SendErrorResponse(ctx, fiber.StatusBadRequest, "Failed to register", err)
	}

	// Create user
	user := entity.User{
		Name:  request.Name,
		Email: request.Email,
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(request.Password)
	if err != nil {
		return utils.SendErrorResponse(ctx, fiber.StatusInternalServerError, "Failed to register", err)
	}
	user.Password = hashedPassword

	// Save user
	if err := database.DB.Create(&user).Error; err != nil {
		return utils.SendErrorResponse(ctx, fiber.StatusInternalServerError, "Failed to register", err)
	}

	return utils.SendSuccessResponse(ctx, fiber.StatusCreated, "Sucessfully registered", fiber.Map{
		"user": user,
	})
}

func Login(ctx *fiber.Ctx) error {
	request := new(request.LoginRequest)

	// Parse request body
	if err := ctx.BodyParser(request); err != nil {
		return utils.SendErrorResponse(ctx, fiber.StatusBadRequest, "Failed to login", err)
	}

	// Validate request
	if err := utils.Validate.Struct(request); err != nil {
		return utils.SendErrorResponse(ctx, fiber.StatusBadRequest, "Failed to login", err)
	}

	// Find user
	var user entity.User
	if err := database.DB.Where("email = ?", request.Email).First(&user).Error; err != nil {
		return utils.SendErrorResponse(ctx, fiber.StatusUnauthorized, "Failed to login", err)
	}

	// Check password
	if err := utils.VerifyPassword(request.Password, user.Password); err != nil {
		return utils.SendErrorResponse(ctx, fiber.StatusUnauthorized, "Failed to login", err)
	}

	// Generate JWT token
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"iat":     time.Now().Unix(),
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	}

	token, err := utils.GenerateToken(&claims)
	if err != nil {
		return utils.SendErrorResponse(ctx, fiber.StatusNotFound, "Failed to login", err)
	}

	return utils.SendSuccessResponse(ctx, fiber.StatusOK, "Successfully logged in", fiber.Map{
		"token": token,
		"user":  user,
	})
}

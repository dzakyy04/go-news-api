package middleware

import (
	"errors"
	"go-news-api/database"
	"go-news-api/models/entity"
	"go-news-api/utils"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(ctx *fiber.Ctx) error {
	// Check token from Authorization header
	authHeader := ctx.Get("Authorization")
	if authHeader == "" {
		return utils.SendErrorResponse(ctx, fiber.StatusUnauthorized, "Unauthorized", errors.New("authorization header is missing"))
	}

	// Split header to get token
	tokenParts := strings.Split(authHeader, "Bearer ")
	if len(tokenParts) != 2 {
		return utils.SendErrorResponse(ctx, fiber.StatusUnauthorized, "Unauthorized", errors.New("invalid authorization header format"))
	}

	tokenString := tokenParts[1]
	if tokenString == "" {
		return utils.SendErrorResponse(ctx, fiber.StatusUnauthorized, "Unauthorized", errors.New("token is empty"))
	}

	// Parse and validate token
	token, err := utils.ParseToken(tokenString)
	if err != nil {
		return utils.SendErrorResponse(ctx, fiber.StatusUnauthorized, "Unauthorized", errors.New("invalid or expired token"))
	}

	// Get claims from token
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return utils.SendErrorResponse(ctx, fiber.StatusUnauthorized, "Unauthorized", errors.New("invalid or expired token"))
	}

	// Get user_id from claims
	userID, ok := claims["user_id"].(float64)
	if !ok {
		return utils.SendErrorResponse(ctx, fiber.StatusUnauthorized, "Unauthorized", errors.New("invalid or expired token"))
	}

	// Find user
	var user entity.User
	if err := database.DB.First(&user, "id = ?", userID).Error; err != nil {
		return utils.SendErrorResponse(ctx, fiber.StatusUnauthorized, "Unauthorized", err)
	}

	// Attach user to context
	ctx.Locals("user", &user)

	return ctx.Next()
}

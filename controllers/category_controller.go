package controllers

import (
	"go-news-api/database"
	"go-news-api/models/entity"

	"github.com/gofiber/fiber/v2"
)

func GetAllCategories(ctx *fiber.Ctx) error {
	var categories []entity.Category

	// Fetch all categories
	if err := database.DB.Find(&categories).Error; err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to fetch categories",
			"error":   err.Error(),
		})
	}

	return ctx.JSON(fiber.Map{
		"success": true,
		"message": "Categories fetched successfully",
		"data": fiber.Map{
			"categories":      categories,
			"totalCategories": len(categories),
		},
	})
}

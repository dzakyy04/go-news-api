package controllers

import (
	"go-news-api/database"
	"go-news-api/models/entity"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
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

func GetCategoryById(ctx *fiber.Ctx) error {
	// Get category ID from URL parameter
	categoryId := ctx.Params("id")

	// Find category by ID
	var category entity.Category
	err := database.DB.First(&category, "id = ?", categoryId).Error
	if err != nil {
		// If category not found
		if err == gorm.ErrRecordNotFound {
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"success": false,
				"message": "Category not found",
				"error":   err.Error(),
			})
		}

		// If error occurred
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to fetch category",
			"error":   err.Error(),
		})
	}

	return ctx.JSON(fiber.Map{
		"success": true,
		"message": "Category fetched successfully",
		"data": fiber.Map{
			"category": category,
		},
	})
}

package controllers

import (
	"go-news-api/database"
	"go-news-api/models/entity"
	"go-news-api/utils"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func GetAllCategories(ctx *fiber.Ctx) error {
	var categories []entity.Category

	// Fetch all categories
	if err := database.DB.Find(&categories).Error; err != nil {
		return utils.SendErrorResponse(ctx, fiber.StatusInternalServerError, "Failed to fetch categories", err)
	}

	return utils.SendSuccessResponse(ctx, fiber.StatusOK, "Categories fetched successfully", fiber.Map{
		"categories":      categories,
		"totalCategories": len(categories),
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
			return utils.SendErrorResponse(ctx, fiber.StatusNotFound, "Category not found", err)
		}

		// If error occurred
		return utils.SendErrorResponse(ctx, fiber.StatusInternalServerError, "Failed to fetch category", err)
	}

	return utils.SendSuccessResponse(ctx, fiber.StatusOK, "Category fetched successfully", fiber.Map{
		"category": category,
	})
}

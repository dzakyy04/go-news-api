package controllers

import (
	"go-news-api/database"
	"go-news-api/models/entity"
	"go-news-api/models/request"
	"go-news-api/utils"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

var validate = validator.New()

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
	if err := database.DB.First(&category, "id = ?", categoryId).Error; err != nil {
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

func CreateCategory(ctx *fiber.Ctx) error {
	request := new(request.CreateCategoryRequest)

	// Parse request body
	if err := ctx.BodyParser(request); err != nil {
		return utils.SendErrorResponse(ctx, fiber.StatusBadRequest, "Invalid request body", err)
	}

	// Validate request
	if err := validate.Struct(request); err != nil {
		return utils.SendValidationErrorResponse(ctx, err)
	}

	// Create category
	category := entity.Category{
		Name:        request.Name,
		Description: request.Description,
	}

	if err := database.DB.Create(&category).Error; err != nil {
		return utils.SendErrorResponse(ctx, fiber.StatusInternalServerError, "Failed to create category", err)
	}

	return utils.SendSuccessResponse(ctx, fiber.StatusCreated, "Category created successfully", fiber.Map{
		"category": category,
	})
}

func UpdateCategory(ctx *fiber.Ctx) error {
	categoryId := ctx.Params("id")

	// Check if category exists
	var category entity.Category
	if err := database.DB.First(&category, "id = ?", categoryId).Error; err != nil {
		// If category not found
		if err == gorm.ErrRecordNotFound {
			return utils.SendErrorResponse(ctx, fiber.StatusNotFound, "Category not found", err)
		}
		// If error occurred
		return utils.SendErrorResponse(ctx, fiber.StatusInternalServerError, "Failed to fetch category", err)
	}

	// Parse request body
	request := new(request.UpdateCategoryRequest)
	if err := ctx.BodyParser(request); err != nil {
		return utils.SendErrorResponse(ctx, fiber.StatusBadRequest, "Invalid request body", err)
	}

	// Validate request
	if err := validate.Struct(request); err != nil {
		return utils.SendValidationErrorResponse(ctx, err)
	}

	// Update category
	category.Name = request.Name
	category.Description = request.Description

	if err := database.DB.Save(&category).Error; err != nil {
		return utils.SendErrorResponse(ctx, fiber.StatusInternalServerError, "Failed to update category", err)
	}

	return utils.SendSuccessResponse(ctx, fiber.StatusOK, "Category updated successfully", fiber.Map{
		"category": category,
	})
}

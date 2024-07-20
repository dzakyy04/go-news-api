package controllers

import (
	"go-news-api/database"
	"go-news-api/models/entity"
	"go-news-api/models/request"
	"go-news-api/utils"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// GetAllCategories godoc
// @Summary Get all categories
// @Description Fetches all categories from the database.
// @Tags Categories
// @Accept  json
// @Produce  json
// @Router /categories [get]
func GetAllCategories(ctx *fiber.Ctx) error {
	var categories []entity.Category

	// Fetch all categories
	if err := database.DB.Find(&categories).Error; err != nil {
		return utils.SendErrorResponse(ctx, fiber.StatusInternalServerError, "Failed to fetch categories", err)
	}

	return utils.SendSuccessResponseWithData(ctx, fiber.StatusOK, "Successfully fetched categories", fiber.Map{
		"categories":      categories,
		"total_categories": len(categories),
	})
}

// GetCategoryById godoc
// @Summary Get category by ID
// @Description Fetches a category by its ID from the database.
// @Tags Categories
// @Accept  json
// @Produce  json
// @Param id path int true "Category ID"
// @Router /categories/{id} [get]
func GetCategoryById(ctx *fiber.Ctx) error {
	// Get category ID from URL parameter
	categoryId := ctx.Params("id")

	// Find category by ID
	var category entity.Category
	if err := database.DB.First(&category, "id = ?", categoryId).Error; err != nil {
		// If category not found
		if err == gorm.ErrRecordNotFound {
			return utils.SendErrorResponse(ctx, fiber.StatusNotFound, "Failed to fetch category", err)
		}
		// If error occurred
		return utils.SendErrorResponse(ctx, fiber.StatusInternalServerError, "Failed to fetch category", err)
	}

	return utils.SendSuccessResponseWithData(ctx, fiber.StatusOK, "Succesfully fetched category", fiber.Map{
		"category": category,
	})
}

// CreateCategory godoc
// @Summary Create category
// @Description Creates a new category in the database.
// @Tags Categories
// @Accept  multipart/form-data
// @Produce  json
// @Param name formData string true "Category Name"
// @Param description formData string true "Category Description"
// @Router /categories [post]
func CreateCategory(ctx *fiber.Ctx) error {
	request := new(request.CategoryRequest)

	// Parse request body
	if err := ctx.BodyParser(request); err != nil {
		return utils.SendErrorResponse(ctx, fiber.StatusBadRequest, "Failed to create category", err)
	}

	// Validate request
	if err := utils.Validate.Struct(request); err != nil {
		return utils.SendErrorResponse(ctx, fiber.StatusBadRequest, "Failed to create category", err)
	}

	// Create category
	category := entity.Category{
		Name:        request.Name,
		Description: request.Description,
	}

	if err := database.DB.Create(&category).Error; err != nil {
		return utils.SendErrorResponse(ctx, fiber.StatusInternalServerError, "Failed to create category", err)
	}

	return utils.SendSuccessResponse(ctx, fiber.StatusCreated, "Successfully created category")
}

// UpdateCategory godoc
// @Summary Update category
// @Description Updates a category in the database.
// @Tags Categories
// @Accept  multipart/form-data
// @Produce  json
// @Param id path int true "Category ID"
// @Param name formData string true "Category Name"
// @Param description formData string true "Category Description"
// @Router /categories/{id} [put]
func UpdateCategory(ctx *fiber.Ctx) error {
	categoryId := ctx.Params("id")

	// Check if category exists
	var category entity.Category
	if err := database.DB.First(&category, "id = ?", categoryId).Error; err != nil {
		// If category not found
		if err == gorm.ErrRecordNotFound {
			return utils.SendErrorResponse(ctx, fiber.StatusNotFound, "Failed to update category", err)
		}
		// If error occurred
		return utils.SendErrorResponse(ctx, fiber.StatusInternalServerError, "Failed to update category", err)
	}

	// Parse request body
	request := new(request.CategoryRequest)
	if err := ctx.BodyParser(request); err != nil {
		return utils.SendErrorResponse(ctx, fiber.StatusBadRequest, "Failed to update category", err)
	}

	// Validate request
	if err := utils.Validate.Struct(request); err != nil {
		return utils.SendErrorResponse(ctx, fiber.StatusBadRequest, "Failed to update category", err)
	}

	// Update category
	category.Name = request.Name
	category.Description = request.Description

	if err := database.DB.Save(&category).Error; err != nil {
		return utils.SendErrorResponse(ctx, fiber.StatusInternalServerError, "Failed to update category", err)
	}

	return utils.SendSuccessResponse(ctx, fiber.StatusOK, "Successfully updated category")
}

// DeleteCategory godoc
// @Summary Delete category
// @Description Deletes a category from the database.
// @Tags Categories
// @Accept  json
// @Produce  json
// @Param id path int true "Category ID"
// @Router /categories/{id} [delete]
func DeleteCategory(ctx *fiber.Ctx) error {
	categoryId := ctx.Params("id")

	// Check if category exists
	var category entity.Category
	if err := database.DB.First(&category, "id = ?", categoryId).Error; err != nil {
		// If category not found
		if err == gorm.ErrRecordNotFound {
			return utils.SendErrorResponse(ctx, fiber.StatusNotFound, "Failed to delete category", err)
		}
		// If error occurred
		return utils.SendErrorResponse(ctx, fiber.StatusInternalServerError, "Failed to delete category", err)
	}

	if err := database.DB.Delete(&category).Error; err != nil {
		return utils.SendErrorResponse(ctx, fiber.StatusInternalServerError, "Failed to delete category", err)
	}

	return utils.SendSuccessResponse(ctx, fiber.StatusOK, "Successfully delete category")
}

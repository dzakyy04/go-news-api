package controllers

import (
	"go-news-api/database"
	"go-news-api/models/entity"
	"go-news-api/models/request"
	"go-news-api/utils"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// GetAllTags godoc
// @Summary Get all tags
// @Description Fetches all tags from the database.
// @Tags Tags
// @Produce  json
// @Router /tags [get]
func GetAllTags(ctx *fiber.Ctx) error {
	var tags []entity.Tag

	// Fetch all tags
	if err := database.DB.Find(&tags).Error; err != nil {
		return utils.SendErrorResponse(ctx, fiber.StatusInternalServerError, "Failed to fetch tags", err)
	}

	return utils.SendSuccessResponseWithData(ctx, fiber.StatusOK, "Successfully fetched tags", fiber.Map{
		"tags":       tags,
		"total_tags": len(tags),
	})
}

// GetTagById godoc
// @Summary Get a tag by ID
// @Description Fetches a tag by its ID from the database.
// @Tags Tags
// @Produce  json
// @Param id path int true "Tag ID"
// @Router /tags/{id} [get]
func GetTagById(ctx *fiber.Ctx) error {
	// Get tag ID from URL parameter
	tagId := ctx.Params("id")

	// Find tag by ID
	var tag entity.Tag
	if err := database.DB.First(&tag, "id = ?", tagId).Error; err != nil {
		// If tag not found
		if err == gorm.ErrRecordNotFound {
			return utils.SendErrorResponse(ctx, fiber.StatusNotFound, "Failed to fetch tag", err)
		}
		// If error occurred
		return utils.SendErrorResponse(ctx, fiber.StatusInternalServerError, "Failed to fetch tag", err)
	}

	return utils.SendSuccessResponseWithData(ctx, fiber.StatusOK, "Succesfully fetched tag", fiber.Map{
		"tag": tag,
	})
}

// CreateTag godoc
// @Summary Create a new tag
// @Description Creates a new tag and saves it to the database.
// @Tags Tags
// @Accept  multipart/form-data
// @Produce  json
// @Param name formData string true "Tag Name"
// @Router /tags [post]
func CreateTag(ctx *fiber.Ctx) error {
	request := new(request.TagRequest)

	// Parse request body
	if err := ctx.BodyParser(request); err != nil {
		return utils.SendErrorResponse(ctx, fiber.StatusBadRequest, "Failed to create tag", err)
	}

	// Validate request
	if err := utils.Validate.Struct(request); err != nil {
		return utils.SendErrorResponse(ctx, fiber.StatusBadRequest, "Failed to create tag", err)
	}

	// Create tag
	tag := entity.Tag{
		Name: request.Name,
	}

	if err := database.DB.Create(&tag).Error; err != nil {
		return utils.SendErrorResponse(ctx, fiber.StatusInternalServerError, "Failed to create tag", err)
	}

	return utils.SendSuccessResponse(ctx, fiber.StatusCreated, "Successfully created tag")
}

// UpdateTag godoc
// @Summary Update an existing tag
// @Description Updates an existing tag identified by ID and saves changes to the database.
// @Tags Tags
// @Accept  multipart/form-data
// @Produce  json
// @Param id path int true "Tag ID"
// @Param name formData string true "Tag Name"
// @Router /tags/{id} [put]
func UpdateTag(ctx *fiber.Ctx) error {
	tagId := ctx.Params("id")

	// Check if tag exist
	var tag entity.Tag
	if err := database.DB.First(&tag, "id = ?", tagId).Error; err != nil {
		// If tag not found
		if err == gorm.ErrRecordNotFound {
			return utils.SendErrorResponse(ctx, fiber.StatusNotFound, "Failed to update tag", err)
		}
		// If error occurred
		return utils.SendErrorResponse(ctx, fiber.StatusInternalServerError, "Failed to update tag", err)
	}

	// Parse request body
	request := new(request.TagRequest)
	if err := ctx.BodyParser(request); err != nil {
		return utils.SendErrorResponse(ctx, fiber.StatusBadRequest, "Failed to update tag", err)
	}

	// Validate request
	if err := utils.Validate.Struct(request); err != nil {
		return utils.SendErrorResponse(ctx, fiber.StatusBadRequest, "Failed to update tag", err)
	}

	// Update tag
	if err := database.DB.Model(&tag).Updates(request).Error; err != nil {
		return utils.SendErrorResponse(ctx, fiber.StatusInternalServerError, "Failed to update tag", err)
	}

	return utils.SendSuccessResponse(ctx, fiber.StatusOK, "Successfully updated tag")
}

// DeleteTag godoc
// @Summary Delete a tag
// @Description Deletes a tag identified by ID from the database.
// @Tags Tags
// @Produce  json
// @Param id path int true "Tag ID"
// @Router /tags/{id} [delete]
func DeleteTag(ctx *fiber.Ctx) error {
	tagId := ctx.Params("id")

	// Check if tag exist
	var tag entity.Tag
	if err := database.DB.First(&tag, "id = ?", tagId).Error; err != nil {
		// If tag not found
		if err == gorm.ErrRecordNotFound {
			return utils.SendErrorResponse(ctx, fiber.StatusNotFound, "Failed to delete tag", err)
		}
		// If error occurred
		return utils.SendErrorResponse(ctx, fiber.StatusInternalServerError, "Failed to delete tag", err)
	}

	// Delete tag
	if err := database.DB.Delete(&tag).Error; err != nil {
		return utils.SendErrorResponse(ctx, fiber.StatusInternalServerError, "Failed to delete tag", err)
	}

	return utils.SendSuccessResponse(ctx, fiber.StatusOK, "Successfully deleted tag")
}

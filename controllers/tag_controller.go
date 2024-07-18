package controllers

import (
	"go-news-api/database"
	"go-news-api/models/entity"
	"go-news-api/utils"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func GetAllTags(ctx *fiber.Ctx) error {
	var tags []entity.Tag

	// Fetch all tags
	if err := database.DB.Find(&tags).Error; err != nil {
		return utils.SendErrorResponse(ctx, fiber.StatusInternalServerError, "Failed to fetch tags", err)
	}

	return utils.SendSuccessResponse(ctx, fiber.StatusOK, "Successfully fetched tags", fiber.Map{
		"tags":      tags,
		"totalTags": len(tags),
	})
}

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

	return utils.SendSuccessResponse(ctx, fiber.StatusOK, "Succesfully fetched tag", fiber.Map{
		"tag": tag,
	})
}

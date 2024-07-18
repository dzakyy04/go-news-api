package controllers

import (
	"go-news-api/database"
	"go-news-api/models/entity"
	"go-news-api/utils"

	"github.com/gofiber/fiber/v2"
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

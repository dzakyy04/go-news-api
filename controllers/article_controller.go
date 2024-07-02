package controllers

import (
	"go-news-api/database"
	"go-news-api/models/entity"
	"go-news-api/utils"

	"github.com/gofiber/fiber/v2"
)

func GetAllArticles(ctx *fiber.Ctx) error {
	var articles []entity.Article

	// Fetch all articles
	if err := database.DB.Preload("Category").Preload("Author").Find(&articles).Error; err != nil {
		return utils.SendErrorResponse(ctx, fiber.StatusInternalServerError, "Failed to fetch articles", err)
	}

	return utils.SendSuccessResponse(ctx, fiber.StatusOK, "Successfully fetched articles", fiber.Map{
		"articles":      articles,
		"totalArticles": len(articles),
	})
}

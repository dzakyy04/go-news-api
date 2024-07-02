package controllers

import (
	"go-news-api/database"
	"go-news-api/models/entity"
	"go-news-api/utils"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
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

func GetArticleById(ctx *fiber.Ctx) error {
	// Get article ID from URL parameter
	articleId := ctx.Params("id")

	// Find article by ID
	var article entity.Article
	if err := database.DB.Preload("Category").Preload("Author").First(&article, "id = ?", articleId).Error; err != nil {
		// If article not found
		if err == gorm.ErrRecordNotFound {
			return utils.SendErrorResponse(ctx, fiber.StatusNotFound, "Failed to fetch article", err)
		}
		// If error occurred
		return utils.SendErrorResponse(ctx, fiber.StatusInternalServerError, "Failed to fetch article", err)
	}

	return utils.SendSuccessResponse(ctx, fiber.StatusOK, "Succesfully fetched article", fiber.Map{
		"article": article,
	})
}

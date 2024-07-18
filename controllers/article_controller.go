package controllers

import (
	"go-news-api/database"
	"go-news-api/models/entity"
	"go-news-api/models/request"
	"go-news-api/utils"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func GetAllArticles(ctx *fiber.Ctx) error {
	var articles []entity.Article

	// Fetch all articles
	if err := database.DB.Preload("Category").
		Preload("Author").
		Preload("Comments.User").
		Find(&articles).Error; err != nil {
		return utils.SendErrorResponse(ctx, fiber.StatusInternalServerError, "Failed to fetch articles", err)
	}

	return utils.SendSuccessResponseWithData(ctx, fiber.StatusOK, "Successfully fetched articles", fiber.Map{
		"articles":      articles,
		"totalArticles": len(articles),
	})
}

func GetArticleById(ctx *fiber.Ctx) error {
	// Get article ID from URL parameter
	articleId := ctx.Params("id")

	// Find article by ID
	var article entity.Article
	if err := database.DB.Preload("Category").
		Preload("Author").
		Preload("Comments.User").
		First(&article, "id = ?", articleId).Error; err != nil {
		// If article not found
		if err == gorm.ErrRecordNotFound {
			return utils.SendErrorResponse(ctx, fiber.StatusNotFound, "Failed to fetch article", err)
		}
		// If error occurred
		return utils.SendErrorResponse(ctx, fiber.StatusInternalServerError, "Failed to fetch article", err)
	}

	return utils.SendSuccessResponseWithData(ctx, fiber.StatusOK, "Succesfully fetched article", fiber.Map{
		"article": article,
	})
}

func CreateArticle(ctx *fiber.Ctx) error {
	request := new(request.CreateArticleRequest)

	// Parse request body
	if err := ctx.BodyParser(request); err != nil {
		return utils.SendErrorResponse(ctx, fiber.StatusBadRequest, "Failed to create article", err)
	}

	// Validate request
	if err := utils.Validate.Struct(request); err != nil {
		return utils.SendErrorResponse(ctx, fiber.StatusBadRequest, "Failed to create article", err)
	}

	// Save the thumbnail file
	thumbnailPath, err := utils.SaveImageFile(ctx, "thumbnail", "./public/uploads/thumbnails")
	if err != nil {
		return utils.SendErrorResponse(ctx, fiber.StatusInternalServerError, "Failed to save thumbnail", err)
	}

	// Create article
	article := entity.Article{
		Title:      request.Title,
		Slug:       request.Slug,
		Thumbnail:  thumbnailPath,
		Content:    request.Content,
		CategoryID: request.CategoryID,
		AuthorID:   request.AuthorID,
	}

	if err := database.DB.Create(&article).Error; err != nil {
		return utils.SendErrorResponse(ctx, fiber.StatusInternalServerError, "Failed to create article", err)
	}

	if err := database.DB.Preload("Category").Preload("Author").Preload("Comments.User").First(&article, article.ID).Error; err != nil {
		return utils.SendErrorResponse(ctx, fiber.StatusInternalServerError, "Failed to create article", err)
	}

	return utils.SendSuccessResponse(ctx, fiber.StatusCreated, "Successfully created article")
}

func UpdateArticle(ctx *fiber.Ctx) error {
	articleSlug := ctx.Params("slug")

	// Check if article exist
	var article entity.Article
	if err := database.DB.First(&article, "slug = ?", articleSlug).Error; err != nil {
		// If article not found
		if err == gorm.ErrRecordNotFound {
			return utils.SendErrorResponse(ctx, fiber.StatusNotFound, "Failed to update article", err)
		}
		// If error occurred
		return utils.SendErrorResponse(ctx, fiber.StatusInternalServerError, "Failed to update article", err)
	}

	// Parse request body
	request := new(request.UpdateArticleRequest)
	if err := ctx.BodyParser(request); err != nil {
		return utils.SendErrorResponse(ctx, fiber.StatusBadRequest, "Failed to update article", err)
	}

	// Validate request
	if err := utils.Validate.Struct(request); err != nil {
		return utils.SendErrorResponse(ctx, fiber.StatusBadRequest, "Failed to update article", err)
	}

	// Update article
	article.Title = request.Title
	article.Slug = request.Slug
	article.Content = request.Content
	article.CategoryID = request.CategoryID
	article.AuthorID = request.AuthorID

	// Save the thumbnail file if provided
	if _, err := ctx.FormFile("thumbnail"); err == nil {
		// Delete old thumbnail
		if err := utils.DeleteFile(article.Thumbnail); err != nil {
			return utils.SendErrorResponse(ctx, fiber.StatusInternalServerError, "Failed to update article", err)
		}

		// Save new thumbnail
		thumbnailPath, err := utils.SaveImageFile(ctx, "thumbnail", "./public/uploads/thumbnails")
		if err != nil {
			return utils.SendErrorResponse(ctx, fiber.StatusInternalServerError, "Failed to update article", err)
		}
		article.Thumbnail = thumbnailPath
	}

	if err := database.DB.Save(&article).Error; err != nil {
		return utils.SendErrorResponse(ctx, fiber.StatusInternalServerError, "Failed to update article", err)
	}

	if err := database.DB.Preload("Category").Preload("Author").Preload("Comments.User").First(&article, article.ID).Error; err != nil {
		return utils.SendErrorResponse(ctx, fiber.StatusInternalServerError, "Failed to update article", err)
	}

	return utils.SendSuccessResponse(ctx, fiber.StatusOK, "Successfully updated article")
}

func DeleteArticle(ctx *fiber.Ctx) error {
	articleSlug := ctx.Params("slug")

	// Check if article exist
	var article entity.Article
	if err := database.DB.First(&article, "slug = ?", articleSlug).Error; err != nil {
		// If article not found
		if err == gorm.ErrRecordNotFound {
			return utils.SendErrorResponse(ctx, fiber.StatusNotFound, "Failed to delete article", err)
		}
		// If error occurred
		return utils.SendErrorResponse(ctx, fiber.StatusInternalServerError, "Failed to delete article", err)
	}

	// Delete thumbnail
	if err := utils.DeleteFile(article.Thumbnail); err != nil {
		return utils.SendErrorResponse(ctx, fiber.StatusInternalServerError, "Failed to delete article", err)
	}

	// Delete article
	if err := database.DB.Delete(&article).Error; err != nil {
		return utils.SendErrorResponse(ctx, fiber.StatusInternalServerError, "Failed to delete article", err)
	}

	return utils.SendSuccessResponse(ctx, fiber.StatusOK, "Successfully deleted article")
}

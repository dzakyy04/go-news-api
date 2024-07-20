package controllers

import (
	"go-news-api/database"
	"go-news-api/models/entity"
	"go-news-api/models/request"
	"go-news-api/utils"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// GetAllArticles godoc
// @Summary Get all articles
// @Description Retrieves a list of all articles along with their related category, author, comments, and tags.
// @Tags Articles
// @Accept  json
// @Produce  json
// @Router /articles [get]
func GetAllArticles(ctx *fiber.Ctx) error {
	var articles []entity.Article

	// Fetch all articles
	if err := database.DB.Preload("Category").
		Preload("Author").
		Preload("Comments.User").
		Preload("Tags").
		Find(&articles).Error; err != nil {
		return utils.SendErrorResponse(ctx, fiber.StatusInternalServerError, "Failed to fetch articles", err)
	}

	return utils.SendSuccessResponseWithData(ctx, fiber.StatusOK, "Successfully fetched articles", fiber.Map{
		"articles":       articles,
		"total_articles": len(articles),
	})
}

// GetArticleBySlug godoc
// @Summary Get an article by its slug
// @Description Retrieves a single article based on the provided slug, including its related category, author, comments, and tags.
// @Tags Articles
// @Accept  json
// @Produce  json
// @Param slug path string true "Article Slug"
// @Router /articles/{slug} [get]
func GetArticleBySlug(ctx *fiber.Ctx) error {
	// Get article slug from URL parameter
	articleSlug := ctx.Params("slug")

	// Find article by slug
	var article entity.Article
	if err := database.DB.Preload("Category").
		Preload("Author").
		Preload("Comments.User").
		Preload("Tags").
		First(&article, "slug = ?", articleSlug).Error; err != nil {
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

// CreateArticle godoc
// @Summary Create a new article
// @Description Creates a new article with the provided title, slug, content, category, author, thumbnail, and tags. The thumbnail is uploaded as a file and saved to the server.
// @Tags Articles
// @Accept  multipart/form-data
// @Produce  json
// @Param title formData string true "Article Title"
// @Param slug formData string true "Article Slug"
// @Param thumbnail formData file true "Article Thumbnail"
// @Param content formData string true "Article Content"
// @Param category_id formData int true "Category ID"
// @Param author_id formData int true "Author ID"
// @Param tags formData []string true "Article Tags (can be multiple)" collectionFormat(multi)
// @Router /articles [post]
func CreateArticle(ctx *fiber.Ctx) error {
	request := new(request.ArticleRequest)

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

	// Handle tags
	tags, err := utils.CreateOrFindTags(request.Tags)
	if err != nil {
		return utils.SendErrorResponse(ctx, fiber.StatusInternalServerError, "Failed to create article", err)
	}

	if err := database.DB.Create(&article).Error; err != nil {
		return utils.SendErrorResponse(ctx, fiber.StatusInternalServerError, "Failed to create article", err)
	}

	// Associate tags
	if err := utils.AssociateTagsWithArticle(article.ID, tags); err != nil {
		return utils.SendErrorResponse(ctx, fiber.StatusInternalServerError, "Failed to create article", err)
	}

	return utils.SendSuccessResponse(ctx, fiber.StatusCreated, "Successfully created article")
}

// UpdateArticle godoc
// @Summary Update an existing article by its slug
// @Description Updates the details of an existing article, including its title, slug, content, category, author, thumbnail, and tags. If a new thumbnail is provided, the old one will be replaced.
// @Tags Articles
// @Accept  multipart/form-data
// @Produce  json
// @Param slug path string true "Article Slug"
// @Param title formData string true "Article Title"
// @Param slug formData string true "Article Slug"
// @Param thumbnail formData file true "Article Thumbnail"
// @Param content formData string true "Article Content"
// @Param category_id formData int true "Category ID"
// @Param author_id formData int true "Author ID"
// @Param tags formData []string true "Article Tags (can be multiple)" collectionFormat(multi)
// @Router /articles/{slug} [put]
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
	request := new(request.ArticleRequest)
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

	// Handle tags
	tags, err := utils.CreateOrFindTags(request.Tags)
	if err != nil {
		return utils.SendErrorResponse(ctx, fiber.StatusInternalServerError, "Failed to create article", err)
	}

	if err := database.DB.Save(&article).Error; err != nil {
		return utils.SendErrorResponse(ctx, fiber.StatusInternalServerError, "Failed to update article", err)
	}

	// Associate tags
	if err := utils.AssociateTagsWithArticle(article.ID, tags); err != nil {
		return utils.SendErrorResponse(ctx, fiber.StatusInternalServerError, "Failed to create article", err)
	}

	return utils.SendSuccessResponse(ctx, fiber.StatusOK, "Successfully updated article")
}

// DeleteArticle godoc
// @Summary Delete an article by its slug
// @Description Deletes an article specified by the slug from the database. Also deletes the associated thumbnail image from the server.
// @Tags Articles
// @Produce  json
// @Param slug path string true "Article Slug"
// @Router /articles/{slug} [delete]
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

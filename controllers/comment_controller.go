package controllers

import (
	"errors"
	"go-news-api/database"
	"go-news-api/models/entity"
	"go-news-api/models/request"
	"go-news-api/utils"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func CreateComment(ctx *fiber.Ctx) error {
	// Get User
	user := ctx.Locals("user").(*entity.User)
	if user == nil {
		return utils.SendErrorResponse(ctx, fiber.StatusUnauthorized, "Failed to create comment", errors.New("user not found"))
	}

	// Check if article exist
	articleSlug := ctx.Params("slug")
	var article entity.Article
	if err := database.DB.First(&article, "slug = ?", articleSlug).Error; err != nil {
		// If article not found
		if err == gorm.ErrRecordNotFound {
			return utils.SendErrorResponse(ctx, fiber.StatusNotFound, "Failed to create comment", err)
		}
		// If error occurred
		return utils.SendErrorResponse(ctx, fiber.StatusInternalServerError, "Failed to create comment", err)
	}

	// Parse request body
	request := new(request.CreateCommentRequest)
	if err := ctx.BodyParser(request); err != nil {
		return utils.SendErrorResponse(ctx, fiber.StatusBadRequest, "Failed to create comment", err)
	}

	// Validate request
	if err := utils.Validate.Struct(request); err != nil {
		return utils.SendErrorResponse(ctx, fiber.StatusBadRequest, "Failed to create comment", err)
	}

	// Create comment
	comment := entity.Comment{
		ArticleID: article.ID,
		UserID:    user.ID,
		Content:   request.Content,
	}

	if err := database.DB.Create(&comment).Error; err != nil {
		return utils.SendErrorResponse(ctx, fiber.StatusInternalServerError, "Failed to create comment", err)
	}

	if err := database.DB.Preload("User").First(&comment, "id = ?", comment.ID).Error; err != nil {
		return utils.SendErrorResponse(ctx, fiber.StatusInternalServerError, "Failed to create comment", err)
	}

	return utils.SendSuccessResponse(ctx, fiber.StatusCreated, "Successfully created comment", fiber.Map{
		"comment": comment,
	})
}

func DeleteComment(ctx *fiber.Ctx) error {
	// Get User
	user := ctx.Locals("user").(*entity.User)
	if user == nil {
		return utils.SendErrorResponse(ctx, fiber.StatusUnauthorized, "Failed to delete comment", errors.New("user not found"))
	}

	// Check if comment exist
	commentID := ctx.Params("id")
	var comment entity.Comment
	if err := database.DB.First(&comment, "id = ?", commentID).Error; err != nil {
		// If comment not found
		if err == gorm.ErrRecordNotFound {
			return utils.SendErrorResponse(ctx, fiber.StatusNotFound, "Failed to delete comment", err)
		}
		// If error occurred
		return utils.SendErrorResponse(ctx, fiber.StatusInternalServerError, "Failed to delete comment", err)
	}

	// Delete comment
	if err := database.DB.Delete(&comment).Error; err != nil {
		return utils.SendErrorResponse(ctx, fiber.StatusInternalServerError, "Failed to delete comment", err)
	}

	return utils.SendSuccessResponse(ctx, fiber.StatusOK, "Successfully deleted comment", nil)
}

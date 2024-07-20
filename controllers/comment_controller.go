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

// CreateComment godoc
// @Summary Create a comment for an article
// @Description Creates a new comment for the specified article. Requires user to be authenticated and the article to exist.
// @Tags Comments
// @Produce  json
// @Param slug path string true "Article Slug"
// @Param content formData string true "Comment content"
// @Router /articles/{slug}/comments [post]
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
	request := new(request.CommentRequest)
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

	return utils.SendSuccessResponse(ctx, fiber.StatusCreated, "Successfully created comment")
}

// UpdateComment godoc
// @Summary Update an existing comment
// @Description Updates the specified comment if the user is the owner. Requires user to be authenticated and the comment to exist.
// @Tags Comments
// @Produce  json
// @Param id path string true "Comment ID"
// @Param content formData string true "Updated comment content"
// @Router /comments/{id} [put]
func UpdateComment(ctx *fiber.Ctx) error {
	// Get User
	user := ctx.Locals("user").(*entity.User)
	if user == nil {
		return utils.SendErrorResponse(ctx, fiber.StatusUnauthorized, "Failed to update comment", errors.New("user not found"))
	}

	// Check if comment exist
	commentID := ctx.Params("id")
	var comment entity.Comment
	if err := database.DB.First(&comment, "id = ?", commentID).Error; err != nil {
		// If comment not found
		if err == gorm.ErrRecordNotFound {
			return utils.SendErrorResponse(ctx, fiber.StatusNotFound, "Failed to update comment", err)
		}
		// If error occurred
		return utils.SendErrorResponse(ctx, fiber.StatusInternalServerError, "Failed to update comment", err)
	}

	// Check if the user is the owner of the comment
	if comment.UserID != user.ID {
		return utils.SendErrorResponse(ctx, fiber.StatusForbidden, "Failed to update comment", errors.New("you are not allowed to update this comment"))
	}

	// Parse request body
	request := new(request.CommentRequest)
	if err := ctx.BodyParser(request); err != nil {
		return utils.SendErrorResponse(ctx, fiber.StatusBadRequest, "Failed to update comment", err)
	}

	// Validate request
	if err := utils.Validate.Struct(request); err != nil {
		return utils.SendErrorResponse(ctx, fiber.StatusBadRequest, "Failed to update comment", err)
	}

	// Update comment
	comment.Content = request.Content

	if err := database.DB.Save(&comment).Error; err != nil {
		return utils.SendErrorResponse(ctx, fiber.StatusInternalServerError, "Failed to update comment", err)
	}

	return utils.SendSuccessResponse(ctx, fiber.StatusOK, "Successfully updated comment")
}

// DeleteComment godoc
// @Summary Delete an existing comment
// @Description Deletes an existing comment. Requires user to be authenticated and authorized to delete the comment.
// @Tags Comments
// @Produce  json
// @Param id path string true "Comment ID"
// @Router /comments/{id} [delete]
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

	// Check if the user is the owner of the comment
	if comment.UserID != user.ID {
		return utils.SendErrorResponse(ctx, fiber.StatusForbidden, "Failed to delete comment", errors.New("you are not allowed to delete this comment"))
	}

	// Delete comment
	if err := database.DB.Delete(&comment).Error; err != nil {
		return utils.SendErrorResponse(ctx, fiber.StatusInternalServerError, "Failed to delete comment", err)
	}

	return utils.SendSuccessResponse(ctx, fiber.StatusOK, "Successfully deleted comment")
}

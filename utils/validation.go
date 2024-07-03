package utils

import (
	"go-news-api/database"
	"go-news-api/models/entity"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

var Validate *validator.Validate

func init() {
	Validate = validator.New()
	Validate.RegisterValidation("category_exists", CategoryExists)
	Validate.RegisterValidation("author_exists", AuthorExists)
}

func CategoryExists(fl validator.FieldLevel) bool {
	categoryID := fl.Field().Uint()
	var category entity.Category
	if err := database.DB.First(&category, "id = ?", categoryID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return false
		}
	}
	return true
}

func AuthorExists(fl validator.FieldLevel) bool {
	authorID := fl.Field().Uint()
	var author entity.User
	if err := database.DB.First(&author, "id = ?", authorID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return false
		}
	}
	return true
}

func FormatValidationError(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return err.Field() + " is required"
	case "min":
		return err.Field() + " must be at least " + err.Param() + " characters long"
	case "max":
		return err.Field() + " must be at most " + err.Param() + " characters long"
	case "email":
		return err.Field() + " must be a valid email address"
	case "eqfield":
		return err.Field() + " must be equal to " + err.Param()
	case "category_exists":
		return "Category does not exist"
	case "author_exists":
		return "Author does not exist"
	default:
		return err.Error()
	}
}

package utils

import "github.com/go-playground/validator/v10"

func FormatValidationError(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return err.Field() + " is required"
	case "min":
		return err.Field() + " must be at least " + err.Param() + " characters long"
	case "max":
		return err.Field() + " must be at most " + err.Param() + " characters long"
	default:
		return err.Error()
	}
}

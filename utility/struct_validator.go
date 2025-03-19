package utils

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

// Validator instance
var validateInstance = validator.New()

// ValidationErrors represents a map of field names to error messages.
type ValidationErrors map[string]string

// ParseValidationErrors parses validator.ValidationErrors into a ValidationErrors map.
func ParseValidationErrors(err error) ValidationErrors {
	validationErrors := make(ValidationErrors)
	if errs, ok := err.(validator.ValidationErrors); ok {
		for _, e := range errs {
			fieldName := e.Field()
			switch e.Tag() {
			case "required":
				validationErrors[fieldName] = fmt.Sprintf("%s is required", fieldName)
			case "email":
				validationErrors[fieldName] = fmt.Sprintf("%s must be a valid email", fieldName)
			case "min":
				validationErrors[fieldName] = fmt.Sprintf("%s must be at least %s characters long", fieldName, e.Param())
			case "max":
				validationErrors[fieldName] = fmt.Sprintf("%s must be at most %s characters long", fieldName, e.Param())
			case "gt":
				validationErrors[fieldName] = fmt.Sprintf("%s must be greater than %s", fieldName, e.Param())
			case "strongpwd":
				validationErrors[fieldName] = fmt.Sprintf("%s must contain at least one uppercase letter, one lowercase letter, and one number", fieldName)
			default:
				validationErrors[fieldName] = fmt.Sprintf("%s is not valid", fieldName)
			}
		}
	}
	return validationErrors
}

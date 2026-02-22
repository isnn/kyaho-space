package common

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

// Validate is the singleton validator instance for the entire application.
var Validate = validator.New()

// ValidateStruct validates a struct and returns human-readable error messages.
// Returns nil if validation passes, or a slice of readable strings like:
//
//	["target is required", "target_year must be at least 2000"]
func ValidateStruct(s interface{}) []string {
	err := Validate.Struct(s)
	if err == nil {
		return nil
	}

	var errs []string
	for _, e := range err.(validator.ValidationErrors) {
		errs = append(errs, formatFieldError(e))
	}
	return errs
}

// formatFieldError converts a single validator.FieldError into a human-readable string.
func formatFieldError(e validator.FieldError) string {
	field := toSnakeCase(e.Field())

	switch e.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", field)
	case "min":
		return fmt.Sprintf("%s must be at least %s", field, e.Param())
	case "max":
		return fmt.Sprintf("%s must be at most %s", field, e.Param())
	case "uuid":
		return fmt.Sprintf("%s must be a valid UUID", field)
	case "email":
		return fmt.Sprintf("%s must be a valid email address", field)
	case "url":
		return fmt.Sprintf("%s must be a valid URL", field)
	case "oneof":
		return fmt.Sprintf("%s must be one of: %s", field, e.Param())
	default:
		return fmt.Sprintf("%s failed validation: %s", field, e.Tag())
	}
}

// toSnakeCase converts PascalCase/camelCase field names to snake_case for readability.
func toSnakeCase(s string) string {
	var result strings.Builder
	for i, r := range s {
		if r >= 'A' && r <= 'Z' {
			if i > 0 {
				result.WriteRune('_')
			}
			result.WriteRune(r + 32) // tolower
		} else {
			result.WriteRune(r)
		}
	}
	return result.String()
}

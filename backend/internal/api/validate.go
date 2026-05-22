package api

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

// ValidateStruct runs validator tags and returns user-facing error strings.
func ValidateStruct(s interface{}) []string {
	err := validate.Struct(s)
	if err == nil {
		return nil
	}
	var errs []string
	for _, e := range err.(validator.ValidationErrors) {
		errs = append(errs, formatValidationError(e))
	}
	return errs
}

func formatValidationError(e validator.FieldError) string {
	field := strings.ToLower(e.Field())
	switch e.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", field)
	case "min":
		return fmt.Sprintf("%s is too short", field)
	case "max":
		return fmt.Sprintf("%s is too long", field)
	case "oneof":
		return fmt.Sprintf("%s has an invalid value", field)
	case "url":
		return "One or more official links is not a valid URL"
	case "uuid":
		return fmt.Sprintf("%s must be a valid identifier", field)
	default:
		return fmt.Sprintf("%s failed validation (%s)", field, e.Tag())
	}
}

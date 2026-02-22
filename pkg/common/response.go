package common

import (
	"net/http"

	"github.com/gofiber/fiber/v3"
)

// ErrorBody is the standard error response shape.
// Status code comes from the HTTP header — not duplicated in the body.
type ErrorBody struct {
	Message string   `json:"message"`
	Details []string `json:"details,omitempty"`
}

// --- Error Responses ---

// BadRequestResponse returns 400 with a message.
func BadRequestResponse(c fiber.Ctx, msg string) error {
	return c.Status(http.StatusBadRequest).JSON(ErrorBody{Message: msg})
}

// ValidationErrorResponse returns 400 with a message and field-level details.
func ValidationErrorResponse(c fiber.Ctx, msg string, details []string) error {
	return c.Status(http.StatusBadRequest).JSON(ErrorBody{Message: msg, Details: details})
}

// UnauthorizedResponse returns 401.
func UnauthorizedResponse(c fiber.Ctx, msg string) error {
	return c.Status(http.StatusUnauthorized).JSON(ErrorBody{Message: msg})
}

// ForbiddenResponse returns 403.
func ForbiddenResponse(c fiber.Ctx, msg string) error {
	return c.Status(http.StatusForbidden).JSON(ErrorBody{Message: msg})
}

// NotFoundResponse returns 404.
func NotFoundResponse(c fiber.Ctx, msg string) error {
	return c.Status(http.StatusNotFound).JSON(ErrorBody{Message: msg})
}

// InternalErrorResponse returns 500.
func InternalErrorResponse(c fiber.Ctx, msg string) error {
	return c.Status(http.StatusInternalServerError).JSON(ErrorBody{Message: msg})
}

// --- Success Responses ---

// SuccessResponse returns 200 with data.
func SuccessResponse(c fiber.Ctx, data interface{}) error {
	return c.Status(http.StatusOK).JSON(data)
}

// CreatedResponse returns 201 with data.
func CreatedResponse(c fiber.Ctx, data interface{}) error {
	return c.Status(http.StatusCreated).JSON(data)
}

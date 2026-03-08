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
func BadRequestResponse(c fiber.Ctx, msg ...string) error {
	message := "Bad Request"
	if len(msg) > 0 && msg[0] != "" {
		message = msg[0]
	}
	return c.Status(http.StatusBadRequest).JSON(ErrorBody{Message: message})
}

// ValidationErrorResponse returns 400 with a message and field-level details.
func ValidationErrorResponse(c fiber.Ctx, msg string, details []string) error {
	if msg == "" {
		msg = "Validation Failed"
	}
	return c.Status(http.StatusBadRequest).JSON(ErrorBody{Message: msg, Details: details})
}

// UnauthorizedResponse returns 401.
func UnauthorizedResponse(c fiber.Ctx, msg ...string) error {
	message := "Unauthorized"
	if len(msg) > 0 && msg[0] != "" {
		message = msg[0]
	}
	return c.Status(http.StatusUnauthorized).JSON(ErrorBody{Message: message})
}

// ForbiddenResponse returns 403.
func ForbiddenResponse(c fiber.Ctx, msg ...string) error {
	message := "Forbidden"
	if len(msg) > 0 && msg[0] != "" {
		message = msg[0]
	}
	return c.Status(http.StatusForbidden).JSON(ErrorBody{Message: message})
}

// NotFoundResponse returns 404.
func NotFoundResponse(c fiber.Ctx, msg ...string) error {
	message := "Not Found"
	if len(msg) > 0 && msg[0] != "" {
		message = msg[0]
	}
	return c.Status(http.StatusNotFound).JSON(ErrorBody{Message: message})
}

// InternalErrorResponse returns 500.
func InternalErrorResponse(c fiber.Ctx, msg ...string) error {
	message := "Internal Server Error"
	if len(msg) > 0 && msg[0] != "" {
		message = msg[0]
	}
	return c.Status(http.StatusInternalServerError).JSON(ErrorBody{Message: message})
}

// --- Success Responses ---

// SuccessBody is the standard success response shape.
type SuccessBody struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// SuccessResponse returns 200 with data.
func SuccessResponse(c fiber.Ctx, data interface{}, msg ...string) error {
	message := "success"
	if len(msg) > 0 && msg[0] != "" {
		message = msg[0]
	}
	return c.Status(http.StatusOK).JSON(SuccessBody{
		Message: message,
		Data:    data,
	})
}

// CreatedResponse returns 201 with data.
func CreatedResponse(c fiber.Ctx, data interface{}, msg ...string) error {
	message := "success"
	if len(msg) > 0 && msg[0] != "" {
		message = msg[0]
	}
	return c.Status(http.StatusCreated).JSON(SuccessBody{
		Message: message,
		Data:    data,
	})
}

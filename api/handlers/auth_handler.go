package handlers

import (
	. "kyaho-space/pkg/common"
	"kyaho-space/pkg/entities"
	"kyaho-space/pkg/user"
	"time"

	"github.com/gofiber/fiber/v3"
)

func Signup(service user.Service) fiber.Handler {
	return func(c fiber.Ctx) error {
		var req entities.SignupRequest
		if err := c.Bind().Body(&req); err != nil {
			return BadRequestResponse(c, "Invalid request body")
		}

		if errs := ValidateStruct(req); errs != nil {
			return ValidationErrorResponse(c, "Validation failed", errs)
		}

		u, err := service.Register(req)
		if err != nil {
			return BadRequestResponse(c, err.Error())
		}

		return CreatedResponse(c, u)
	}
}

func Login(service user.Service, jwtSecret string, env string) fiber.Handler {
	return func(c fiber.Ctx) error {
		var req entities.LoginRequest
		if err := c.Bind().Body(&req); err != nil {
			return BadRequestResponse(c, "Invalid request body")
		}

		if errs := ValidateStruct(req); errs != nil {
			return ValidationErrorResponse(c, "Validation failed", errs)
		}

		accessToken, refreshToken, err := service.Login(req, jwtSecret)
		if err != nil {
			return UnauthorizedResponse(c, err.Error())
		}

		cookie := &fiber.Cookie{
			Name:     "refresh_token",
			Value:    refreshToken,
			HTTPOnly: true,
			Secure:   env == "production",
			SameSite: fiber.CookieSameSiteLaxMode,
			MaxAge:   int((7 * 24 * time.Hour).Seconds()),
			Path:     "/",
		}

		if env == "production" {
			cookie.SameSite = fiber.CookieSameSiteNoneMode
		}
		c.Cookie(cookie)

		return SuccessResponse(c, entities.AuthLoginResponse{
			AccessToken: accessToken,
		})
	}
}

func Refresh(service user.Service, jwtSecret string) fiber.Handler {
	return func(c fiber.Ctx) error {
		refreshTokenOld := c.Cookies("refresh_token")

		if refreshTokenOld == "" {
			return UnauthorizedResponse(c, "Missing refresh token")
		}

		newAccessToken, err := service.RefreshToken(refreshTokenOld, jwtSecret)
		if err != nil {
			return UnauthorizedResponse(c, err.Error())
		}

		return SuccessResponse(c, entities.AuthRefreshResponse{
			RefreshToken: newAccessToken,
		})
	}
}

func Logout(service user.Service) fiber.Handler {
	return func(c fiber.Ctx) error {
		refreshToken := c.Cookies("refresh_token")
		if refreshToken != "" {
			_ = service.Logout(refreshToken)
		}

		// Clear the cookie
		cookie := &fiber.Cookie{
			Name:     "refresh_token",
			Value:    "",
			Expires:  time.Now().Add(-time.Hour),
			HTTPOnly: true,
			Secure:   true, // Set to true for deletion robustness
			SameSite: fiber.CookieSameSiteNoneMode,
			Path:     "/",
		}
		c.Cookie(cookie)

		return SuccessResponse(c, nil, "Logged out")
	}
}

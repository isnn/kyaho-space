package middleware

import (
	"strings"

	. "kyaho-space/pkg/common"

	"github.com/golang-jwt/jwt/v5"

	"github.com/gofiber/fiber/v3"
)

// JWTAuth returns a Fiber middleware that validates the access token from the
// Authorization header and sets user_id in c.Locals.
func JWTAuth(secret string) fiber.Handler {
	return func(c fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return UnauthorizedResponse(c, "Unauthorized")
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			return UnauthorizedResponse(c, "Unauthorized")
		}

		tokenStr := parts[1]

		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})
		if err != nil || !token.Valid {
			return UnauthorizedResponse(c, "Invalid or expired token")
		}

		userID, ok := claims["user_id"].(string)
		if !ok || userID == "" {
			return UnauthorizedResponse(c, "Invalid token claims")
		}

		c.Locals("user_id", userID)
		return c.Next()
	}
}

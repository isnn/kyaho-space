package middleware

import (
	"kyaho-space/pkg/config"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
)

func CorsSetup(cfg config.Config) fiber.Handler {
	rawOrigins := strings.Split(cfg.FrontendURL, ",")
	var envOrigins []string
	for _, o := range rawOrigins {
		trimmed := strings.TrimSpace(o)
		if trimmed != "" {
			envOrigins = append(envOrigins, trimmed)
		}
	}

	corsConfig := cors.Config{
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"},
		AllowCredentials: true,
	}

	if cfg.Env == "development" {
		corsConfig.AllowOrigins = append([]string{
			"http://localhost:3000",
			"http://127.0.0.1:3000",
		}, envOrigins...)
	} else {
		corsConfig.AllowOrigins = envOrigins
	}

	return cors.New(corsConfig)
}

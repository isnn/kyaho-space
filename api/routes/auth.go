package routes

import (
	"kyaho-space/api/handlers"
	"kyaho-space/pkg/user"

	"github.com/gofiber/fiber/v3"
)

func AuthRouter(api fiber.Router, service user.Service, jwtSecret string, env string) {
	auth := api.Group("/auth")

	auth.Post("/signup", handlers.Signup(service))
	auth.Post("/login", handlers.Login(service, jwtSecret, env))
	auth.Post("/refresh", handlers.Refresh(service, jwtSecret, env))
	auth.Post("/logout", handlers.Logout(service))
}

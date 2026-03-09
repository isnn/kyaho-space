package routes

import (
	"kyaho-space/api/handlers"
	"kyaho-space/pkg/user"

	"github.com/gofiber/fiber/v3"
)

func AuthRouter(api fiber.Router, service user.Service, jwtSecret string) {
	auth := api.Group("/auth")

	auth.Post("/signup", handlers.Signup(service))
	auth.Post("/login", handlers.Login(service, jwtSecret))
	auth.Post("/refresh", handlers.Refresh(service, jwtSecret))
	auth.Post("/logout", handlers.Logout(service))
}

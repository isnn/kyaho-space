package routes

import (
	"kyaho-space/api/handlers"
	"kyaho-space/pkg/job"

	"github.com/gofiber/fiber/v3"
)

func JobRouter(api fiber.Router, service job.Service) {
	api.Get("/jobs", handlers.GetJob(service))
	api.Post("/jobs", handlers.AddJob(service))
}

package routes

import (
	"kyaho-space/api/handlers"
	"kyaho-space/pkg/job"

	"github.com/gofiber/fiber/v3"
)

func JobRouter(api fiber.Router, service job.Service) {
	api.Get("/jobs", handlers.GetJobs(service))
	api.Post("/jobs", handlers.AddJob(service))
	api.Get("/jobs/statistics", handlers.GetJobStatistics(service))
	api.Get("/jobs/:id", handlers.GetJobDetail(service))
	api.Put("/jobs/:id", handlers.UpdateJob(service))
	api.Delete("/jobs/:id", handlers.DeleteJob(service))
}

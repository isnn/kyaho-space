package handlers

import (
	"kyaho-space/pkg/entities"
	"kyaho-space/pkg/job"
	"net/http"

	"github.com/gofiber/fiber/v3"
)

func AddJob(service job.Service) fiber.Handler {
	return func(c fiber.Ctx) error {
		var requestBody entities.JobRequest
		if err := c.Bind().Body(&requestBody); err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		newJob := entities.Job{
			CompanyName: requestBody.CompanyName,
			JobTitle:    requestBody.JobTitle,
		}

		err := service.AddJob(&newJob)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.Status(http.StatusCreated).JSON(newJob)
	}
}

func GetJob(service job.Service) fiber.Handler {
	return func(c fiber.Ctx) error {
		jobs, err := service.GetJob()
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.JSON(jobs)
	}
}

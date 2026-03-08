package handlers

import (
	. "kyaho-space/pkg/common"
	"kyaho-space/pkg/entities"
	"kyaho-space/pkg/job"

	"gorm.io/gorm"

	"github.com/gofiber/fiber/v3"
)

func AddJob(service job.Service) fiber.Handler {
	return func(c fiber.Ctx) error {
		var req entities.JobRequest
		if err := c.Bind().Body(&req); err != nil {
			return BadRequestResponse(c, "Invalid request body")
		}

		if errs := ValidateStruct(req); errs != nil {
			return ValidationErrorResponse(c, "Bad Request", errs)
		}

		newJob := entities.Job{
			CompanyName: req.CompanyName,
			JobTitle:    req.JobTitle,
			JobPostURL:  req.JobPostURL,
			Status:      req.Status,
			Notes:       req.Notes,
		}

		if err := service.AddJob(&newJob); err != nil {
			return InternalErrorResponse(c, "Failed to create job")
		}

		return CreatedResponse(c, newJob)
	}
}

func GetJobs(service job.Service) fiber.Handler {
	return func(c fiber.Ctx) error {
		jobs, err := service.GetJobs()
		if err != nil {
			return InternalErrorResponse(c, "Failed to retrieve jobs")
		}

		return SuccessResponse(c, jobs)
	}
}

func GetJobDetail(service job.Service) fiber.Handler {
	return func(c fiber.Ctx) error {
		id := c.Params("id")
		if id == "" {
			return BadRequestResponse(c, "Job ID is required")
		}

		j, err := service.GetJobByID(id)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return NotFoundResponse(c, "Job not found")
			}
			return InternalErrorResponse(c, "Failed to retrieve job")
		}

		return SuccessResponse(c, j)
	}
}

func UpdateJob(service job.Service) fiber.Handler {
	return func(c fiber.Ctx) error {
		id := c.Params("id")
		if id == "" {
			return BadRequestResponse(c, "Job ID is required")
		}

		var req entities.JobUpdateRequest
		if err := c.Bind().Body(&req); err != nil {
			return BadRequestResponse(c, "Invalid request body")
		}

		if errs := ValidateStruct(req); errs != nil {
			return ValidationErrorResponse(c, "Validation failed", errs)
		}

		updates := make(map[string]interface{})
		if req.CompanyName != nil {
			updates["company_name"] = *req.CompanyName
		}
		if req.JobTitle != nil {
			updates["job_title"] = *req.JobTitle
		}
		if req.JobPostURL != nil {
			updates["job_post_url"] = *req.JobPostURL
		}
		if req.Status != nil {
			updates["status"] = *req.Status
		}
		if req.Notes != nil {
			updates["notes"] = *req.Notes
		}

		if len(updates) == 0 {
			return BadRequestResponse(c, "No fields to update")
		}

		j := entities.Job{ID: id}
		if err := service.UpdateJob(&j, updates); err != nil {
			if err == gorm.ErrRecordNotFound {
				return NotFoundResponse(c, "Job not found")
			}
			return InternalErrorResponse(c, "Failed to update job")
		}

		return SuccessResponse(c, fiber.Map{"message": "Job updated successfully"})
	}
}

func DeleteJob(service job.Service) fiber.Handler {
	return func(c fiber.Ctx) error {
		id := c.Params("id")
		if id == "" {
			return BadRequestResponse(c, "Job ID is required")
		}

		if err := service.DeleteJob(id); err != nil {
			if err == gorm.ErrRecordNotFound {
				return NotFoundResponse(c, "Job not found")
			}
			return InternalErrorResponse(c, "Failed to delete job")
		}

		return SuccessResponse(c, fiber.Map{"message": "Job deleted successfully"})
	}
}

func GetJobStatistics(service job.Service) fiber.Handler {
	return func(c fiber.Ctx) error {
		stats, err := service.GetStatistics()
		if err != nil {
			return InternalErrorResponse(c, "Failed to retrieve job statistics")
		}

		return SuccessResponse(c, stats)
	}
}

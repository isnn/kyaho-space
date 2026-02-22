package handlers

import (
	. "kyaho-space/pkg/common"
	"kyaho-space/pkg/entities"
	"kyaho-space/pkg/meet_plan"

	"github.com/gofiber/fiber/v3"
)

func CreateMeetPlan(service meet_plan.Service) fiber.Handler {
	return func(c fiber.Ctx) error {
		var req entities.MeetPlanRequest
		if err := c.Bind().Body(&req); err != nil {
			return BadRequestResponse(c, "Invalid request body")
		}

		if errs := ValidateStruct(req); errs != nil {
			return ValidationErrorResponse(c, "Bad Request", errs)
		}

		m := entities.MeetPlan{
			Organization: req.Organization,
			Email:        req.Email,
			InquiryType:  req.InquiryType,
			Brief:        req.Brief,
		}

		if err := service.AddMeetPlan(&m); err != nil {
			return InternalErrorResponse(c, "Failed to create meeting plan")
		}

		return CreatedResponse(c, m)
	}
}

func ListMeetPlans(service meet_plan.Service) fiber.Handler {
	return func(c fiber.Ctx) error {
		meetPlans, err := service.GetMeetPlans()
		if err != nil {
			return InternalErrorResponse(c, "Failed to retrieve meeting plans")
		}

		return SuccessResponse(c, meetPlans)
	}
}

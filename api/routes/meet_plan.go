package routes

import (
	"kyaho-space/api/handlers"
	"kyaho-space/pkg/meet_plan"

	"github.com/gofiber/fiber/v3"
)

func MeetPlanRouter(api fiber.Router, service meet_plan.Service) {
	api.Get("/sync", handlers.ListMeetPlans(service))
	api.Post("/sync", handlers.CreateMeetPlan(service))
}

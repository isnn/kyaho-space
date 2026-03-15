package meet_plan

import "kyaho-space/pkg/entities"

type Service interface {
	AddMeetPlan(meetPlan *entities.MeetPlan) error
	GetMeetPlans() ([]entities.MeetPlan, error)
	DeleteMeetPlan(id string) error
}

type service struct {
	repository Repository
}

// NewService creates a single instance of the service
func NewService(r Repository) Service {
	return &service{
		repository: r,
	}
}

func (s *service) AddMeetPlan(meetPlan *entities.MeetPlan) error {
	return s.repository.CreateMeetPlan(meetPlan)
}

func (s *service) GetMeetPlans() ([]entities.MeetPlan, error) {
	return s.repository.ListMeetPlans()
}

func (s *service) DeleteMeetPlan(id string) error {
	return s.repository.DeleteMeetPlan(id)
}

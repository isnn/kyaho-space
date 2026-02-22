package job

import "kyaho-space/pkg/entities"

type Service interface {
	AddJob(job *entities.Job) error
	GetJob() ([]entities.Job, error)
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

func (s *service) AddJob(job *entities.Job) error {
	return s.repository.CreateJob(job)
}

func (s *service) GetJob() ([]entities.Job, error) {
	return s.repository.ReadJob()
}
